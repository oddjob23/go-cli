package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFromFile(t *testing.T) {
	tests := []struct {
		name          string
		configContent string
		configFile    string
		wantErr       bool
		wantBranch    string
		wantRepoCount int
	}{
		{
			name: "should load valid config with custom branch",
			configContent: `{
				"repositories": [
					{"path": "/path/to/repo1", "name": "repo1"},
					{"path": "/path/to/repo2", "name": "repo2"}
				],
				"gitBranch": "develop"
			}`,
			wantErr:       false,
			wantBranch:    "develop",
			wantRepoCount: 2,
		},
		{
			name: "should default to main branch when not specified",
			configContent: `{
				"repositories": [
					{"path": "/path/to/repo1", "name": "repo1"}
				]
			}`,
			wantErr:       false,
			wantBranch:    "main",
			wantRepoCount: 1,
		},
		{
			name: "should load config with empty repositories list",
			configContent: `{
				"repositories": []
			}`,
			wantErr:       false,
			wantBranch:    "main",
			wantRepoCount: 0,
		},
		{
			name:          "should return error when json is invalid",
			configContent: `{"repositories": [}`,
			wantErr:       true,
		},
		{
			name:          "should return error when file does not exist",
			configFile:    "nonexistent.json",
			configContent: "",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file for test
			var configPath string
			if tt.configFile == "nonexistent.json" {
				configPath = tt.configFile
			} else {
				tmpDir := t.TempDir()
				configPath = filepath.Join(tmpDir, "config.json")
				if tt.configContent != "" {
					if err := os.WriteFile(configPath, []byte(tt.configContent), 0644); err != nil {
						t.Fatalf("failed to create test config file: %v", err)
					}
				}
			}

			// Test LoadFromFile
			config, err := LoadFromFile(configPath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("LoadFromFile() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("LoadFromFile() unexpected error: %v", err)
				return
			}

			if config.GitBranch != tt.wantBranch {
				t.Errorf("LoadFromFile() GitBranch = %v, want %v", config.GitBranch, tt.wantBranch)
			}

			if len(config.Repositories) != tt.wantRepoCount {
				t.Errorf("LoadFromFile() repository count = %v, want %v", len(config.Repositories), tt.wantRepoCount)
			}
		})
	}
}

func TestLoadFromFileDefaultPath(t *testing.T) {
	// Save current directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Create temp directory and change to it
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	// Create config.json in temp directory
	configContent := `{
		"repositories": [
			{"path": "/path/to/repo", "name": "repo"}
		]
	}`
	if err := os.WriteFile("config.json", []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to create config.json: %v", err)
	}

	// Test with empty string (should default to config.json)
	config, err := LoadFromFile("")
	if err != nil {
		t.Errorf("LoadFromFile(\"\") unexpected error: %v", err)
		return
	}

	if len(config.Repositories) != 1 {
		t.Errorf("LoadFromFile(\"\") repository count = %v, want 1", len(config.Repositories))
	}
}

func TestValidate(t *testing.T) {
	// Create a temporary git repository for testing
	tmpDir := t.TempDir()
	gitRepo := filepath.Join(tmpDir, "valid-repo")
	if err := os.MkdirAll(filepath.Join(gitRepo, ".git"), 0755); err != nil {
		t.Fatalf("failed to create test git repo: %v", err)
	}

	// Create a non-git directory
	nonGitDir := filepath.Join(tmpDir, "non-git")
	if err := os.MkdirAll(nonGitDir, 0755); err != nil {
		t.Fatalf("failed to create non-git directory: %v", err)
	}

	// Create a file (not a directory)
	testFile := filepath.Join(tmpDir, "file.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	tests := []struct {
		name    string
		config  *Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "should validate successfully when config is valid",
			config: &Config{
				Repositories: []Repository{
					{Path: gitRepo, Name: "valid-repo"},
				},
				GitBranch: "main",
			},
			wantErr: false,
		},
		{
			name: "should return error when no repositories configured",
			config: &Config{
				Repositories: []Repository{},
				GitBranch:    "main",
			},
			wantErr: true,
			errMsg:  "no repositories configured",
		},
		{
			name: "should return error when repository path is missing",
			config: &Config{
				Repositories: []Repository{
					{Path: "", Name: "repo"},
				},
				GitBranch: "main",
			},
			wantErr: true,
			errMsg:  "path is required",
		},
		{
			name: "should return error when repository name is missing",
			config: &Config{
				Repositories: []Repository{
					{Path: gitRepo, Name: ""},
				},
				GitBranch: "main",
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "should return error when repository path does not exist",
			config: &Config{
				Repositories: []Repository{
					{Path: "/nonexistent/path", Name: "repo"},
				},
				GitBranch: "main",
			},
			wantErr: true,
			errMsg:  "does not exist",
		},
		{
			name: "should return error when repository path is not a directory",
			config: &Config{
				Repositories: []Repository{
					{Path: testFile, Name: "repo"},
				},
				GitBranch: "main",
			},
			wantErr: true,
			errMsg:  "is not a directory",
		},
		{
			name: "should return error when repository path is not a git repository",
			config: &Config{
				Repositories: []Repository{
					{Path: nonGitDir, Name: "repo"},
				},
				GitBranch: "main",
			},
			wantErr: true,
			errMsg:  "is not a git repository",
		},
		{
			name: "should return error when one of multiple repositories is invalid",
			config: &Config{
				Repositories: []Repository{
					{Path: gitRepo, Name: "valid-repo"},
					{Path: nonGitDir, Name: "invalid-repo"},
				},
				GitBranch: "main",
			},
			wantErr: true,
			errMsg:  "is not a git repository",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() expected error containing %q, got nil", tt.errMsg)
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("Validate() error = %q, want error containing %q", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestIsDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a directory
	dir := filepath.Join(tmpDir, "testdir")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}

	// Create a file
	file := filepath.Join(tmpDir, "testfile.txt")
	if err := os.WriteFile(file, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "should return true when path is a valid directory",
			path: dir,
			want: true,
		},
		{
			name: "should return false when path is a file",
			path: file,
			want: false,
		},
		{
			name: "should return false when path does not exist",
			path: "/nonexistent/path",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isDirectory(tt.path)
			if got != tt.want {
				t.Errorf("isDirectory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsRepository(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a git repository
	gitRepo := filepath.Join(tmpDir, "git-repo")
	if err := os.MkdirAll(filepath.Join(gitRepo, ".git"), 0755); err != nil {
		t.Fatalf("failed to create test git repo: %v", err)
	}

	// Create a non-git directory
	nonGitDir := filepath.Join(tmpDir, "non-git")
	if err := os.MkdirAll(nonGitDir, 0755); err != nil {
		t.Fatalf("failed to create non-git directory: %v", err)
	}

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "should return true when path is a valid git repository",
			path: gitRepo,
			want: true,
		},
		{
			name: "should return false when path is not a git repository",
			path: nonGitDir,
			want: false,
		},
		{
			name: "should return false when path does not exist",
			path: "/nonexistent/path",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isRepository(tt.path)
			if got != tt.want {
				t.Errorf("isRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && indexOf(s, substr) >= 0))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}