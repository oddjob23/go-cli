package commands

import (
	"fmt"
	"os"

	"github.com/oddjob23/go-cli/internal/git"
	"github.com/oddjob23/go-cli/pkg/config"
	"github.com/oddjob23/go-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync Git repositories in a directory",
	Long: `Scans a directory for Git repositories and syncs them by checking out
the main branch and pulling the latest changes. Processes repositories in parallel.`,
	RunE: runSync,
}

func runSync(cmd *cobra.Command, args []string) error {
	// Get flags
	directory, _ := cmd.Flags().GetString("directory")
	branch, _ := cmd.Flags().GetString("branch")
	envFile, _ := cmd.Flags().GetString("env-file")

	// Load configuration
	cfg, err := config.LoadFromFile(envFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Override with command line flags if provided
	if directory != "" {
		cfg.ScanDirectory = directory
	}
	if branch != "" {
		cfg.GitBranch = branch
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Create output handler
	output := utils.NewCliOutput(false) // Set to true for verbose mode if needed

	// Create syncer and run sync
	syncer := git.NewSyncer(output)

	output.Info("Starting Git repository sync in: %s", cfg.ScanDirectory)
	output.Info("Target branch: %s", cfg.GitBranch)

	result, err := syncer.SyncRepositories(cfg.ScanDirectory, cfg.GitBranch)
	if err != nil {
		output.Error("Sync failed: %v", err)
		return err
	}

	// Print results
	syncer.PrintSummary(result)

	// Print final summary
	if result.TotalRepositories == 0 {
		output.Warning("No Git repositories found in the specified directory")
		return nil
	}

	if result.FailureCount == 0 {
		output.Success("All %d repositories synced successfully!", result.SuccessCount)
	} else {
		output.Warning("Synced %d/%d repositories successfully. %d failed.",
			result.SuccessCount, result.TotalRepositories, result.FailureCount)

		// Exit with error code if any repositories failed
		os.Exit(1)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
