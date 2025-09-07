package domain_test

import (
	"context"
	"fmt"

	"github.com/pedrobarco/nomuz/internal/domain"
)

type mockConnector struct {
	Playlists []*domain.Playlist
	Tracks    []domain.Track
}

var _ domain.Connector = (*mockConnector)(nil)

func (m *mockConnector) CreatePlaylist(ctx context.Context, name string) (*domain.Playlist, error) {
	pl := &domain.Playlist{
		ID:     fmt.Sprintf("pl%d", len(m.Playlists)),
		Name:   name,
		Tracks: []domain.Track{},
	}
	m.Playlists = append(m.Playlists, pl)
	return pl, nil
}

func (m *mockConnector) GetPlaylists(ctx context.Context) ([]*domain.Playlist, error) {
	return m.Playlists, nil
}

func (m *mockConnector) GetPlaylistByName(ctx context.Context, name string) (*domain.Playlist, error) {
	for _, pl := range m.Playlists {
		if pl.Name == name {
			return pl, nil
		}
	}
	return nil, nil
}

func (m *mockConnector) SearchTrack(ctx context.Context, filters domain.TrackFilters) ([]domain.Track, error) {
	for _, tr := range m.Tracks {
		if tr.ID == filters.ID {
			return []domain.Track{tr}, nil
		}
	}
	return nil, nil
}

func (m *mockConnector) AddTracksToPlaylist(ctx context.Context, id string, tracks []domain.Track) error {
	for _, pl := range m.Playlists {
		if pl.ID == id {
			pl.Tracks = append(pl.Tracks, tracks...)
			return nil
		}
	}
	return fmt.Errorf("playlist with id %s not found", id)
}

func (m *mockConnector) DeleteTracksFromPlaylist(ctx context.Context, id string, tracks []domain.Track) error {
	for _, pl := range m.Playlists {
		if pl.ID == id {
			trackMap := make(map[string]struct{})
			for _, tr := range tracks {
				trackMap[tr.ID] = struct{}{}
			}

			var newTracks []domain.Track
			for _, tr := range pl.Tracks {
				if _, found := trackMap[tr.ID]; !found {
					newTracks = append(newTracks, tr)
				}
			}

			pl.Tracks = newTracks
			return nil
		}
	}
	return fmt.Errorf("playlist with id %s not found", id)
}
