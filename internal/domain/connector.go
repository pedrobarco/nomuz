package domain

import "context"

type TrackFilters struct {
	ID string
}

type Connector interface {
	CreatePlaylist(ctx context.Context, name string) (*Playlist, error)
	GetPlaylists(ctx context.Context) ([]*Playlist, error)
	GetPlaylistByName(ctx context.Context, name string) (*Playlist, error)
	AddTracksToPlaylist(ctx context.Context, id string, tracks []Track) error
	DeleteTracksFromPlaylist(ctx context.Context, id string, tracks []Track) error
	SearchTrack(ctx context.Context, filters TrackFilters) ([]Track, error)
}
