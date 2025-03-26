package main

import (
	"log"
	"os"
	"strconv"

	pb "github.com/EdmilsonRodrigues/ophelia-ci"
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
)

const (
	configPath = "/etc/ophelia-ci/server-config.toml"
	homePath   = "/var/lib/ophelia"
)

// LoadConfig reads the server configuration from a TOML file located at
// configPath. It uses a sync.Once to ensure the configuration is loaded only
// once and caches the result. If reading the file or unmarshalling the TOML data
// fails, the function panics. It returns the cached configuration.
func LoadConfig() Config {
	var err error
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config := loadConfigFromEnv()
		if err := pb.SaveConfig(configPath, config); err != nil {
			panic(err)
		}
		return config
	}

	configCache, err = pb.LoadConfigFromFile(configPath, configCache)
	if err != nil {
		panic(err)
	}

	return configCache
}

// loadConfigFromEnv loads the server configuration from environment variables.
// It retrieves the server address, authentication secret, database path, and
// SSL key file path from the corresponding environment variables and populates
// the Config struct. If the environment variables are not set, it falls back
// to default values.
func loadConfigFromEnv() (config Config) {
	port, err := strconv.Atoi(os.Getenv("OPHELIA_CI_SERVER_PORT"))
	if err != nil || port <= 0 {
		log.Printf("OPHELIA_CI_SERVER_PORT is not set or invalid. Using default port 50051.")
		port = 50051
	}
	config.Server.Port = port

	secret := os.Getenv("OPHELIA_CI_SERVER_SECRET")
	if secret == "" {
		log.Printf("OPHELIA_CI_SERVER_SECRET is not set. Using random secret.")
		secret = randomKey()
	}
	config.Server.Secret = secret

	dbPath := os.Getenv("OPHELIA_CI_SERVER_DB_PATH")
	if dbPath == "" {
		dbPath = homePath + "/ophelia.db"
		log.Printf("OPHELIA_CI_SERVER_DB_PATH is not set. Using default path %s.", dbPath)
	}

	config.Server.DBPath = dbPath
	expirationTime, err := strconv.Atoi(os.Getenv("OPHELIA_CI_SERVER_EXPIRATION_TIME"))
	if err != nil || expirationTime <= 0 {
		log.Printf("OPHELIA_CI_SERVER_EXPIRATION_TIME is not set or invalid. Using default expiration time 30 days.")
		expirationTime = 30
	}
	config.Server.ExpirationTime = expirationTime

	config.SSL.CertFile = os.Getenv("OPHELIA_CI_SERVER_CERT_FILE")
	config.SSL.KeyFile = os.Getenv("OPHELIA_CI_SERVER_KEY_FILE")
	return
}
