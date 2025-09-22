package commands

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/oddjob23/go-cli/internal/docker"
	"github.com/oddjob23/go-cli/pkg/config"
	"github.com/oddjob23/go-cli/pkg/utils"
)

var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Manage Docker containers and dependencies",
	Long:  `Manage Docker containers, dependencies, and microservices using Docker Compose`,
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Docker services",
	Long:  `Start Docker dependencies and services using Docker Compose files`,
}

var startDepsCmd = &cobra.Command{
	Use:   "deps",
	Short: "Start dependencies only",
	Long:  `Start shared dependencies (databases, message queues, etc.) using docker-compose.dependencies.yml`,
	RunE:  runStartDependencies,
}

var startServicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Start microservices only",
	Long:  `Start microservices using docker-compose.services.yml`,
	RunE:  runStartServices,
}

var startAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Start all services",
	Long:  `Start both dependencies and microservices`,
	RunE:  runStartAll,
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop all Docker services",
	Long:  `Stop all running Docker services and dependencies`,
	RunE:  runStop,
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Docker services status",
	Long:  `Display the status of all Docker services and dependencies`,
	RunE:  runStatus,
}

var logsCmd = &cobra.Command{
	Use:   "logs [service-name]",
	Short: "Show logs for Docker services",
	Long:  `Display logs for Docker services. If no service name is provided, shows logs for all services`,
	RunE:  runLogs,
}

func runStartDependencies(cmd *cobra.Command, args []string) error {
	workingDir, err := loadDockerConfig(cmd)
	if err != nil {
		return err
	}

	manager := docker.NewManager(workingDir)
	return manager.StartDependencies()
}

func runStartServices(cmd *cobra.Command, args []string) error {
	workingDir, err := loadDockerConfig(cmd)
	if err != nil {
		return err
	}

	manager := docker.NewManager(workingDir)
	return manager.StartServices()
}

func runStartAll(cmd *cobra.Command, args []string) error {
	workingDir, err := loadDockerConfig(cmd)
	if err != nil {
		return err
	}

	manager := docker.NewManager(workingDir)
	return manager.StartAll()
}

func runStop(cmd *cobra.Command, args []string) error {
	workingDir, err := loadDockerConfig(cmd)
	if err != nil {
		return err
	}

	manager := docker.NewManager(workingDir)
	return manager.Stop()
}

func runStatus(cmd *cobra.Command, args []string) error {
	workingDir, err := loadDockerConfig(cmd)
	if err != nil {
		return err
	}

	manager := docker.NewManager(workingDir)
	return manager.Status()
}

func runLogs(cmd *cobra.Command, args []string) error {
	workingDir, err := loadDockerConfig(cmd)
	if err != nil {
		return err
	}

	serviceName := ""
	if len(args) > 0 {
		serviceName = args[0]
	}

	manager := docker.NewManager(workingDir)
	return manager.Logs(serviceName)
}

func loadDockerConfig(cmd *cobra.Command) (string, error) {
	directory, _ := cmd.Flags().GetString("directory")
	envFile, _ := cmd.Flags().GetString("env-file")

	cfg, err := config.LoadFromFile(envFile)
	if err != nil {
		utils.Error("Failed to load configuration: " + err.Error())
		return "", err
	}

	workingDir := directory
	if workingDir == "" {
		if cfg.ScanDirectory != "" {
			workingDir = cfg.ScanDirectory
		} else {
			wd, err := os.Getwd()
			if err != nil {
				return "", err
			}
			workingDir = wd
		}
	}

	workingDir, err = filepath.Abs(workingDir)
	if err != nil {
		return "", err
	}

	return workingDir, nil
}

func init() {
	rootCmd.AddCommand(dockerCmd)

	dockerCmd.AddCommand(startCmd)
	startCmd.AddCommand(startDepsCmd)
	startCmd.AddCommand(startServicesCmd)
	startCmd.AddCommand(startAllCmd)

	dockerCmd.AddCommand(stopCmd)
	dockerCmd.AddCommand(statusCmd)
	dockerCmd.AddCommand(logsCmd)
}