package domain

import (
	"fmt"
	"log/slog"
)

type trackChangelog struct {
	Added   []Track
	Removed []Track
	Missing []Track
}

type playlistChangelog struct {
	Name   string
	Added  bool
	Tracks trackChangelog
}

func (cl *playlistChangelog) HasChanges() bool {
	return cl.Added ||
		len(cl.Tracks.Added) > 0 ||
		len(cl.Tracks.Removed) > 0 ||
		len(cl.Tracks.Missing) > 0
}

func PlanSync(from, to Connector) ([]playlistChangelog, error) {
	pls, err := from.GetPlaylists()
	if err != nil {
		return nil, fmt.Errorf("failed to get playlists from source: %w", err)
	}

	var cls []playlistChangelog
	for _, pl := range pls {
		slog.Info("syncing playlist", "id", pl.ID, "name", pl.Name)
		cl, err := syncPlaylist(*pl, to)
		if err != nil {
			return nil, fmt.Errorf("failed to sync playlist %s: %w", pl.Name, err)
		}

		cls = append(cls, *cl)
	}

	for _, cl := range cls {
		slog.Info("playlist changelog",
			"name", cl.Name,
			"added", cl.Added,
			"tracks_added", len(cl.Tracks.Added),
			"tracks_removed", len(cl.Tracks.Removed),
			"tracks_missing", len(cl.Tracks.Missing),
		)
	}

	return cls, nil
}

func syncPlaylist(pl Playlist, to Connector) (*playlistChangelog, error) {
	cl := new(playlistChangelog)
	cl.Name = pl.Name

	dst, err := to.GetPlaylistByName(pl.Name)
	if err != nil {
		dst = &Playlist{
			ID:   pl.ID,
			Name: pl.Name,
		}
		cl.Added = true
	}

	dstLookup := make(map[string]Track)
	for _, tr := range dst.Tracks {
		dstLookup[tr.ID] = tr
	}

	srcLookup := make(map[string]Track)
	for _, tr := range pl.Tracks {
		srcLookup[tr.ID] = tr

		if _, found := dstLookup[tr.ID]; found {
			continue
		}

		tracks, err := to.SearchTrack(TrackFilters{ID: tr.ID})
		if err != nil {
			return nil, fmt.Errorf("failed to search track %s in destination: %w", tr.ID, err)
		}

		if len(tracks) == 0 {
			cl.Tracks.Missing = append(cl.Tracks.Missing, tr)
			slog.Info("track not found in destination, skipping", "isrc", tr.ID)
			continue
		}

		cl.Tracks.Added = append(cl.Tracks.Added, tracks[0])
	}

	for _, tr := range dst.Tracks {
		if _, found := srcLookup[tr.ID]; !found {
			cl.Tracks.Removed = append(cl.Tracks.Removed, tr)
		}
	}

	return cl, nil
}
