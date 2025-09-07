package main

import (
	"context"
	"fmt"

	"github.com/pedrobarco/nomuz/internal/domain"
	"github.com/urfave/cli/v3"
)

var playlistsCmd = &cli.Command{
	Name:      "playlists",
	Usage:     "List all available playlists",
	UsageText: `nomuz playlists --from <connector> [--name <playlist name>]`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "from",
			Usage:    "Source connector",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "name",
			Usage: "Playlist name",
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		cfg, err := LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		connector, err := NewConnector(cfg, cmd.String("from"))
		if err != nil {
			return fmt.Errorf("failed to create connector: %w", err)
		}

		var res []*domain.Playlist

		name := cmd.String("name")
		if name != "" {
			playlist, err := connector.GetPlaylistByName(ctx, name)
			if err != nil {
				return fmt.Errorf("failed to get playlist by name: %w", err)
			}
			res = append(res, playlist)
		} else {
			playlists, err := connector.GetPlaylists(ctx)
			if err != nil {
				return fmt.Errorf("failed to get playlists: %w", err)
			}
			res = playlists
		}

		for _, p := range res {
			fmt.Printf("ID : %s\n", p.ID)
			fmt.Printf("Name : %s\n", p.Name)
			fmt.Printf("Number of tracks : %d\n", len(p.Tracks))
			fmt.Println("-----")
		}

		return nil
	},
}
