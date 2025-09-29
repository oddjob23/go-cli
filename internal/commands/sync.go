package commands

import (
	"fmt"
	"os"
	"sync"

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
	configFile, _ := cmd.Flags().GetString("config")
	branch, _ := cmd.Flags().GetString("branch")

	// Load configuration
	cfg, err := config.LoadFromFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Override branch if provided via command line
	if branch != "" {
		cfg.GitBranch = branch
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Create output handler
	output := utils.NewCliOutput(false) // Set to true for verbose mode if needed

	// Create syncer
	syncer := git.NewSyncer(output)

	output.Info("Starting Git repository sync for %d configured repositories", len(cfg.Repositories))
	output.Info("Target branch: %s", cfg.GitBranch)

	// Sync each configured repository in parallel
	var wg sync.WaitGroup
	var mu sync.Mutex
	var successCount, failureCount int

	output.Plain("")

	for _, repo := range cfg.Repositories {
		wg.Add(1)
		go func(r config.Repository) {
			defer wg.Done()

			output.Plain("  📂 %s", r.Name)
			err := syncer.SyncSingleRepository(r.Path, cfg.GitBranch)

			mu.Lock()
			if err != nil {
				output.Plain("     ❌ Failed to sync - %s", err.Error())
				failureCount++
			} else {
				output.Plain("    ✅  Successfully pulled %s branch", cfg.GitBranch)
				successCount++
			}
			mu.Unlock()
		}(repo)
	}

	// Wait for all repositories to complete
	wg.Wait()

	// Print final summary
	if len(cfg.Repositories) == 0 {
		output.Warning("No repositories configured")
		return nil
	}

	if failureCount == 0 {
		output.Success("All %d repositories synced successfully!", successCount)
	} else {
		output.Warning("Synced %d/%d repositories successfully. %d failed.",
			successCount, len(cfg.Repositories), failureCount)

		// Exit with error code if any repositories failed
		os.Exit(1)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
