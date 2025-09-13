package tidal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pedrobarco/nomuz/internal/domain"
	"github.com/pedrobarco/nomuz/pkg/tidal"
	"golang.org/x/oauth2/clientcredentials"
)

const serverURL = "https://openapi.tidal.com/v2"
const tokenURL = "https://auth.tidal.com/v1/oauth2/token"

func NewConnector(clientID, clientSecret, countryCode string) (*connector, error) {
	cfg := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
	}
	httpClient := cfg.Client(context.Background())
	client, err := tidal.NewClientWithResponses(
		serverURL,
		tidal.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create tidal connector: %w", err)
	}

	return &connector{
		client:      client,
		countryCode: countryCode,
	}, nil
}

type connector struct {
	client      tidal.ClientWithResponsesInterface
	countryCode string
}

var _ domain.Connector = (*connector)(nil)

func (c *connector) AddTracksToPlaylist(ctx context.Context, id string, tracks []domain.Track) error {
	var data []tidal.PlaylistItemsRelationshipAddOperationPayloadData

	for _, track := range tracks {
		data = append(data, tidal.PlaylistItemsRelationshipAddOperationPayloadData{
			Id:   track.ID,
			Type: "tracks",
		})
	}

	resp, err := c.client.PostPlaylistsIdRelationshipsItemsWithApplicationVndAPIPlusJSONBodyWithResponse(
		ctx,
		id,
		&tidal.PostPlaylistsIdRelationshipsItemsParams{
			CountryCode: c.countryCode,
		},
		tidal.PostPlaylistsIdRelationshipsItemsApplicationVndAPIPlusJSONRequestBody{
			Data: data,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to add tracks to playlist: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("failed to add tracks to playlist: status code %d", resp.StatusCode())
	}

	return nil
}

func (c *connector) CreatePlaylist(ctx context.Context, name string) (*domain.Playlist, error) {
	panic("unimplemented")
}

func (c *connector) DeleteTracksFromPlaylist(ctx context.Context, id string, tracks []domain.Track) error {
	panic("unimplemented")
}

func (c *connector) GetPlaylistByName(ctx context.Context, name string) (*domain.Playlist, error) {
	panic("unimplemented")
}

func (c *connector) GetPlaylists(ctx context.Context) ([]*domain.Playlist, error) {
	resp, err := c.client.GetPlaylistsWithResponse(
		ctx,
		&tidal.GetPlaylistsParams{
			CountryCode: c.countryCode,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get playlists: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to get playlists: status code %d: %s", resp.StatusCode(), string(resp.Body))
	}

	var playlists []*domain.Playlist
	for _, p := range resp.ApplicationvndApiJSON200.Data {
		playlists = append(playlists, &domain.Playlist{
			ID:     p.Id,
			Name:   p.Attributes.Name,
			Tracks: make([]domain.Track, *p.Attributes.NumberOfItems),
		})
	}

	return playlists, nil
}

func (c *connector) SearchTrack(ctx context.Context, filters domain.TrackFilters) ([]domain.Track, error) {
	panic("unimplemented")
}
