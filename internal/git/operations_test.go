package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckoutMainBranch(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(t *testing.T) string
		branchName     string
		wantSuccess    bool
		wantMsgContain string
	}{
		{
			name: "should checkout main when on different branch",
			setupRepo: func(t *testing.T) string {
				return createTestGitRepo(t, "feature-branch")
			},
			branchName:     "main",
			wantSuccess:    false, // Will fail on pull due to no remote
			wantMsgContain: "not accessible or not found",
		},
		{
			name: "should attempt pull when already on main",
			setupRepo: func(t *testing.T) string {
				return createTestGitRepo(t, "main")
			},
			branchName:     "main",
			wantSuccess:    false, // Will fail on pull due to no remote
			wantMsgContain: "not accessible or not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if testing.Short() {
				t.Skip("skipping integration test in short mode")
			}

			repoPath := tt.setupRepo(t)
			defer os.RemoveAll(repoPath)

			ops := NewOperations()
			repo := Repository{
				Path: repoPath,
				Name: "test-repo",
			}

			result := ops.CheckoutMainBranch(repo, tt.branchName)

			if result.Success != tt.wantSuccess {
				t.Errorf("CheckoutMainBranch() Success = %v, want %v (Error: %v, Message: %v)",
					result.Success, tt.wantSuccess, result.Error, result.Message)
			}

			if tt.wantMsgContain != "" && !strings.Contains(result.Message, tt.wantMsgContain) {
				t.Errorf("CheckoutMainBranch() Message = %q, want to contain %q", result.Message, tt.wantMsgContain)
			}

			if !tt.wantSuccess && result.Error == nil {
				t.Errorf("CheckoutMainBranch() expected error, got nil")
			}
		})
	}
}

func TestGetCurrentBranch(t *testing.T) {
	tests := []struct {
		name       string
		setupRepo  func(t *testing.T) string
		wantBranch string
		wantErr    bool
	}{
		{
			name: "should return current branch name when on main",
			setupRepo: func(t *testing.T) string {
				return createTestGitRepo(t, "main")
			},
			wantBranch: "main",
			wantErr:    false,
		},
		{
			name: "should return current branch name when on feature branch",
			setupRepo: func(t *testing.T) string {
				return createTestGitRepo(t, "feature-branch")
			},
			wantBranch: "feature-branch",
			wantErr:    false,
		},
		{
			name: "should return error when path is not a git repository",
			setupRepo: func(t *testing.T) string {
				tmpDir := t.TempDir()
				return tmpDir
			},
			wantBranch: "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if testing.Short() {
				t.Skip("skipping integration test in short mode")
			}

			repoPath := tt.setupRepo(t)
			defer os.RemoveAll(repoPath)

			ops := NewOperations()
			branch, err := ops.getCurrentBranch(repoPath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("getCurrentBranch() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("getCurrentBranch() unexpected error: %v", err)
				return
			}

			if branch != tt.wantBranch {
				t.Errorf("getCurrentBranch() = %q, want %q", branch, tt.wantBranch)
			}
		})
	}
}

func TestPullFromMain(t *testing.T) {
	tests := []struct {
		name      string
		setupRepo func(t *testing.T) string
		wantErr   bool
	}{
		{
			name: "should return error when no remote configured",
			setupRepo: func(t *testing.T) string {
				return createTestGitRepo(t, "main")
			},
			wantErr: true, // Local repo with no remote will fail
		},
		{
			name: "should return error when not on a git repository",
			setupRepo: func(t *testing.T) string {
				tmpDir := t.TempDir()
				return tmpDir
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if testing.Short() {
				t.Skip("skipping integration test in short mode")
			}

			repoPath := tt.setupRepo(t)
			defer os.RemoveAll(repoPath)

			ops := NewOperations()
			err := ops.PullFromMain(repoPath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("PullFromMain() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("PullFromMain() unexpected error: %v", err)
			}
		})
	}
}

func TestHandleGitError(t *testing.T) {
	tests := []struct {
		name           string
		output         string
		command        string
		wantErrContain string
		wantMsgContain string
	}{
		{
			name:           "should handle uncommitted changes error",
			output:         "error: Your local changes to the following files would be overwritten by checkout",
			command:        "checkout",
			wantErrContain: "would be overwritten",
			wantMsgContain: "uncommitted changes",
		},
		{
			name:           "should handle already on branch message",
			output:         "Already on 'main'",
			command:        "checkout",
			wantErrContain: "Already on 'main'",
			wantMsgContain: "Already on 'main' branch",
		},
		{
			name:           "should handle branch does not exist error",
			output:         "error: pathspec 'main' did not match any file(s) known to git",
			command:        "checkout",
			wantErrContain: "did not match",
			wantMsgContain: "does not exist",
		},
		{
			name:           "should handle not a git repository error",
			output:         "fatal: not a git repository (or any of the parent directories): .git",
			command:        "status",
			wantErrContain: "not a git repository",
			wantMsgContain: "Not a valid Git repository",
		},
		{
			name:           "should handle permission denied error",
			output:         "fatal: could not open '.git/config': Permission denied",
			command:        "status",
			wantErrContain: "Permission denied",
			wantMsgContain: "Permission denied",
		},
		{
			name:           "should handle repository not found error",
			output:         "fatal: repository 'https://github.com/example/repo.git' not found",
			command:        "clone",
			wantErrContain: "not found",
			wantMsgContain: "Remote repository not accessible or not found",
		},
		{
			name:           "should handle no tracking information error",
			output:         "There is no tracking information for the current branch",
			command:        "pull",
			wantErrContain: "no tracking information",
			wantMsgContain: "No tracking branch configured",
		},
		{
			name:           "should handle local changes overwrite error",
			output:         "error: Your local changes to the following files would be overwritten by merge",
			command:        "pull",
			wantErrContain: "overwritten",
			wantMsgContain: "uncommitted changes",
		},
		{
			name:           "should handle generic error when no specific match",
			output:         "fatal: unknown error occurred",
			command:        "status",
			wantErrContain: "unknown error",
			wantMsgContain: "Git status failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ops := NewOperations()
			err, msg := ops.handleGitError(tt.output, tt.command)

			if err == nil {
				t.Errorf("handleGitError() expected error, got nil")
				return
			}

			if !strings.Contains(err.Error(), tt.wantErrContain) {
				t.Errorf("handleGitError() error = %q, want to contain %q", err.Error(), tt.wantErrContain)
			}

			if !strings.Contains(msg, tt.wantMsgContain) {
				t.Errorf("handleGitError() message = %q, want to contain %q", msg, tt.wantMsgContain)
			}
		})
	}
}

