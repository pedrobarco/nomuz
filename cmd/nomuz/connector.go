package main

import (
	"fmt"

	"github.com/pedrobarco/nomuz/internal/domain"
	"github.com/pedrobarco/nomuz/internal/spotify"
)

type ConnectorName string

const (
	ConnectorSpotify ConnectorName = "spotify"
)

func NewConnector(cfg *config, name string) (domain.Connector, error) {
	switch ConnectorName(name) {
	case ConnectorSpotify:
		return spotify.NewConnector(
			cfg.Connectors.Spotify.ClientID,
			cfg.Connectors.Spotify.ClientSecret,
		)
	default:
		return nil, fmt.Errorf("unknown connector: %s", name)
	}
}
