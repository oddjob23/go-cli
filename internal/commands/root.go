package commands

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-cli",
	Short: "A CLI tool for Git repositories and Docker management",
	Long: `A CLI application that:
- Scans directories for Git repositories, checks out main branch, and pulls the latest changes
- Manages Docker containers and microservices using Docker Compose`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "config.json", "Path to config.json file")
	rootCmd.PersistentFlags().StringP("branch", "b", "main", "Git branch to checkout and pull")
}
