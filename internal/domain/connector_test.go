package domain_test

import (
	"fmt"

	"github.com/pedrobarco/nomuz/internal/domain"
)

type mockConnector struct {
	Playlists []*domain.Playlist
	Tracks    []domain.Track
}

var _ domain.Connector = (*mockConnector)(nil)

func (m *mockConnector) GetPlaylists() ([]*domain.Playlist, error) {
	return m.Playlists, nil
}

func (m *mockConnector) GetPlaylistByName(name string) (*domain.Playlist, error) {
	for _, pl := range m.Playlists {
		if pl.Name == name {
			return pl, nil
		}
	}
	return nil, fmt.Errorf("playlist not found")
}

func (m *mockConnector) SavePlaylist(pl *domain.Playlist) error {
	for i, existing := range m.Playlists {
		if existing.ID == pl.ID {
			m.Playlists[i] = pl
			return nil
		}
	}
	m.Playlists = append(m.Playlists, pl)
	return nil
}

func (m *mockConnector) SearchTrack(filters domain.TrackFilters) ([]domain.Track, error) {
	for _, tr := range m.Tracks {
		if tr.ID == filters.ID {
			return []domain.Track{tr}, nil
		}
	}
	return nil, nil
}
