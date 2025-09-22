package docker

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/oddjob23/go-cli/pkg/utils"
)

type ComposeConfig struct {
	FilePath    string
	ProjectName string
	Timeout     time.Duration
}

func NewComposeConfig(filePath string) *ComposeConfig {
	return &ComposeConfig{
		FilePath:    filePath,
		ProjectName: "microservices",
		Timeout:     5 * time.Minute,
	}
}

func (c *ComposeConfig) ValidateFile() error {
	if _, err := os.Stat(c.FilePath); os.IsNotExist(err) {
		return fmt.Errorf("docker-compose file not found: %s", c.FilePath)
	}
	return nil
}

func (c *ComposeConfig) StartDependencies() error {
	utils.Info(fmt.Sprintf("Starting dependencies from %s...", filepath.Base(c.FilePath)))

	if err := c.ValidateFile(); err != nil {
		return err
	}

	composeCmd := GetDockerComposeCommand()
	args := []string{"-f", c.FilePath, "-p", c.ProjectName, "up", "-d", "--build"}

	cmd := exec.Command("sh", "-c", fmt.Sprintf("%s %s", composeCmd, strings.Join(args, " ")))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start dependencies: %w", err)
	}

	utils.Success("Dependencies started successfully")
	return c.WaitForHealthChecks()
}

func (c *ComposeConfig) WaitForHealthChecks() error {
	utils.Info("Waiting for services to become healthy...")

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			utils.Warning("Timeout waiting for all services to become healthy")
			return c.showServiceStatus()
		case <-ticker.C:
			if healthy, err := c.checkAllServicesHealthy(); err != nil {
				return err
			} else if healthy {
				utils.Success("All services are healthy")
				return nil
			}
			utils.Info("Some services are still starting...")
		}
	}
}

func (c *ComposeConfig) checkAllServicesHealthy() (bool, error) {
	composeCmd := GetDockerComposeCommand()
	cmd := exec.Command("sh", "-c", fmt.Sprintf("%s -f %s -p %s ps --format json", composeCmd, c.FilePath, c.ProjectName))

	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check service status: %w", err)
	}

	if len(output) == 0 {
		return false, nil
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if strings.Contains(line, "unhealthy") || strings.Contains(line, "starting") {
			return false, nil
		}
	}

	return true, nil
}

func (c *ComposeConfig) showServiceStatus() error {
	composeCmd := GetDockerComposeCommand()
	cmd := exec.Command("sh", "-c", fmt.Sprintf("%s -f %s -p %s ps", composeCmd, c.FilePath, c.ProjectName))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *ComposeConfig) Stop() error {
	utils.Info("Stopping services...")

	composeCmd := GetDockerComposeCommand()
	cmd := exec.Command("sh", "-c", fmt.Sprintf("%s -f %s -p %s down", composeCmd, c.FilePath, c.ProjectName))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop services: %w", err)
	}

	utils.Success("Services stopped successfully")
	return nil
}

func (c *ComposeConfig) Logs(serviceName string) error {
	composeCmd := GetDockerComposeCommand()
	args := []string{"-f", c.FilePath, "-p", c.ProjectName, "logs", "-f"}

	if serviceName != "" {
		args = append(args, serviceName)
	}

	cmd := exec.Command("sh", "-c", fmt.Sprintf("%s %s", composeCmd, strings.Join(args, " ")))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}