package domain

type TrackFilters struct {
	ID string
}

type Connector interface {
	GetPlaylists() ([]*Playlist, error)
	GetPlaylistByName(name string) (*Playlist, error)
	SavePlaylist(pl *Playlist) error
	SearchTrack(filters TrackFilters) ([]Track, error)
}
