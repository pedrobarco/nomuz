package connector

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

const (
	spotifyAuthRedirectURI = "http://127.0.0.1:8080/callback"
	spotifyAuthState       = "state-string"
	successHTML            = `
	<html>
	<body>
		<h2>Authentication Successful!</h2>
		<p>You can now close this window and return to your terminal.</p>
		<script>window.close();</script>
	</body>
	</html>
	`
)

func NewSpotifyAuthServer(auth *spotifyauth.Authenticator, ch chan<- *oauth2.Token) *http.Server {
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

		defer func() {
			if err := SaveSpotifyAuthToken(token); err != nil {
				log.Printf("Failed to save spotify auth token: %v", err)
			}
		}()

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, successHTML)

		ch <- token
	})

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}

func getSpotifyAuthConfigPath() (string, error) {
	cfg, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home dir: %v", err)
	}
	return path.Join(cfg, "nomuz", "spotify_auth.yaml"), nil
}

type spotifyAuthConfig struct {
	Token *oauth2.Token `yaml:"token"`
}

func SaveSpotifyAuthToken(token *oauth2.Token) error {
	filePath, err := getSpotifyAuthConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get spotify auth config path: %v", err)
	}

	_, err = os.Stat(filePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat spotify auth config file: %v", err)
	}

	if os.IsNotExist(err) {
		if err := os.MkdirAll(path.Dir(filePath), 00755); err != nil {
			return fmt.Errorf("failed to create spotify auth config dir: %v", err)
		}
	}

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create spotify auth config file: %v", err)
	}
	defer f.Close()

	cfg := &spotifyAuthConfig{Token: token}
	if err := yaml.NewEncoder(f).Encode(cfg); err != nil {
		return fmt.Errorf("failed to write spotify auth config file: %v", err)
	}

	return nil
}

func GetSpotifyAuthToken() (*oauth2.Token, error) {
	path, err := getSpotifyAuthConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get spotify auth config path: %v", err)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open spotify auth config file: %v", err)
	}
	defer f.Close()

	var cfg spotifyAuthConfig
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse spotify auth config file: %v", err)
	}

	return cfg.Token, nil
}

func IsInvalidSpotifyAuthToken(token *oauth2.Token) bool {
	return token == nil || !token.Valid()
}
