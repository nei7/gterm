package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Background string `yaml:"background"`
	Font       Font   `yaml:"font"`
}

type Font struct {
	Family string  `yaml:"family"`
	Size   float64 `yaml:"size"`
}

func LoadConfig() (*Config, error) {
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
