package docker

import (
	"fmt"
	"path/filepath"

	"github.com/oddjob23/go-cli/pkg/utils"
)

type Manager struct {
	dependenciesFile string
	servicesFile     string
}

func NewManager(baseDir string) *Manager {
	return &Manager{
		dependenciesFile: filepath.Join(baseDir, "docker-compose.dependencies.yml"),
		servicesFile:     filepath.Join(baseDir, "docker-compose.services.yml"),
	}
}

func (m *Manager) StartDependencies() error {
	utils.Info("Starting Docker dependencies workflow...")

	if err := CheckDockerDaemon(); err != nil {
		return fmt.Errorf("docker daemon check failed: %w", err)
	}

	if err := CheckDockerCompose(); err != nil {
		return fmt.Errorf("docker compose check failed: %w", err)
	}

	config := NewComposeConfig(m.dependenciesFile)
	if err := config.StartDependencies(); err != nil {
		return fmt.Errorf("failed to start dependencies: %w", err)
	}

	utils.Success("Dependencies workflow completed successfully")
	return nil
}

func (m *Manager) StartServices() error {
	utils.Info("Starting microservices...")

	config := NewComposeConfig(m.servicesFile)
	if err := config.StartDependencies(); err != nil {
		return fmt.Errorf("failed to start services: %w", err)
	}

	utils.Success("Services started successfully")
	return nil
}

func (m *Manager) StartAll() error {
	if err := m.StartDependencies(); err != nil {
		return err
	}

	if err := m.StartServices(); err != nil {
		return err
	}

	utils.Success("All services started successfully")
	return nil
}

func (m *Manager) Stop() error {
	utils.Info("Stopping all services...")

	servicesConfig := NewComposeConfig(m.servicesFile)
	if err := servicesConfig.Stop(); err != nil {
		utils.Warning("Failed to stop services: " + err.Error())
	}

	dependenciesConfig := NewComposeConfig(m.dependenciesFile)
	if err := dependenciesConfig.Stop(); err != nil {
		return fmt.Errorf("failed to stop dependencies: %w", err)
	}

	utils.Success("All services stopped successfully")
	return nil
}

func (m *Manager) Status() error {
	utils.Info("Checking service status...")

	dependenciesConfig := NewComposeConfig(m.dependenciesFile)
	utils.Info("Dependencies status:")
	if err := dependenciesConfig.showServiceStatus(); err != nil {
		utils.Warning("Failed to get dependencies status")
	}

	servicesConfig := NewComposeConfig(m.servicesFile)
	utils.Info("Services status:")
	if err := servicesConfig.showServiceStatus(); err != nil {
		utils.Warning("Failed to get services status")
	}

	return nil
}

func (m *Manager) Logs(serviceName string) error {
	dependenciesConfig := NewComposeConfig(m.dependenciesFile)
	return dependenciesConfig.Logs(serviceName)
}