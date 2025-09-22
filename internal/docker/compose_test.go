package docker

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewComposeConfig(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		want     *ComposeConfig
	}{
		{
			name:     "creates new compose config with defaults",
			filePath: "docker-compose.yml",
			want: &ComposeConfig{
				FilePath:    "docker-compose.yml",
				ProjectName: "microservices",
				Timeout:     5 * time.Minute,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewComposeConfig(tt.filePath)
			if got.FilePath != tt.want.FilePath {
				t.Errorf("NewComposeConfig().FilePath = %v, want %v", got.FilePath, tt.want.FilePath)
			}
			if got.ProjectName != tt.want.ProjectName {
				t.Errorf("NewComposeConfig().ProjectName = %v, want %v", got.ProjectName, tt.want.ProjectName)
			}
			if got.Timeout != tt.want.Timeout {
				t.Errorf("NewComposeConfig().Timeout = %v, want %v", got.Timeout, tt.want.Timeout)
			}
		})
	}
}

func TestComposeConfig_ValidateFile(t *testing.T) {
	tempDir := t.TempDir()
	existingFile := filepath.Join(tempDir, "existing.yml")
	nonExistingFile := filepath.Join(tempDir, "non-existing.yml")

	if err := os.WriteFile(existingFile, []byte("version: '3.8'\nservices: {}"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name     string
		filePath string
		wantErr  bool
	}{
		{
			name:     "validates existing file",
			filePath: existingFile,
			wantErr:  false,
		},
		{
			name:     "returns error for non-existing file",
			filePath: nonExistingFile,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ComposeConfig{
				FilePath: tt.filePath,
			}
			err := c.ValidateFile()
			if (err != nil) != tt.wantErr {
				t.Errorf("ComposeConfig.ValidateFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}