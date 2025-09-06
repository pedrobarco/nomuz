package domain

type Playlist struct {
	ID     string
	Name   string
	Tracks []Track
}

type Track struct {
	ID     string
	Title  string
	Artist string
	Album  string
}
