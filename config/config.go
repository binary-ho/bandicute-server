package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var instance *Config

type Config struct {
	Server struct {
		Port int    `yaml:"port"`
		Env  string `yaml:"env"`
	} `yaml:"server"`

	Database struct {
		BaseURL string `yaml:"base_url"`
		Key     string `yaml:"key"`
	} `yaml:"database"`

	GitHub struct {
		Token string `yaml:"token"`
	} `yaml:"github"`

	SummaryAssistant struct {
		APIKey string `yaml:"api_key"`
	} `yaml:"summary_assistant"`

	Logging struct {
		Level string `yaml:"level"`
	} `yaml:"logging"`

	CORS struct {
		AllowOrigins     string `yaml:"allow_origins"`
		AllowMethods     string `yaml:"allow_methods"`
		AllowHeaders     string `yaml:"allow_headers"`
		AllowCredentials bool   `yaml:"allow_credentials"`
	} `yaml:"cors"`
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
