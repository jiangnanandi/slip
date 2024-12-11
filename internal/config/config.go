package config

import (
	"os"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Keys map[string]string `yaml:"keys"`
	Notes NotesConfig `yaml:"notes"`
	Title string `yaml:"title"`
	DataDir string `yaml:"data_dir"`
}

type NotesConfig struct {
	PublishedDir string `yaml:"published_dir"`
	DraftDir string `yaml:"draft_dir"`
	ArchivedDir string `yaml:"archived_dir"`
	DeletedDir string `yaml:"deleted_dir"`
	PrivateDir string `yaml:"private_dir"`
	DefaultAuthor string `yaml:"default_author"`
	DefaultStatus string `yaml:"default_status"`
	DefaultTags []string `yaml:"default_tags"`
}

var AppConfig Config

func LoadConfig() (error) {

	data, err := os.ReadFile("./configs/config.yaml")
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, &AppConfig); err != nil {
		return err
	}
	return nil
}