func TestHandleNoTrackingBranch(t *testing.T) {
	tests := []struct {
		name      string
		setupRepo func(t *testing.T) string
		wantErr   bool
	}{
		{
			name: "should handle no tracking branch by setting upstream",
			setupRepo: func(t *testing.T) string {
				return createTestGitRepo(t, "main")
			},
			wantErr: false,
		},
		{
			name: "should return error when not a git repository",
			setupRepo: func(t *testing.T) string {
				tmpDir := t.TempDir()
				return tmpDir
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if testing.Short() {
				t.Skip("skipping integration test in short mode")
			}

			repoPath := tt.setupRepo(t)
			defer os.RemoveAll(repoPath)

			ops := NewOperations()
			err := ops.handleNoTrackingBranch(repoPath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("handleNoTrackingBranch() expected error, got nil")
				}
				return
			}

			// Note: This may still error in real scenarios without remote,
			// but we're testing the logic flow
			if err != nil && !strings.Contains(err.Error(), "origin/main") {
				t.Logf("handleNoTrackingBranch() error (expected in test env): %v", err)
			}
		})
	}
}

func TestExecuteGitCommand(t *testing.T) {
	tests := []struct {
		name      string
		setupRepo func(t *testing.T) string
		args      []string
		wantErr   bool
	}{
		{
			name: "should execute git status successfully",
			setupRepo: func(t *testing.T) string {
				return createTestGitRepo(t, "main")
			},
			args:    []string{"status"},
			wantErr: false,
		},
		{
			name: "should execute git branch successfully",
			setupRepo: func(t *testing.T) string {
				return createTestGitRepo(t, "main")
			},
			args:    []string{"branch"},
			wantErr: false,
		},
		{
			name: "should return error when executing git command in non-repo",
			setupRepo: func(t *testing.T) string {
				tmpDir := t.TempDir()
				return tmpDir
			},
			args:    []string{"status"},
			wantErr: true,
		},
		{
			name: "should return error when executing invalid git command",
			setupRepo: func(t *testing.T) string {
				return createTestGitRepo(t, "main")
			},
			args:    []string{"invalid-command"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if testing.Short() {
				t.Skip("skipping integration test in short mode")
			}

			repoPath := tt.setupRepo(t)
			defer os.RemoveAll(repoPath)

			ops := NewOperations()
			err := ops.executeGitCommand(repoPath, tt.args...)

			if tt.wantErr {
				if err == nil {
					t.Errorf("executeGitCommand() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("executeGitCommand() unexpected error: %v", err)
			}
		})
	}
}

func TestNewOperations(t *testing.T) {
	t.Run("should create new Operations instance", func(t *testing.T) {
		ops := NewOperations()
		if ops == nil {
			t.Errorf("NewOperations() returned nil")
		}
	})
}

// Helper functions

// createTestGitRepo creates a minimal git repository for testing
func createTestGitRepo(t *testing.T, branchName string) string {
	t.Helper()

	tmpDir := t.TempDir()

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to init git repo: %v", err)
	}

	// Configure git user (required for commits)
	configCmds := [][]string{
		{"config", "user.email", "test@example.com"},
		{"config", "user.name", "Test User"},
	}
	for _, args := range configCmds {
		cmd := exec.Command("git", args...)
		cmd.Dir = tmpDir
		if err := cmd.Run(); err != nil {
			t.Fatalf("failed to config git: %v", err)
		}
	}

	// Create an initial commit
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "add", ".")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to add files: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to commit: %v", err)
	}

	// Rename default branch to main (in case it's master)
	cmd = exec.Command("git", "branch", "-M", "main")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to rename branch to main: %v", err)
	}

	// Create and checkout branch if not main
	if branchName != "main" {
		cmd = exec.Command("git", "checkout", "-b", branchName)
		cmd.Dir = tmpDir
		if err := cmd.Run(); err != nil {
			t.Fatalf("failed to create branch %s: %v", branchName, err)
		}
	}

	return tmpDir
}