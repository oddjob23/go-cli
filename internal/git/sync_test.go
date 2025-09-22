package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/oddjob23/go-cli/pkg/utils"
)

func TestSyncer_SyncRepositories(t *testing.T) {
	tests := []struct {
		name               string
		setupDir           func(string) error
		expectedTotalCount int
		expectError        bool
	}{
		{
			name: "sync multiple repositories",
			setupDir: func(tempDir string) error {
				// Create directories with .git folders
				repo1 := filepath.Join(tempDir, "repo1")
				repo2 := filepath.Join(tempDir, "repo2")

				if err := os.MkdirAll(filepath.Join(repo1, ".git"), 0755); err != nil {
					return err
				}
				if err := os.MkdirAll(filepath.Join(repo2, ".git"), 0755); err != nil {
					return err
				}
				return nil
			},
			expectedTotalCount: 2,
			expectError:        false,
		},
		{
			name: "sync empty directory",
			setupDir: func(tempDir string) error {
				return nil // Empty directory
			},
			expectedTotalCount: 0,
			expectError:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "git-sync-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Setup test directory structure
			if err := tt.setupDir(tempDir); err != nil {
				t.Fatalf("Failed to setup test directory: %v", err)
			}

			// Create output handler
			output := utils.NewCliOutput(false)

			// Test the syncer
			syncer := NewSyncer(output)
			result, err := syncer.SyncRepositories(tempDir, "main")

			// Check error expectation
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			// Check result
			if result == nil {
				t.Fatalf("Expected result but got nil")
			}

			if result.TotalRepositories != tt.expectedTotalCount {
				t.Errorf("Expected %d total repositories, got %d", tt.expectedTotalCount, result.TotalRepositories)
			}

			// Verify that success + failure counts equal total
			if result.SuccessCount+result.FailureCount != result.TotalRepositories {
				t.Errorf("Success (%d) + Failure (%d) counts don't equal total (%d)",
					result.SuccessCount, result.FailureCount, result.TotalRepositories)
			}

			// Verify results slice length matches total
			if len(result.Results) != result.TotalRepositories {
				t.Errorf("Results slice length (%d) doesn't match total repositories (%d)",
					len(result.Results), result.TotalRepositories)
			}
		})
	}
}

func TestSyncer_SyncRepositories_NonExistentDirectory(t *testing.T) {
	output := utils.NewCliOutput(false)

	syncer := NewSyncer(output)
	_, err := syncer.SyncRepositories("/non/existent/directory", "main")

	if err == nil {
		t.Errorf("Expected error for non-existent directory")
	}
}
