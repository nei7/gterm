package term

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Window struct {
		Background string `yaml:"background"`
		Padding    struct {
			X float64 `yaml:"x"`
			Y float64 `yaml:"y"`
		}
	}

	Font struct {
		Family string  `yaml:"family"`
		Size   float64 `yaml:"size"`
	} `yaml:"font"`
}

func loadConfig() (*Config, error) {
	configPath, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	configPath = path.Join(configPath, "gterm", "config.yaml")

	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("failed to read config from %s: %w", configPath, err)
	}

	return &config, nil
}
