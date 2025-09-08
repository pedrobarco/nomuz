package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
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

		res, err := connector.GetPlaylists(ctx)
		if err != nil {
			return fmt.Errorf("failed to get playlists: %w", err)
		}

		var pls []*domain.Playlist
		name := cmd.String("name")
		for _, p := range res {
			if name == "" || p.Name == name {
				pls = append(pls, p)
			}
		}

		t := table.New().
			Border(lipgloss.NormalBorder()).
			StyleFunc(func(row, col int) lipgloss.Style {
				return cellStyle
			})

		t.Headers("ID", "Name", "# Tracks")
		for _, p := range pls {
			t.Row(p.ID, p.Name, strconv.Itoa(len(p.Tracks)))
		}

		fmt.Println(t.Render())

		return nil
	},
}
