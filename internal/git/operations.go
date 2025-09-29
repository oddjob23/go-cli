package git

import (
	"fmt"
	"os/exec"
	"strings"
)

const (
	mainBranch = "main"
)

// OperationResult represents the result of a Git operation
type OperationResult struct {
	Repository Repository
	Success    bool
	Error      error
	Message    string
}

// Operations handles Git operations on repositories
type Operations struct{}

// NewOperations creates a new Operations instance
func NewOperations() *Operations {
	return &Operations{}
}

// CheckoutMainBranch attempts to checkout the main branch for a repository
func (o *Operations) CheckoutMainBranch(repo Repository, branchName string) OperationResult {
	result := OperationResult{
		Repository: repo,
		Success:    false,
	}

	// Get current branch
	currentBranch, err := o.getCurrentBranch(repo.Path)
	if err != nil {
		result.Error = fmt.Errorf("failed to get current branch: %w", err)
		result.Message = result.Error.Error()
		return result
	}

	// Checkout main branch if not already on it
	if currentBranch != mainBranch {
		err = o.executeGitCommand(repo.Path, "checkout", mainBranch)
		if err != nil {
			result.Error, result.Message = o.handleGitError(err.Error(), "checkout")
			return result
		}
	}

	// Pull latest changes from main
	err = o.PullFromMain(repo.Path)
	if err != nil {
		result.Error, result.Message = o.handleGitError(err.Error(), "pull")
		return result
	}

	result.Success = true
	result.Message = fmt.Sprintf("Checked out '%s' and pulled latest changes", mainBranch)
	return result
}


// getCurrentBranch gets the current branch name
func (o *Operations) getCurrentBranch(repoPath string) (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// PullFromMain pulls the latest changes from the main branch
func (o *Operations) PullFromMain(repoPath string) error {
	// Try regular pull first
	err := o.executeGitCommand(repoPath, "pull")
	if err == nil {
		return nil
	}

	// If pull fails, handle tracking issues
	if strings.Contains(err.Error(), "no tracking information") {
		return o.handleNoTrackingBranch(repoPath)
	}

	// Return the original error
	return err
}

// handleNoTrackingBranch handles the case when branch has no tracking information
func (o *Operations) handleNoTrackingBranch(repoPath string) error {
	// First, fetch to make sure we have latest remote info
	err := o.executeGitCommand(repoPath, "fetch")
	if err != nil {
		return fmt.Errorf("failed to fetch: %w", err)
	}

	// Try to set upstream tracking for main
	err = o.executeGitCommand(repoPath, "branch", "--set-upstream-to=origin/"+mainBranch, mainBranch)
	if err != nil {
		// If setting upstream fails, try pull with explicit remote and branch
		err = o.executeGitCommand(repoPath, "pull", "origin", mainBranch)
		if err != nil {
			return fmt.Errorf("failed to pull from origin/%s: %w", mainBranch, err)
		}
		return nil
	}

	// Now try pull again
	err = o.executeGitCommand(repoPath, "pull")
	if err != nil {
		return fmt.Errorf("failed to pull after setting upstream: %w", err)
	}

	return nil
}

// handleGitError analyzes git command output and returns user-friendly messages
func (o *Operations) handleGitError(output string, command string) (error, string) {
	outputLower := strings.ToLower(output)

	// Check for common git errors in the output
	switch {
	case strings.Contains(outputLower, "uncommitted changes") || strings.Contains(outputLower, "would be overwritten"):
		return fmt.Errorf("%s", output), "Skipped: Repository has uncommitted changes. Please commit or stash changes first."
	case strings.Contains(outputLower, "already on") && strings.Contains(outputLower, mainBranch):
		return fmt.Errorf("%s", output), fmt.Sprintf("Already on '%s' branch", mainBranch)
	case strings.Contains(outputLower, "did not match any file") || (strings.Contains(outputLower, "pathspec") && strings.Contains(outputLower, "did not match")):
		return fmt.Errorf("%s", output), fmt.Sprintf("Branch '%s' does not exist in this repository", mainBranch)
	case strings.Contains(outputLower, "not a git repository"):
		return fmt.Errorf("%s", output), "Not a valid Git repository"
	case strings.Contains(outputLower, "no such file or directory"):
		return fmt.Errorf("%s", output), "Repository path does not exist"
	case strings.Contains(outputLower, "permission denied"):
		return fmt.Errorf("%s", output), "Permission denied accessing repository"
	case strings.Contains(outputLower, "repository not found") || strings.Contains(outputLower, "could not read from remote"):
		return fmt.Errorf("%s", output), "Remote repository not accessible or not found"
	case strings.Contains(outputLower, "no tracking information"):
		return fmt.Errorf("%s", output), "No tracking branch configured for this branch"
	case strings.Contains(outputLower, "your local changes to the following files"):
		return fmt.Errorf("%s", output), "Local changes would be overwritten. Please commit or stash changes first."
	default:
		return fmt.Errorf("%s", output), fmt.Sprintf("Git %s failed: %s", command, output)
	}
}

// executeGitCommand executes a git command in the specified directory
func (o *Operations) executeGitCommand(repoPath string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = repoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}

	return nil
}
