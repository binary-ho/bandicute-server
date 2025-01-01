package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var instance *Config

type Config struct {
	Server struct {
		Port    int    `yaml:"port"`
		Env     string `yaml:"env"`
		Address string `yaml:"address"`
	} `yaml:"server"`

	Database struct {
		BaseURL string `yaml:"base_url"`
		Key     string `yaml:"key"`
	} `yaml:"database"`

	GitHub struct {
		Token string `yaml:"token"`
	} `yaml:"github"`

	OpenAI struct {
		APIKey string `yaml:"api_key"`
	} `yaml:"openai"`

	Logging struct {
		Level string `yaml:"level"`
	} `yaml:"logging"`
}

const (
	configPath = "config/property.yml"
)

func GetInstance() *Config {
	if instance == nil {
		return loadConfig()
	}
	return instance
}

func loadConfig() *Config {
	instance = &Config{}
	if err := loadYAML(configPath, instance); err != nil {
		panic(err)
	}
	return instance
}

func loadYAML(filename string, cfg *Config) error {
	file, err := os.ReadFile(filepath.Clean(filename))
	if err != nil {
		return err
	}
	return yaml.Unmarshal(file, cfg)
}
