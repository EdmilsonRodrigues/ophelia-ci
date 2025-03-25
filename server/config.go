package main

import (
	"os"
	"sync"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Server struct {
		Port           int    `toml:"port"`
		Secret         string `toml:"secret"`
		DBPath         string `toml:"db_path"`
		ExpirationTime int    `toml:"expiration_time"`
	} `toml:"server"`
	SSL struct {
		CertFile string `toml:"cert_file"`
		KeyFile  string `toml:"key_file"`
	} `toml:"ssl"`
}

var (
	configCache Config
	once        sync.Once
)

const (
	configPath = "/etc/ophelia-ci/server-config.toml"
	homePath   = "/var/lib/ophelia"
)
// LoadConfig reads the server configuration from a TOML file located at
// the configPath. It uses a sync.Once to ensure the
// configuration is loaded only once and caches the result. If reading the
// file or unmarshalling the TOML data fails, the function panics. It 
// returns the cached configuration.
func LoadConfig() Config {
	once.Do(func() {
		data, err := os.ReadFile(configPath)
		if err != nil {
			panic(err)
		}

		if err := toml.Unmarshal(data, &configCache); err != nil {
			panic(err)
		}
	})
	
	return configCache
}
