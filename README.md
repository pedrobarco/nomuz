# Nomuz 🎶

Nomuz is a CLI tool for transferring/syncing playlists between different music streaming platforms.
It connects to your accounts, fetches playlists, matches tracks across services, and generates a changelog of what was added, removed, or missing.

## Features

- Transfer playlists between supported platforms (Spotify, YouTube Music, Apple Music, Deezer, Tidal, …).
- Track matching by ISRC (preferred) or metadata fallback (title + artist).
- Generate changelogs with Added, Removed, and Missing tracks.
- Modular connector system → easy to add new platforms.

## Installation

Build from source (Go 1.21+ required):

```sh
git clone https://github.com/pedrobarco/nomuz.git
cd nomuz
go build -o nomuz ./cmd/nomuz
```

Now you can run it as:

```sh
./nomuz --help
```

## Usage

### List playlists

```sh
nomuz playlists --from spotify
```

### Transfer a playlist

```sh
nomuz sync --from spotify --to ytmusic --playlist "My Favorites"
```

### Example Output (changelog)

```
Added:   45 tracks
Removed: 3 tracks
Missing: 2 tracks
```
