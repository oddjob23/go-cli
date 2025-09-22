package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// Git sync configuration
	ScanDirectory string
	GitBranch     string
}

func New() *Config {
	return &Config{
		ScanDirectory: getEnv("SCAN_DIRECTORY", "./"),
		GitBranch:     getEnv("GIT_BRANCH", "main"),
	}
}

func LoadFromFile(envFile string) (*Config, error) {
	if envFile != "" {
		if err := godotenv.Load(envFile); err != nil {
			return nil, err
		}
	} else {
		// Try to load .env from current directory if no file specified
		if _, err := os.Stat(".env"); err == nil {
			_ = godotenv.Load(".env")
		}
	}

	return New(), nil
}

func (c *Config) Validate() error {
	if c.ScanDirectory != "" {
		if _, err := os.Stat(c.ScanDirectory); os.IsNotExist(err) {
			return err
		}
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
