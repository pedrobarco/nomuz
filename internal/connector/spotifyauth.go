package connector

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const (
	spotifyAuthRedirectURI = "http://127.0.0.1:8080/callback"
	spotifyAuthState       = "state-string"
)

func NewSpotifyAuthServer(auth *spotifyauth.Authenticator, ch chan<- *spotify.Client) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.Token(r.Context(), spotifyAuthState, r)
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusForbidden)
			log.Fatalf("Couldn't get token: %v", err)
		}

		if st := r.FormValue("state"); st != spotifyAuthState {
			http.NotFound(w, r)
			log.Fatalf("State mismatch: %s != %s\n", st, spotifyAuthState)
		}

		client := spotify.New(auth.Client(r.Context(), token))
		ch <- client
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Got request for:", r.URL.String())
	})

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}
