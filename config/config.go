package config

import (
	"os"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Keys map[string]string `yaml:"keys"`
}

var AppConfig Config
func LoadConfig() (error) {

	data, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, &AppConfig); err != nil {
		return err
	}

	return nil
}
