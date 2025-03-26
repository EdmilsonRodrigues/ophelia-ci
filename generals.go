package ophelia_ci

import (
	"os"
	"strings"
	"sync"

	"github.com/pelletier/go-toml/v2"
)

var (
	once sync.Once
)

func CheckRunningFromImage() bool {
	return os.Getenv("OPHELIA_CI_FROM_IMAGE") != ""
}

func LoadConfigFromFile[CT any](configFile string, config CT) (CT, error) {
	var err error
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
	return config, err
}

func SaveConfig[CT any](configFile string, config CT) (err error) {
	if CheckRunningFromImage() {
		return
	}

	data := new(strings.Builder)
	enc := toml.NewEncoder(data)
	if err = enc.Encode(config); err != nil {
		return
	}

	if err = os.WriteFile(configFile, []byte(data.String()), 0644); err != nil {
		return
	}
	return
}

