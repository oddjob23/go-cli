package git

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanner_ScanDirectory(t *testing.T) {
	tests := []struct {
		name          string
		setupDir      func(string) error
		expectedCount int
		expectError   bool
	}{
		{
			name: "directory with git repositories",
			setupDir: func(tempDir string) error {
				// Create directories with .git folders
				repo1 := filepath.Join(tempDir, "repo1")
				repo2 := filepath.Join(tempDir, "repo2")
				nonRepo := filepath.Join(tempDir, "notarepo")

				if err := os.MkdirAll(filepath.Join(repo1, ".git"), 0755); err != nil {
					return err
				}
				if err := os.MkdirAll(filepath.Join(repo2, ".git"), 0755); err != nil {
					return err
				}
				if err := os.MkdirAll(nonRepo, 0755); err != nil {
					return err
				}
				return nil
			},
			expectedCount: 2,
			expectError:   false,
		},
		{
			name: "directory with no git repositories",
			setupDir: func(tempDir string) error {
				// Create directories without .git folders
				dir1 := filepath.Join(tempDir, "dir1")
				dir2 := filepath.Join(tempDir, "dir2")

				if err := os.MkdirAll(dir1, 0755); err != nil {
					return err
				}
				if err := os.MkdirAll(dir2, 0755); err != nil {
					return err
				}
				return nil
			},
			expectedCount: 0,
			expectError:   false,
		},
		{
			name: "empty directory",
			setupDir: func(tempDir string) error {
				return nil // Empty directory
			},
			expectedCount: 0,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "git-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Setup test directory structure
			if err := tt.setupDir(tempDir); err != nil {
				t.Fatalf("Failed to setup test directory: %v", err)
			}

			// Test the scanner
			scanner := NewScanner()
			repositories, err := scanner.ScanDirectory(tempDir)

			// Check error expectation
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			// Check repository count
			if len(repositories) != tt.expectedCount {
				t.Errorf("Expected %d repositories, got %d", tt.expectedCount, len(repositories))
			}

			// Verify repository details
			for _, repo := range repositories {
				if repo.Path == "" {
					t.Errorf("Repository path is empty")
				}
				if repo.Name == "" {
					t.Errorf("Repository name is empty")
				}
				// Verify that .git directory exists
				gitPath := filepath.Join(repo.Path, ".git")
				if _, err := os.Stat(gitPath); os.IsNotExist(err) {
					t.Errorf("Expected .git directory at %s", gitPath)
				}
			}
		})
	}
}

func TestScanner_ScanDirectory_NonExistentDirectory(t *testing.T) {
	scanner := NewScanner()
	_, err := scanner.ScanDirectory("/non/existent/directory")

	if err == nil {
		t.Errorf("Expected error for non-existent directory")
	}
}
