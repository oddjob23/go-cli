package docker

import (
	"path/filepath"
	"testing"
)

func TestNewManager(t *testing.T) {
	tests := []struct {
		name    string
		baseDir string
		want    *Manager
	}{
		{
			name:    "creates new manager with correct file paths",
			baseDir: "/test/dir",
			want: &Manager{
				dependenciesFile: "/test/dir/docker-compose.dependencies.yml",
				servicesFile:     "/test/dir/docker-compose.services.yml",
			},
		},
		{
			name:    "handles relative paths",
			baseDir: ".",
			want: &Manager{
				dependenciesFile: "./docker-compose.dependencies.yml",
				servicesFile:     "./docker-compose.services.yml",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewManager(tt.baseDir)

			expectedDepsFile := filepath.Join(tt.baseDir, "docker-compose.dependencies.yml")
			expectedServicesFile := filepath.Join(tt.baseDir, "docker-compose.services.yml")

			if got.dependenciesFile != expectedDepsFile {
				t.Errorf("NewManager().dependenciesFile = %v, want %v", got.dependenciesFile, expectedDepsFile)
			}
			if got.servicesFile != expectedServicesFile {
				t.Errorf("NewManager().servicesFile = %v, want %v", got.servicesFile, expectedServicesFile)
			}
		})
	}
}

func TestManager_FilePathGeneration(t *testing.T) {
	baseDir := "/home/user/project"
	manager := NewManager(baseDir)

	expectedDepsFile := "/home/user/project/docker-compose.dependencies.yml"
	expectedServicesFile := "/home/user/project/docker-compose.services.yml"

	if manager.dependenciesFile != expectedDepsFile {
		t.Errorf("Manager.dependenciesFile = %v, want %v", manager.dependenciesFile, expectedDepsFile)
	}

	if manager.servicesFile != expectedServicesFile {
		t.Errorf("Manager.servicesFile = %v, want %v", manager.servicesFile, expectedServicesFile)
	}
}