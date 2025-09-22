package docker

import (
	"os/exec"
	"testing"
)

func TestGetDockerComposeCommand(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "should return docker compose command",
			expected: "docker compose",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetDockerComposeCommand()
			if result != tt.expected && result != "docker-compose" {
				t.Errorf("GetDockerComposeCommand() = %v, want %v or docker-compose", result, tt.expected)
			}
		})
	}
}

func TestCheckDockerCompose(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		setup   func()
	}{
		{
			name:    "should check docker compose availability",
			wantErr: false,
			setup:   func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			err := CheckDockerCompose()

			if _, lookErr := exec.LookPath("docker"); lookErr != nil {
				if err == nil {
					t.Error("Expected error when docker is not available, but got nil")
				}
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("CheckDockerCompose() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}