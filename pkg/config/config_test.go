package config

import (
	"os"
	"testing"
)

func TestConfig_New(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected Config
	}{
		{
			name:    "default values",
			envVars: map[string]string{},
			expected: Config{
				ScanDirectory: "./",
				GitBranch:     "main",
			},
		},
		{
			name: "custom env values",
			envVars: map[string]string{
				"SCAN_DIRECTORY": "/test/path",
				"GIT_BRANCH":     "develop",
			},
			expected: Config{
				ScanDirectory: "/test/path",
				GitBranch:     "develop",
			},
		},
		{
			name: "partial env values",
			envVars: map[string]string{
				"SCAN_DIRECTORY": "/custom/path",
			},
			expected: Config{
				ScanDirectory: "/custom/path",
				GitBranch:     "main",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			cfg := New()

			if cfg.ScanDirectory != tt.expected.ScanDirectory {
				t.Errorf("Expected ScanDirectory %s, got %s", tt.expected.ScanDirectory, cfg.ScanDirectory)
			}

			if cfg.GitBranch != tt.expected.GitBranch {
				t.Errorf("Expected GitBranch %s, got %s", tt.expected.GitBranch, cfg.GitBranch)
			}

			// Clean up environment variables
			for key := range tt.envVars {
				os.Unsetenv(key)
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name      string
		config    Config
		wantError bool
	}{
		{
			name: "valid config with current directory",
			config: Config{
				ScanDirectory: "./",
				GitBranch:     "main",
			},
			wantError: false,
		},
		{
			name: "valid config with empty directory",
			config: Config{
				ScanDirectory: "",
				GitBranch:     "main",
			},
			wantError: false,
		},
		{
			name: "invalid config with non-existent directory",
			config: Config{
				ScanDirectory: "/non/existent/path",
				GitBranch:     "main",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.wantError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}
