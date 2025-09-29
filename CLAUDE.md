# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go CLI application with two main capabilities:
1. **Git Repository Sync**: Scans a directory for Git repositories, checks out main branch, pulls latest changes, and provides a summary of successful/failed operations with parallel processing
2. **Microservices Management**: Starts shared dependencies and microservices using Docker Compose, including Docker health checks and ngrok configuration

## Project Structure

```
.
├── cmd/
│   └── main.go          # Application entry point
├── internal/            # Private application code
│   ├── commands/        # CLI commands and cobra setup
│   ├── git/             # Git operations (scan, sync, pull)
│   ├── docker/          # Docker and Docker Compose management
│   └── ngrok/           # ngrok tunnel management
├── pkg/                 # Public libraries that can be imported
│   ├── config/          # Configuration management
│   ├── logger/          # Structured logging utilities
│   └── utils/           # Common utilities and CLI output
├── test/
│   ├── unit/            # Unit tests
│   └── integration/     # Integration tests
├── scripts/             # Build and deployment scripts
│   └── integration-test.sh  # Comprehensive integration test script
├── docs/                # Documentation
├── Makefile            # Build, test, and integration commands
├── go.mod              # Go module definition
└── go.sum              # Go module checksums
```

## Development Commands

### Makefile Commands (Recommended)
```bash
# Build the application
make build

# Run the application
make run

# Run unit tests
make test

# Run integration tests with script
make test-integration

# Demo integration test with real repos
make integration-demo

# Clean build artifacts
make clean

# Show all available commands
make help
```

### Manual Build and Run
```bash
# Set Go toolchain (required)
export GOTOOLCHAIN=go1.25.1

# Build manually
go build -o bin/go-cli cmd/main.go

# Run the CLI
./bin/go-cli --help
./bin/go-cli sync --help
./bin/go-cli sync -d /path/to/repositories
```

### Testing
```bash
# Run all unit tests
export GOTOOLCHAIN=go1.25.1 && go test ./...

# Run unit tests only
go test ./internal/git/... -v

# Run integration tests
go test ./test/integration/... -v -timeout=5m

# Run tests with coverage
go test -cover ./...

# Integration test with real repositories
make test-integration
```

### Code Quality
```bash
# Format code
export GOTOOLCHAIN=go1.25.1 && go fmt ./...

# Vet code
export GOTOOLCHAIN=go1.25.1 && go vet ./...

# Tidy modules
export GOTOOLCHAIN=go1.25.1 && go mod tidy
```

## Dependencies

- **Go Version**: Requires Go 1.25.1 (configured in go.mod)
- **CLI Framework**: Cobra for command-line interface
- **Coloring**: fatih/color for terminal output
- **Configuration**: godotenv for environment file support
- **Git Operations**: Uses `git` command via os/exec
- **Output**: Custom CLI output system with colored, user-friendly messages

## Architecture Notes

### Git Sync Module (`internal/git/`)
- **Scanner**: Recursively scans directory for `.git` folders
- **Operations**: Handles Git operations (checkout, pull) with advanced error handling
- **Sync**: Orchestrates parallel repository processing using goroutines
- **Error Handling**: Detects uncommitted changes, tracking branch issues, and provides actionable error messages

### CLI Interface (`internal/commands/`)
- **Root Command**: Main application entry point with global flags
- **Sync Command**: Git repository synchronization with directory and branch options
- **Output System**: Clean, minimalistic CLI output with emoji indicators and colors

### Utilities (`pkg/utils/`)
- **CLI Output**: Simple, colored output system (Info, Success, Warning, Error, Debug)
- **Configuration**: Environment variable and CLI flag support
- **Legacy Support**: Backward-compatible print functions

## Usage Examples

```bash
# Sync repositories in a directory
./bin/go-cli sync -d /Users/username/projects

# Sync with specific branch
./bin/go-cli sync -d /path/to/repos -b develop

# Use environment file
./bin/go-cli sync -e .env

# Integration testing
make integration-demo
make test-integration
```

