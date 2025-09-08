package connector

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/pedrobarco/nomuz/internal/domain"
	"github.com/toqueteos/webbrowser"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

func NewSpotifyConnector(clientID, clientSecret string) (*spotifyConnector, error) {
	ctx := context.Background()

	auth := spotifyauth.New(
		spotifyauth.WithClientID(clientID),
		spotifyauth.WithClientSecret(clientSecret),
		spotifyauth.WithRedirectURL(spotifyAuthRedirectURI),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPrivate,
		),
	)

	token, err := GetSpotifyAuthToken()
	if err != nil || IsInvalidSpotifyAuthToken(token) {
		ch := make(chan *oauth2.Token)
		server := NewSpotifyAuthServer(auth, ch)

		go func() {
			if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("failed to start server: %v", err)
			}
		}()
		defer func() {
			if err := server.Close(); err != nil {
				log.Fatalf("failed to close server: %v", err)
			}
		}()

		webbrowser.Open(auth.AuthURL(spotifyAuthState))
		token = <-ch
	}

	client := spotify.New(auth.Client(ctx, token))

	user, err := client.CurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}

	return &spotifyConnector{
		client: client,
		user:   user,
	}, nil
}

type spotifyConnector struct {
	client *spotify.Client
	user   *spotify.PrivateUser
}

var _ domain.Connector = (*spotifyConnector)(nil)

func (s *spotifyConnector) CreatePlaylist(ctx context.Context, name string) (*domain.Playlist, error) {
	pl, err := s.client.CreatePlaylistForUser(ctx, s.user.ID, name, "", false, false)
	if err != nil {
		return nil, fmt.Errorf("failed to create playlist: %w", err)
	}

	return &domain.Playlist{
		ID:   pl.ID.String(),
		Name: pl.Name,
	}, nil
}

func (s *spotifyConnector) GetPlaylistByName(ctx context.Context, name string) (*domain.Playlist, error) {
	res, err := s.client.GetPlaylistsForUser(ctx, s.user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get playlists for user: %w", err)
	}

	var p *domain.Playlist
	for _, pl := range res.Playlists {
		if pl.Name == name {
			p = &domain.Playlist{
				ID:   pl.ID.String(),
				Name: pl.Name,
			}
			break
		}
	}

	if p == nil {
		return nil, nil
	}

	tracks, err := s.getTracksByPlaylistID(ctx, p.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tracks for playlist %s: %w", p.Name, err)
	}

	p.Tracks = tracks
	return p, nil
}

func (s *spotifyConnector) GetPlaylists(ctx context.Context) ([]*domain.Playlist, error) {
	res, err := s.client.GetPlaylistsForUser(ctx, s.user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get playlists for user: %w", err)
	}

	var pls []*domain.Playlist
	for _, pl := range res.Playlists {
		p := &domain.Playlist{
			ID:     pl.ID.String(),
			Name:   pl.Name,
			Tracks: make([]domain.Track, int(pl.Tracks.Total)),
		}
		pls = append(pls, p)
	}

	return pls, nil
}

func (s *spotifyConnector) AddTracksToPlaylist(ctx context.Context, id string, tracks []domain.Track) error {
	var spotifyTracks []spotify.ID
	for _, t := range tracks {
		spotifyTracks = append(spotifyTracks, spotify.ID(t.ID))
	}

	_, err := s.client.AddTracksToPlaylist(ctx, spotify.ID(id), spotifyTracks...)
	if err != nil {
		return fmt.Errorf("failed to add tracks to playlist: %w", err)
	}

	return nil
}

func (s *spotifyConnector) DeleteTracksFromPlaylist(ctx context.Context, id string, tracks []domain.Track) error {
	var spotifyTracks []spotify.ID
	for _, t := range tracks {
		spotifyTracks = append(spotifyTracks, spotify.ID(t.ID))
	}

	_, err := s.client.RemoveTracksFromPlaylist(ctx, spotify.ID(id), spotifyTracks...)
	if err != nil {
		return fmt.Errorf("failed to remove tracks from playlist: %w", err)
	}

	return nil
}

func (s *spotifyConnector) SearchTrack(ctx context.Context, filters domain.TrackFilters) ([]domain.Track, error) {
	res, err := s.client.Search(ctx, filters.ID, spotify.SearchTypeTrack)
	if err != nil {
		return nil, fmt.Errorf("failed to search track: %w", err)
	}

	if res.Tracks == nil {
		return nil, nil
	}

	var tracks []domain.Track
	for _, t := range res.Tracks.Tracks {
		tracks = append(tracks, s.toDomainTrack(t))
	}

	return tracks, nil
}

func (s *spotifyConnector) getTracksByPlaylistID(ctx context.Context, playlistID string) ([]domain.Track, error) {
	res, err := s.client.GetPlaylistItems(ctx, spotify.ID(playlistID))
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist items: %w", err)
	}

	var tracks []domain.Track
	for _, item := range res.Items {
		if item.Track.Track == nil {
			continue
		}
		tracks = append(tracks, s.toDomainTrack(*item.Track.Track))
	}

	return tracks, nil
}

func (s *spotifyConnector) toDomainTrack(t spotify.FullTrack) domain.Track {
	return domain.Track{
		ID:     t.ID.String(),
		Title:  t.Name,
		Artist: t.Artists[0].Name,
		Album:  t.Album.Name,
	}
}
