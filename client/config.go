package main

import (
	"os"
	"strings"
	"sync"

	"github.com/pelletier/go-toml/v2"
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
	once        sync.Once
)

const configFile = "/etc/ophelia-ci/client-config.toml"

// LoadConfig reads the client configuration from a TOML file located at
// configFile. It uses a sync.Once to ensure the
// configuration is loaded only once and caches the result. If reading the
// file or unmarshalling the TOML data fails, the function panics. It
// returns the cached configuration.
func LoadConfig() Config {
	var err error
	if checkRunningFromImage() {
		return loadConfigFromEnv()
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		config := loadConfigFromEnv()
		SaveConfig(config)
		return config
	}

	configCache, err = loadConfigFromFile()
	if err != nil {
		panic(err)
	}

	return configCache
}

func loadConfigFromFile() (config Config, err error) {
	once.Do(func() {
		var data []byte
		data, err = os.ReadFile(configFile)
		if err != nil {
			return
		}

		if err = toml.Unmarshal(data, &config); err != nil {
			return
		}
	})
	return
}

func loadConfigFromEnv() (config Config) {
	config.Client.Server = os.Getenv("OPHELIA_CI_SERVER")
	config.Client.AuthToken = os.Getenv("OPHELIA_CI_AUTH_TOKEN")
	config.SSL.KeyFile = os.Getenv("OPHELIA_CI_SSL_KEY_FILE")
	return
}

func checkRunningFromImage() bool {
	return os.Getenv("OPHELIA_CI_FROM_IMAGE") != ""
}

// SaveConfig saves the given configuration to a TOML file located at
// configFile. If marshalling the configuration to TOML data or writing the
// file fails, the function panics.
func SaveConfig(config Config) {
	if checkRunningFromImage() {
		return
	}

	data := new(strings.Builder)
	enc := toml.NewEncoder(data)
	if err := enc.Encode(config); err != nil {
		panic(err)
	}

	if err := os.WriteFile(configFile, []byte(data.String()), 0644); err != nil {
		panic(err)
	}
}
