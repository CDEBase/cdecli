package main

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type Config struct {
	Registry    string `json:"registry"`
	Endpoint    string `json:"endpoint"`
	AccessToken string `json:"access_token"`
}

func loadConfig(configPath string, flags *ApplicationFlags) (*Config, error) {
	cfgPath := configPath
	userSpecified := cfgPath != ""

	if !userSpecified {
		user, err := user.Current()
		if err != nil {
			return nil, err
		}

		cfgPath = filepath.Join(user.HomeDir, "cde-config.json")
	}

	data, err := os.ReadFile(os.ExpandEnv(cfgPath))

	if err != nil && (!os.IsNotExist(err) || userSpecified) {
		return nil, err
	}

	config := &Config{}
	if err == nil {
		if err := json.Unmarshal(data, config); err != nil {
			return nil, err
		}
	}

	// Apply config overrides.
	if envToken := os.Getenv("SRC_ACCESS_TOKEN"); envToken != "" {
		config.AccessToken = envToken
	}

	// Override endpoint from flags
	if flags.Endpoint != "" {
		config.Endpoint = strings.TrimSuffix(flags.Endpoint, "/")
	}

	// Override access_key from flags
	if flags.AccessToken != "" {
		config.AccessToken = flags.AccessToken
	}

	return config, nil
}
