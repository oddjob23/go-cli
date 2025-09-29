package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Repository struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

type Config struct {
	Repositories []Repository `json:"repositories"`
	GitBranch    string       `json:"gitBranch,omitempty"`
}

func LoadFromFile(configFile string) (*Config, error) {
	if configFile == "" {
		configFile = "config.json"
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configFile, err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", configFile, err)
	}

	// Set default branch if not specified
	if config.GitBranch == "" {
		config.GitBranch = "main"
	}

	return &config, nil
}

func (c *Config) Validate() error {
	if len(c.Repositories) == 0 {
		return fmt.Errorf("no repositories configured")
	}

	for i, repo := range c.Repositories {
		if repo.Path == "" {
			return fmt.Errorf("repository %d: path is required", i)
		}
		if repo.Name == "" {
			return fmt.Errorf("repository %d: name is required", i)
		}
		if _, err := os.Stat(repo.Path); os.IsNotExist(err) {
			return fmt.Errorf("repository %s: path %s does not exist", repo.Name, repo.Path)
		}
		if !isDirectory(repo.Path) {
			return fmt.Errorf("repository %s: path %s is not a directory", repo.Name, repo.Path)
		}
		if !isRepository(repo.Path) {
			return fmt.Errorf("repository %s: path %s is not a git repository", repo.Name, repo.Path)
		}
	}

	return nil
}

func isDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func isRepository(path string) bool {
	gitDir := filepath.Join(path, ".git")
	_, err := os.Stat(gitDir)
	return err == nil
}
