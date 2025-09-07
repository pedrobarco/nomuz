package domain

import (
	"context"
	"fmt"
	"log/slog"
)

type playlistTracksChangelog struct {
	Added   []Track
	Removed []Track
	Missing []Track
}

type playlistChangelog struct {
	Added []PlaylistRef
}

type changelog struct {
	Playlists        playlistChangelog
	TracksByPlaylist map[PlaylistRef]playlistTracksChangelog
}

type PlaylistRef struct {
	ID   string
	Name string
}

func PlanSync(ctx context.Context, from, to Connector) (*changelog, error) {
	pls, err := from.GetPlaylists(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get playlists from source: %w", err)
	}

	changelog := &changelog{
		Playlists:        playlistChangelog{},
		TracksByPlaylist: make(map[PlaylistRef]playlistTracksChangelog),
	}

	for _, src := range pls {
		slog.Info("syncing playlist", "id", src.ID, "name", src.Name)

		dst, err := to.GetPlaylistByName(ctx, src.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to get playlist %s from destination: %w", src.Name, err)
		}

		if dst == nil {
			dst = &Playlist{
				ID:   src.ID,
				Name: src.Name,
			}
			changelog.Playlists.Added = append(changelog.Playlists.Added, PlaylistRef{
				ID:   src.ID,
				Name: src.Name,
			})
		}

		cl, err := syncPlaylist(ctx, *src, *dst, to)
		if err != nil {
			return nil, fmt.Errorf("failed to sync playlist %s: %w", src.Name, err)
		}

		ref := PlaylistRef{
			ID:   src.ID,
			Name: src.Name,
		}
		changelog.TracksByPlaylist[ref] = *cl
	}

	return changelog, nil
}

func Sync(ctx context.Context, from, to Connector, cl changelog) error {
	createdPl := make(map[string]PlaylistRef)
	for _, ref := range cl.Playlists.Added {
		pl, err := to.CreatePlaylist(ctx, ref.Name)
		if err != nil {
			return fmt.Errorf("failed to create playlist %s: %w", ref.Name, err)
		}
		createdPl[ref.Name] = PlaylistRef{
			ID:   pl.ID,
			Name: pl.Name,
		}
	}

	for ref, tracks := range cl.TracksByPlaylist {
		if _, found := createdPl[ref.Name]; found {
			ref = createdPl[ref.Name]
		}

		if len(tracks.Added) > 0 {
			if err := to.AddTracksToPlaylist(ctx, ref.ID, tracks.Added); err != nil {
				return fmt.Errorf("failed to add tracks to playlist %s: %w", ref.Name, err)
			}
		}

		if len(tracks.Removed) > 0 {
			if err := to.DeleteTracksFromPlaylist(ctx, ref.ID, tracks.Removed); err != nil {
				return fmt.Errorf("failed to remove tracks from playlist %s: %w", ref.Name, err)
			}
		}
	}

	return nil
}

func syncPlaylist(ctx context.Context, src, dst Playlist, to Connector) (*playlistTracksChangelog, error) {
	dstLookup := make(map[string]Track)
	for _, tr := range dst.Tracks {
		dstLookup[tr.ID] = tr
	}

	cl := new(playlistTracksChangelog)

	srcLookup := make(map[string]Track)
	for _, tr := range src.Tracks {
		srcLookup[tr.ID] = tr

		if _, found := dstLookup[tr.ID]; found {
			continue
		}

		tracks, err := to.SearchTrack(ctx, TrackFilters{ID: tr.ID})
		if err != nil {
			return nil, fmt.Errorf("failed to search track %s in destination: %w", tr.ID, err)
		}

		if len(tracks) == 0 {
			cl.Missing = append(cl.Missing, tr)
			slog.Info("track not found in destination, skipping", "isrc", tr.ID)
			continue
		}

		cl.Added = append(cl.Added, tracks[0])
	}

	for _, tr := range dst.Tracks {
		if _, found := srcLookup[tr.ID]; !found {
			cl.Removed = append(cl.Removed, tr)
		}
	}

	return cl, nil
}
