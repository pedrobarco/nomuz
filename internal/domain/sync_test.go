package domain_test

import (
	"context"
	"testing"

	"github.com/pedrobarco/nomuz/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestPlanSync(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	tracks := []domain.Track{
		{
			ID:     "t1",
			Title:  "Track 1",
			Artist: "Artist A",
			Album:  "Album X",
		},
		{
			ID:     "t2",
			Title:  "Track 2",
			Artist: "Artist B",
			Album:  "Album Y",
		},
		{
			ID:     "t3",
			Title:  "Track 3",
			Artist: "Artist C",
			Album:  "Album Z",
		},
	}

	t.Run("no playlists to sync", func(t *testing.T) {
		src := &mockConnector{}
		dst := &mockConnector{}

		cl, err := domain.PlanSync(ctx, src, dst)
		assert.NoError(err)
		assert.Len(cl.Playlists.Added, 0)
		assert.Len(cl.TracksByPlaylist, 0)
	})

	t.Run("no changes to sync", func(t *testing.T) {
		src := &mockConnector{
			Playlists: []*domain.Playlist{
				{
					ID:     "pl1",
					Name:   "Playlist 1",
					Tracks: tracks,
				},
			},
		}

		dst := &mockConnector{
			Tracks: tracks,
			Playlists: []*domain.Playlist{
				{
					ID:     "pl1",
					Name:   "Playlist 1",
					Tracks: tracks,
				},
			},
		}

		cl, err := domain.PlanSync(ctx, src, dst)
		assert.NoError(err)
		assert.Len(cl.Playlists.Added, 0)
		assert.Len(cl.TracksByPlaylist, 0)
	})

	t.Run("full sync", func(t *testing.T) {
		src := &mockConnector{
			Playlists: []*domain.Playlist{
				{
					ID:   "pl1",
					Name: "Playlist 1",
					Tracks: append(tracks, domain.Track{
						ID:     "t4",
						Title:  "Rare Track",
						Artist: "Unknown Artist",
						Album:  "Unknown Album",
					}),
				},
				{
					ID:     "pl2",
					Name:   "Playlist 2",
					Tracks: tracks,
				},
			},
		}

		dst := &mockConnector{
			Tracks: tracks,
			Playlists: []*domain.Playlist{
				{
					ID:   "pl1",
					Name: "Playlist 1",
					Tracks: []domain.Track{
						{
							ID:     "t10",
							Title:  "Track 10",
							Artist: "Artist AA",
							Album:  "Album XX",
						},
					},
				},
			},
		}

		ref0 := domain.PlaylistRef{
			ID:   src.Playlists[0].ID,
			Name: src.Playlists[0].Name,
		}

		ref1 := domain.PlaylistRef{
			ID:   src.Playlists[1].ID,
			Name: src.Playlists[1].Name,
		}

		cl, err := domain.PlanSync(ctx, src, dst)
		assert.NoError(err)
		assert.Len(cl.Playlists.Added, 1)
		assert.Equal("pl2", cl.Playlists.Added[0].ID)
		assert.Equal("Playlist 2", cl.Playlists.Added[0].Name)
		assert.Len(cl.TracksByPlaylist[ref0].Added, 3)
		assert.Len(cl.TracksByPlaylist[ref0].Removed, 1)
		assert.Len(cl.TracksByPlaylist[ref0].Missing, 1)
		assert.Len(cl.TracksByPlaylist[ref1].Added, 3)
		assert.Len(cl.TracksByPlaylist[ref1].Removed, 0)
		assert.Len(cl.TracksByPlaylist[ref1].Missing, 0)
	})
}
