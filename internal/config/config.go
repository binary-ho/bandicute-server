package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	GitHub   GitHubConfig
	OpenAI   OpenAIConfig
	Logging  LoggingConfig
}

type ServerConfig struct {
	Port int
	Env  string
}

type DatabaseConfig struct {
	BaseURL string
	Key     string
}

type GitHubConfig struct {
	Token string
}

type OpenAIConfig struct {
	APIKey string
}

type LoggingConfig struct {
	Level string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		logrus.Warn("Error loading .env file, using environment variables")
	}

	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		return nil, err
	}

	return &Config{
		Server: ServerConfig{
			Port: port,
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			BaseURL: getEnv("SUPABASE_URL", ""),
			Key:     getEnv("SUPABASE_KEY", ""),
		},
		GitHub: GitHubConfig{
			Token: getEnv("GITHUB_TOKEN", ""),
		},
		OpenAI: OpenAIConfig{
			APIKey: getEnv("OPENAI_API_KEY", ""),
		},
		Logging: LoggingConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
