package config

import (
	"os"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Keys map[string]string `yaml:"keys"`
	Notes NotesConfig `yaml:"notes"`
}

type NotesConfig struct {
	Dir string `yaml:"dir"`
	DefaultAuthor string `yaml:"default_author"`
	DefaultStatus string `yaml:"default_status"`
	DefaultTags []string `yaml:"default_tags"`
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
