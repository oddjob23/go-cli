package git

import (
	"fmt"
	"sync"

	"github.com/oddjob23/go-cli/pkg/utils"
)

// SyncResult represents the overall result of syncing multiple repositories
type SyncResult struct {
	TotalRepositories int
	SuccessCount      int
	FailureCount      int
	Results           []OperationResult
}

// Syncer orchestrates the Git synchronization process
type Syncer struct {
	scanner    *Scanner
	operations *Operations
	output     *utils.CliOutput
}

// NewSyncer creates a new Syncer instance
func NewSyncer(output *utils.CliOutput) *Syncer {
	return &Syncer{
		scanner:    NewScanner(),
		operations: NewOperations(),
		output:     output,
	}
}

// SyncRepositories scans the directory and syncs all Git repositories in parallel
func (s *Syncer) SyncRepositories(rootDir string, branchName string) (*SyncResult, error) {
	// Scan for repositories
	repositories, err := s.scanner.ScanDirectory(rootDir)
	if err != nil {
		return nil, fmt.Errorf("failed to scan directory: %w", err)
	}

	if len(repositories) == 0 {
		return &SyncResult{
			TotalRepositories: 0,
			SuccessCount:      0,
			FailureCount:      0,
			Results:           []OperationResult{},
		}, nil
	}

	s.output.Info("Found %d repositories", len(repositories))
	s.output.Plain("")

	// Process repositories in parallel
	results := s.processRepositoriesParallel(repositories, branchName)

	// Calculate summary
	syncResult := &SyncResult{
		TotalRepositories: len(repositories),
		Results:           results,
	}

	for _, result := range results {
		if result.Success {
			syncResult.SuccessCount++
		} else {
			syncResult.FailureCount++
		}
	}

	return syncResult, nil
}

// processRepositoriesParallel processes multiple repositories concurrently using goroutines
func (s *Syncer) processRepositoriesParallel(repositories []Repository, branchName string) []OperationResult {
	var wg sync.WaitGroup
	results := make([]OperationResult, len(repositories))

	// Process each repository in a separate goroutine
	for i, repo := range repositories {
		wg.Add(1)
		go func(index int, repository Repository) {
			defer wg.Done()

			s.output.Plain("  📂 %s", repository.Name)

			result := s.operations.CheckoutMainBranch(repository, branchName)
			results[index] = result

			if result.Success {
				s.output.Plain("     ✅ %s", result.Message)
			} else {
				s.output.Plain("     ❌ %s", result.Message)
			}
		}(i, repo)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	return results
}

// PrintSummary prints a summary of the sync operation results
func (s *Syncer) PrintSummary(result *SyncResult) {
	s.output.Plain("")
	s.output.Plain("Summary:")
	s.output.Plain("  Total: %d", result.TotalRepositories)
	s.output.Plain("  Successful: %d", result.SuccessCount)
	s.output.Plain("  Failed: %d", result.FailureCount)
}
