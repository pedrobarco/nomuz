package domain_test

import (
	"testing"

	"github.com/pedrobarco/nomuz/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestSync(t *testing.T) {
	assert := assert.New(t)

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

		cl, err := domain.PlanSync(src, dst)
		assert.NoError(err)
		assert.Len(cl, 0)
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

		cl, err := domain.PlanSync(src, dst)
		assert.NoError(err)
		assert.Len(cl, 1)
		assert.False(cl[0].HasChanges())
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

		cl, err := domain.PlanSync(src, dst)
		assert.NoError(err)
		assert.Len(cl, 2)
		assert.True(cl[0].HasChanges())
		assert.False(cl[0].Added)
		assert.Len(cl[0].Tracks.Added, 3)
		assert.Len(cl[0].Tracks.Removed, 1)
		assert.Len(cl[0].Tracks.Missing, 1)
		assert.True(cl[1].HasChanges())
		assert.True(cl[1].Added)
		assert.Len(cl[1].Tracks.Added, 3)
	})
}
