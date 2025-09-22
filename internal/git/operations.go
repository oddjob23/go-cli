package git

import (
	"fmt"
	"os/exec"
	"strings"
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

	// Check for uncommitted changes first
	if hasUncommittedChanges, err := o.hasUncommittedChanges(repo.Path); err != nil {
		result.Error = fmt.Errorf("failed to check repository status: %w", err)
		result.Message = result.Error.Error()
		return result
	} else if hasUncommittedChanges {
		result.Error = fmt.Errorf("repository has uncommitted changes")
		result.Message = "Skipped: Repository has uncommitted changes. Please commit or stash changes first."
		return result
	}

	// First, try to determine the default branch (main or master)
	defaultBranch, err := o.getDefaultBranch(repo.Path)
	if err != nil {
		result.Error = fmt.Errorf("failed to determine default branch: %w", err)
		result.Message = result.Error.Error()
		return result
	}

	// Use the provided branch name if specified, otherwise use detected default
	targetBranch := branchName
	if targetBranch == "" || targetBranch == "main" {
		targetBranch = defaultBranch
	}

	// Check if we're already on the target branch
	currentBranch, err := o.getCurrentBranch(repo.Path)
	if err != nil {
		result.Error = fmt.Errorf("failed to get current branch: %w", err)
		result.Message = result.Error.Error()
		return result
	}

	// Checkout the target branch if not already on it
	if currentBranch != targetBranch {
		err = o.executeGitCommand(repo.Path, "checkout", targetBranch)
		if err != nil {
			result.Error = fmt.Errorf("failed to checkout branch '%s': %w", targetBranch, err)
			result.Message = result.Error.Error()
			return result
		}
	}

	// Handle pull with better error handling
	err = o.pullWithFallback(repo.Path, targetBranch)
	if err != nil {
		result.Error = fmt.Errorf("failed to pull latest changes: %w", err)
		result.Message = result.Error.Error()
		return result
	}

	result.Success = true
	if currentBranch != targetBranch {
		result.Message = fmt.Sprintf("Checked out '%s' and pulled latest changes", targetBranch)
	} else {
		result.Message = fmt.Sprintf("Already on '%s', pulled latest changes", targetBranch)
	}
	return result
}

// hasUncommittedChanges checks if repository has uncommitted changes
func (o *Operations) hasUncommittedChanges(repoPath string) (bool, error) {
	// Check for staged changes
	cmd := exec.Command("git", "diff", "--cached", "--quiet")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		// Non-zero exit means there are staged changes
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
			return true, nil
		}
		return false, fmt.Errorf("failed to check staged changes: %w", err)
	}

	// Check for unstaged changes
	cmd = exec.Command("git", "diff", "--quiet")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		// Non-zero exit means there are unstaged changes
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
			return true, nil
		}
		return false, fmt.Errorf("failed to check unstaged changes: %w", err)
	}

	return false, nil
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

// pullWithFallback attempts to pull with fallback for tracking issues
func (o *Operations) pullWithFallback(repoPath string, branchName string) error {
	// Try regular pull first
	err := o.executeGitCommand(repoPath, "pull")
	if err == nil {
		return nil
	}

	// If pull fails, check if it's a tracking issue
	if strings.Contains(err.Error(), "no tracking information") {
		// Try to set upstream and pull
		return o.setupTrackingAndPull(repoPath, branchName)
	}

	// If it's not a tracking issue, return the original error
	return err
}

// setupTrackingAndPull sets up tracking branch and pulls
func (o *Operations) setupTrackingAndPull(repoPath string, branchName string) error {
	// First, fetch to make sure we have latest remote info
	err := o.executeGitCommand(repoPath, "fetch")
	if err != nil {
		return fmt.Errorf("failed to fetch: %w", err)
	}

	// Try to set upstream tracking
	err = o.executeGitCommand(repoPath, "branch", "--set-upstream-to=origin/"+branchName, branchName)
	if err != nil {
		// If setting upstream fails, try pull with explicit remote and branch
		err = o.executeGitCommand(repoPath, "pull", "origin", branchName)
		if err != nil {
			return fmt.Errorf("failed to pull from origin/%s: %w", branchName, err)
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

// getDefaultBranch determines the default branch (main or master)
func (o *Operations) getDefaultBranch(repoPath string) (string, error) {
	// Try to get the default branch from remote
	cmd := exec.Command("git", "symbolic-ref", "refs/remotes/origin/HEAD")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err == nil {
		// Extract branch name from refs/remotes/origin/HEAD -> refs/remotes/origin/main
		parts := strings.Split(strings.TrimSpace(string(output)), "/")
		if len(parts) > 0 {
			return parts[len(parts)-1], nil
		}
	}

	// Fallback: check if main branch exists
	err = o.executeGitCommand(repoPath, "show-ref", "--verify", "--quiet", "refs/heads/main")
	if err == nil {
		return "main", nil
	}

	// Fallback: check if master branch exists
	err = o.executeGitCommand(repoPath, "show-ref", "--verify", "--quiet", "refs/heads/master")
	if err == nil {
		return "master", nil
	}

	// Final fallback
	return "main", nil
}

// executeGitCommand executes a git command in the specified directory
func (o *Operations) executeGitCommand(repoPath string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = repoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git command failed: %s (output: %s)", err.Error(), string(output))
	}

	return nil
}
