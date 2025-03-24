package main

import (
	"os"
	"sync"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Client struct {
		Server string `toml:"server"`
	} `toml:"client"`
	SSL struct {
		KeyFile string `toml:"key_file"`
	} `toml:"ssl"`
}

var (
	configCache Config
	once        sync.Once
)

// LoadConfig reads the client configuration from a TOML file located at
// "/etc/ophelia-ci/client-config.toml". It uses a sync.Once to ensure the
// configuration is loaded only once and caches the result. If reading the
// file or unmarshalling the TOML data fails, the function panics. It
// returns the cached configuration.
func LoadConfig() Config {
	once.Do(func() {
		data, err := os.ReadFile("/etc/ophelia-ci/client-config.toml")
		if err != nil {
			panic(err)
		}

		if err := toml.Unmarshal(data, &configCache); err != nil {
			panic(err)
		}
	})
	return configCache
}
