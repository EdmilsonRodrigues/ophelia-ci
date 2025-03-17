package main

import (
	"os"

	"github.com/pelletier/go-toml/v2"

)

type Config struct {
	Client struct {
		Server	int	`toml:"server"`
	} `toml:"client"`
	SSL struct {
		KeyFile 	string	`toml:"key_file"`
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
