package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestSyncIntegration tests the sync command with real git repositories
func TestSyncIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "go-cli-integration-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Use the pre-built binary
	binaryPath := "../../bin/go-cli"
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Skipf("Binary not found at %s. Run 'make build' first.", binaryPath)
	}

	// Clone a small test repository
	repoDir := filepath.Join(tempDir, "test-repos")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatalf("Failed to create repo directory: %v", err)
	}

	// Clone Hello-World repository (small and fast)
	cloneCmd := exec.Command("git", "clone", "--quiet", "https://github.com/octocat/Hello-World.git", "hello-world")
	cloneCmd.Dir = repoDir
	if err := cloneCmd.Run(); err != nil {
		t.Skipf("Failed to clone test repository (network issue?): %v", err)
	}

	// Run the sync command
	syncCmd := exec.Command(binaryPath, "sync", "-d", repoDir)
	output, err := syncCmd.CombinedOutput()

	// Validate results
	if err != nil {
		t.Errorf("Sync command failed: %v\nOutput: %s", err, string(output))
		return
	}

	outputStr := string(output)

	// Check that the command found repositories
	if !strings.Contains(outputStr, "Found 1 repositories") {
		t.Errorf("Expected to find 1 repository, but output was: %s", outputStr)
	}

	// Check that sync was successful
	if !strings.Contains(outputStr, "repositories synced successfully") {
		t.Errorf("Expected successful sync message, but output was: %s", outputStr)
	}

	// Verify the repository is on the correct branch
	checkBranchCmd := exec.Command("git", "branch", "--show-current")
	checkBranchCmd.Dir = filepath.Join(repoDir, "hello-world")
	branchOutput, err := checkBranchCmd.Output()
	if err != nil {
		t.Errorf("Failed to check current branch: %v", err)
		return
	}

	currentBranch := strings.TrimSpace(string(branchOutput))
	if currentBranch != "master" && currentBranch != "main" {
		t.Errorf("Expected repository to be on master or main branch, but was on: %s", currentBranch)
	}

	t.Logf("Integration test completed successfully. Repository is on branch: %s", currentBranch)
}

// TestSyncWithMultipleRepos tests sync with multiple repositories
func TestSyncWithMultipleRepos(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Set a timeout for this test since it involves network operations
	timeout := 2 * time.Minute
	done := make(chan bool)

	go func() {
		defer func() { done <- true }()

		// Create temporary directory for test
		tempDir, err := os.MkdirTemp("", "go-cli-multi-integration-*")
		if err != nil {
			t.Errorf("Failed to create temp directory: %v", err)
			return
		}
		defer os.RemoveAll(tempDir)

		// Use the pre-built binary
		binaryPath := "../../bin/go-cli"
		if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
			t.Errorf("Binary not found at %s. Run 'make build' first.", binaryPath)
			return
		}

		// Create repo directory
		repoDir := filepath.Join(tempDir, "test-repos")
		if err := os.MkdirAll(repoDir, 0755); err != nil {
			t.Errorf("Failed to create repo directory: %v", err)
			return
		}

		// Clone multiple small repositories
		repos := []struct {
			url  string
			name string
		}{
			{"https://github.com/octocat/Hello-World.git", "hello-world"},
			{"https://github.com/github/gitignore.git", "gitignore"},
		}

		clonedCount := 0
		for _, repo := range repos {
			cloneCmd := exec.Command("git", "clone", "--quiet", "--depth", "1", repo.url, repo.name)
			cloneCmd.Dir = repoDir
			if err := cloneCmd.Run(); err == nil {
				clonedCount++
			}
		}

		if clonedCount == 0 {
			t.Skip("No repositories could be cloned (network issues?)")
			return
		}

		// Run the sync command
		syncCmd := exec.Command(binaryPath, "sync", "-d", repoDir)
		output, err := syncCmd.CombinedOutput()

		if err != nil {
			t.Errorf("Sync command failed: %v\nOutput: %s", err, string(output))
			return
		}

		outputStr := string(output)

		// Check that the command found the expected number of repositories
		expectedMsg := "Found " + string(rune('0'+clonedCount)) + " repositories"
		if !strings.Contains(outputStr, expectedMsg) {
			t.Errorf("Expected to find %d repositories, but output was: %s", clonedCount, outputStr)
		}

		t.Logf("Multi-repository integration test completed successfully with %d repositories", clonedCount)
	}()

	select {
	case <-done:
		// Test completed successfully
	case <-time.After(timeout):
		t.Fatal("Integration test timed out after", timeout)
	}
}

// TestSyncEmptyDirectory tests sync with an empty directory
func TestSyncEmptyDirectory(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "go-cli-empty-integration-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Use the pre-built binary
	binaryPath := "../../bin/go-cli"
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Skipf("Binary not found at %s. Run 'make build' first.", binaryPath)
	}

	// Create empty repo directory
	repoDir := filepath.Join(tempDir, "empty-repos")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatalf("Failed to create repo directory: %v", err)
	}

	// Run the sync command on empty directory
	syncCmd := exec.Command(binaryPath, "sync", "-d", repoDir)
	output, err := syncCmd.CombinedOutput()

	if err != nil {
		t.Errorf("Sync command failed: %v\nOutput: %s", err, string(output))
		return
	}

	outputStr := string(output)

	// Check for expected empty directory message
	if !strings.Contains(outputStr, "No Git repositories found") {
		t.Errorf("Expected 'No Git repositories found' message, but output was: %s", outputStr)
	}

	t.Log("Empty directory integration test completed successfully")
}