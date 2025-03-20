package main

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Server struct {
		Port   int    `toml:"port"`
		Secret string `toml:"secret"`
		DBPath string `toml:"db_path"`
	} `toml:"server"`
	SSL struct {
		CertFile string `toml:"cert_file"`
		KeyFile  string `toml:"key_file"`
	} `toml:"ssl"`
}

func LoadConfig() (config Config) {
	data, err := os.ReadFile("/etc/ophelia-ci/server-config.toml")
	if err != nil {
		panic(err)
	}

	if err := toml.Unmarshal(data, &config); err != nil {
		panic(err)
	}

	return
}
