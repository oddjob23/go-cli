package docker

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/oddjob23/go-cli/pkg/utils"
)

func CheckDockerDaemon() error {
	utils.Info("Checking Docker daemon status...")

	cmd := exec.Command("docker", "info")
	err := cmd.Run()
	if err != nil {
		utils.Warning("Docker daemon is not running")
		return StartDockerDaemon()
	}

	utils.Success("Docker daemon is running")
	return nil
}

func StartDockerDaemon() error {
	utils.Info("Starting Docker daemon...")

	cmd := exec.Command("open", "-a", "Docker")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to start Docker daemon: %w", err)
	}

	return waitForDockerDaemon()
}

func waitForDockerDaemon() error {
	utils.Info("Waiting for Docker daemon to start...")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for Docker daemon to start")
		case <-ticker.C:
			cmd := exec.Command("docker", "info")
			if cmd.Run() == nil {
				utils.Success("Docker daemon started successfully")
				return nil
			}
		}
	}
}

func CheckDockerCompose() error {
	utils.Info("Checking Docker Compose availability...")

	cmd := exec.Command("docker", "compose", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		cmd = exec.Command("docker-compose", "version")
		output, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("docker-compose not found: %w", err)
		}
		utils.Success("docker-compose (standalone) is available")
		return nil
	}

	version := strings.TrimSpace(string(output))
	utils.Success(fmt.Sprintf("Docker Compose is available: %s", version))
	return nil
}

func GetDockerComposeCommand() string {
	cmd := exec.Command("docker", "compose", "version")
	if cmd.Run() == nil {
		return "docker compose"
	}
	return "docker-compose"
}