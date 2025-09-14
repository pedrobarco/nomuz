package main

import (
	"fmt"

	"github.com/pedrobarco/nomuz/internal/domain"
	"github.com/pedrobarco/nomuz/internal/spotify"
	"github.com/pedrobarco/nomuz/internal/tidal"
)

type ConnectorName string

const (
	ConnectorSpotify ConnectorName = "spotify"
	ConnectorTidal   ConnectorName = "tidal"
)

func NewConnector(cfg *config, name string) (domain.Connector, error) {
	switch ConnectorName(name) {
	case ConnectorSpotify:
		return spotify.NewConnector(
			cfg.Connectors.Spotify.ClientID,
			cfg.Connectors.Spotify.ClientSecret,
		)
	case ConnectorTidal:
		return tidal.NewConnector(
			cfg.Connectors.Tidal.ClientID,
			cfg.Connectors.Tidal.ClientSecret,
			cfg.Connectors.Tidal.CountryCode,
		)
	default:
		return nil, fmt.Errorf("unknown connector: %s", name)
	}
}
