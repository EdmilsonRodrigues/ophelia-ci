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

func LoadConfig() Config {
	once.Do(func() {
		data, err := os.ReadFile("/etc/ophelia-ci/server-config.toml")
		if err != nil {
			panic(err)
		}

		if err := toml.Unmarshal(data, &configCache); err != nil {
			panic(err)
		}
	})
	return configCache
}