## Development Milestones

### ✅ Milestone 1: Setup project in Go - COMPLETED
- [x] Recommended folder and file structure
- [x] Basic main.go with logging
- [x] Go module initialization

### ✅ Milestone 2: Install CLI dependencies - COMPLETED
- [x] Add minimal CLI framework for clear input/output (Cobra)
- [x] Implement colored output for success/error states
- [x] Setup configuration management (ENV + CLI flags)
- [x] Create unit tests with table-driven approach
- [x] Add Makefile with build, run, test, clean commands
- [x] Configure colored test output (green PASS, red FAIL)

### ✅ Milestone 3: Git repository sync - COMPLETED
- [x] Implement directory scanning for Git repositories
- [x] Add Git operations (checkout main, pull) with error handling
- [x] Parallel processing with goroutines
- [x] Advanced error handling (uncommitted changes, tracking branches)
- [x] Clean CLI output system replacing verbose slog
- [x] Comprehensive unit and integration tests
- [x] Integration test automation with real repositories

### ⏳ Milestone 4: Docker Compose management
- [ ] Docker daemon health check and startup
- [ ] Docker Compose file parsing and execution
- [ ] Shared dependencies management
- [ ] Unit and integration tests

### ⏳ Milestone 5: ngrok integration
- [ ] ngrok tunnel management with --all configuration
- [ ] Integration with service startup workflow
- [ ] Unit and integration tests

### ⏳ Milestone 6: Remaining services
- [ ] Complete microservices startup orchestration
- [ ] Final integration testing
- [ ] Performance optimization

## Testing Strategy

### Unit Tests
- **Coverage**: All core functionality in `internal/git/` module
- **Approach**: Table-driven tests with mock Git repositories
- **Validation**: Repository scanning, Git operations, parallel processing

### Integration Tests
- **Script-based**: Comprehensive bash script with real repository cloning
- **Go-based**: Native Go integration tests with network operations
- **Automation**: Makefile commands for easy execution
- **Real Repositories**: Tests against actual GitHub repositories

### Test Naming Convention
All test cases MUST follow a descriptive naming pattern that clearly explains the expected behavior:

**Format**: `"should [expected action/result] when [condition/scenario]"`

**Examples**:
- ✅ `"should return true when path is a valid directory"`
- ✅ `"should return error when file does not exist"`
- ✅ `"should load valid config with custom branch"`
- ✅ `"should handle uncommitted changes error"`
- ✅ `"should checkout main and pull when on different branch"`
- ❌ `"valid directory"` (too vague)
- ❌ `"test error handling"` (not descriptive)
- ❌ `"happy path"` (unclear what is being tested)

**Benefits**:
- Makes test failures immediately understandable
- Serves as living documentation of expected behavior
- Forces developers to think about what they're testing
- Easy to scan test output and understand what broke

### Test Commands
```bash
# Quick unit tests
make test

# Comprehensive integration testing
make test-integration

# Demo with real repositories
make integration-demo

# Go integration tests
go test ./test/integration/... -v
```

## Configuration

The application supports:
- **Directory Input**: Via CLI flag (`-d`) or environment variable (`SCAN_DIRECTORY`)
- **Branch Selection**: Via CLI flag (`-b`) or environment variable (`GIT_BRANCH`)
- **Environment Files**: Via CLI flag (`-e`) for loading `.env` files
- **Verbose Output**: Configurable CLI output levels
- **Parallel Processing**: Automatic concurrent repository processing

## Current Status

**Milestone 3 Complete**: The Git repository sync functionality is fully implemented and tested with:
- ✅ **Parallel Processing**: Concurrent Git operations using goroutines
- ✅ **Smart Error Handling**: Detects uncommitted changes, missing tracking branches
- ✅ **User-Friendly Output**: Clean CLI interface with colored indicators
- ✅ **Comprehensive Testing**: Unit tests + integration tests with real repositories
- ✅ **Automation**: Makefile commands for testing and demonstration

**Ready for**: Milestone 4 (Docker Compose management) implementation.