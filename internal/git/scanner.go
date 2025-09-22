package git

import (
	"os"
	"path/filepath"
)

// Repository represents a Git repository with its path
type Repository struct {
	Path string
	Name string
}

// Scanner handles scanning directories for Git repositories
type Scanner struct{}

// NewScanner creates a new Scanner instance
func NewScanner() *Scanner {
	return &Scanner{}
}

// ScanDirectory scans the given directory for Git repositories
// Returns a slice of Repository structs representing found Git repos
func (s *Scanner) ScanDirectory(rootDir string) ([]Repository, error) {
	var repositories []Repository

	// Check if root directory exists
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		return nil, err
	}

	// Read all entries in the root directory
	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return nil, err
	}

	// Check each subdirectory for .git folder
	for _, entry := range entries {
		if entry.IsDir() {
			dirPath := filepath.Join(rootDir, entry.Name())
			gitPath := filepath.Join(dirPath, ".git")

			// Check if .git directory exists
			if _, err := os.Stat(gitPath); err == nil {
				repositories = append(repositories, Repository{
					Path: dirPath,
					Name: entry.Name(),
				})
			}
		}
	}

	return repositories, nil
}
