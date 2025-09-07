package main

import (
	"fmt"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

type config struct {
	Connectors connectorsConfig `yaml:"connectors"`
}

type connectorsConfig struct {
	Spotify spotifyConfig `yaml:"spotify"`
}

type spotifyConfig struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

var defaultConfig = config{}

func LoadConfig() (*config, error) {
	config := defaultConfig

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user config dir: %v", err)
	}

	cfgFile := path.Join(home, ".config", "nomuz", "config.yaml")
	if _, err := os.Stat(cfgFile); err != nil && os.IsNotExist(err) {
		if err := os.MkdirAll(path.Dir(cfgFile), 00755); err != nil {
			return nil, fmt.Errorf("failed to create config dir: %v", err)
		}

		f, err := os.Create(cfgFile)
		if err != nil {
			return nil, fmt.Errorf("failed to create config file: %v", err)
		}
		defer f.Close()

		if err := yaml.NewEncoder(f).Encode(&config); err != nil {
			return nil, fmt.Errorf("failed to write default config: %v", err)
		}
	} else {
		f, err := os.Open(cfgFile)
		if err != nil {
			return nil, fmt.Errorf("failed to open config file: %v", err)
		}
		defer f.Close()

		if err := yaml.NewDecoder(f).Decode(&config); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %v", err)
		}
	}

	return &config, nil
}
