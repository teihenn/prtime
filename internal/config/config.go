package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Bitbucket struct {
		AllRepoApiKeyPath string `yaml:"allRepoApiKeyPath"`
		Repositories      []struct {
			Owner string `yaml:"owner"`
			Name  string `yaml:"name"`
		} `yaml:"repositories"`
	} `yaml:"bitbucket"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
