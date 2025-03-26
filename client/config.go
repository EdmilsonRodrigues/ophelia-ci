package main

import (
	"os"

	pb "github.com/EdmilsonRodrigues/ophelia-ci"
)

type Config struct {
	Client struct {
		Server    string `toml:"server"`
		AuthToken string `toml:"auth_token"`
	} `toml:"client"`
	SSL struct {
		KeyFile string `toml:"key_file"`
	} `toml:"ssl"`
}

var (
	configCache Config
)

const configFile = "/etc/ophelia-ci/client-config.toml"

// LoadConfig reads the client configuration from a TOML file located at
// configFile. It uses a sync.Once to ensure the
// configuration is loaded only once and caches the result. If reading the
// file or unmarshalling the TOML data fails, the function panics. It
// returns the cached configuration.
func LoadConfig() Config {
	var err error
	if pb.CheckRunningFromImage() {
		return loadConfigFromEnv()
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		config := loadConfigFromEnv()
		if err := pb.SaveConfig(configFile, config); err != nil {
			panic(err)
		}
		return config
	}

	configCache, err = pb.LoadConfigFromFile(configFile, configCache)
	if err != nil {
		panic(err)
	}

	return configCache
}

// loadConfigFromEnv loads the client configuration from environment variables.
// It retrieves the server address, authentication token, and SSL key file path
// from the corresponding environment variables and populates the Config struct.
func loadConfigFromEnv() (config Config) {
	server := os.Getenv("OPHELIA_CI_SERVER")
	if server == "" {
		server = "localhost:50051"
	}
	config.Client.Server = server
	config.Client.AuthToken = os.Getenv("OPHELIA_CI_AUTH_TOKEN")
	config.SSL.KeyFile = os.Getenv("OPHELIA_CI_SSL_KEY_FILE")
	return
}
