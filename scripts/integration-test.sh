#!/bin/bash

# Integration test script for go-cli
# This script creates a temporary directory, clones some popular repositories,
# runs the sync command, and validates the results

set -e  # Exit on any error

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test configuration
TEST_DIR="/tmp/go-cli-integration-test-$$"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BINARY_PATH="$PROJECT_ROOT/bin/go-cli"
LOG_FILE="$TEST_DIR/test.log"

# Repository URLs (using smaller, popular repos for faster cloning)
REPOS=(
    "https://github.com/octocat/Hello-World.git"
    "https://github.com/github/gitignore.git"
    "https://github.com/microsoft/vscode.git --depth 1"
)

REPO_NAMES=(
    "Hello-World"
    "gitignore"
    "vscode"
)

# Cleanup function
cleanup() {
    echo -e "${YELLOW}Cleaning up test environment...${NC}"
    if [ -d "$TEST_DIR" ]; then
        rm -rf "$TEST_DIR"
    fi
    echo -e "${GREEN}Cleanup completed${NC}"
}

# Trap to ensure cleanup on exit
trap cleanup EXIT

# Function to print test status
print_status() {
    local status=$1
    local message=$2
    if [ "$status" = "PASS" ]; then
        echo -e "${GREEN}✅ PASS${NC}: $message"
    elif [ "$status" = "FAIL" ]; then
        echo -e "${RED}❌ FAIL${NC}: $message"
        exit 1
    else
        echo -e "${YELLOW}ℹ️  INFO${NC}: $message"
    fi
}

# Function to validate sync results
validate_sync_results() {
    local exit_code=$1
    local output="$2"

    # Check if sync completed successfully
    if [ $exit_code -ne 0 ]; then
        print_status "FAIL" "Sync command exited with code $exit_code"
        return 1
    fi

    # Check if output contains expected success indicators
    if echo "$output" | grep -q "repositories synced successfully"; then
        print_status "PASS" "Sync completed with success message"
    else
        print_status "FAIL" "Success message not found in output"
        return 1
    fi

    # Check if all repositories were processed
    local expected_count=${#REPO_NAMES[@]}
    if echo "$output" | grep -q "Found $expected_count repositories"; then
        print_status "PASS" "Expected number of repositories found ($expected_count)"
    else
        print_status "FAIL" "Expected $expected_count repositories, but count doesn't match"
        return 1
    fi

    return 0
}

# Main test execution
main() {
    echo -e "${GREEN}🚀 Starting Go CLI Integration Tests${NC}"
    echo "Test directory: $TEST_DIR"
    echo "Binary path: $BINARY_PATH"
    echo ""

    # Check if binary exists
    if [ ! -f "$BINARY_PATH" ]; then
        print_status "FAIL" "Binary not found at $BINARY_PATH. Run 'make build' first."
        exit 1
    fi

    # Create test directory
    print_status "INFO" "Creating test directory: $TEST_DIR"
    mkdir -p "$TEST_DIR"
    cd "$TEST_DIR"

    # Clone test repositories
    print_status "INFO" "Cloning test repositories..."
    local clone_count=0

    for i in "${!REPOS[@]}"; do
        local repo_url="${REPOS[$i]}"
        local repo_name="${REPO_NAMES[$i]}"

        echo "  📂 Cloning $repo_name..."
        if git clone --quiet $repo_url "$repo_name" 2>/dev/null; then
            ((clone_count++))
            echo "     ✅ Successfully cloned $repo_name"
        else
            echo -e "     ${YELLOW}⚠️  Failed to clone $repo_name (continuing...)${NC}"
        fi
    done

    if [ $clone_count -eq 0 ]; then
        print_status "FAIL" "No repositories were successfully cloned"
        exit 1
    fi

    print_status "PASS" "Successfully cloned $clone_count repositories"

    # Run the sync command
    print_status "INFO" "Running sync command on test repositories..."
    echo "Command: $BINARY_PATH sync -d $TEST_DIR"
    echo ""

    # Capture both output and exit code
    set +e  # Don't exit on error for this command
    sync_output=$("$BINARY_PATH" sync -d "$TEST_DIR" 2>&1)
    sync_exit_code=$?
    set -e  # Re-enable exit on error

    # Display the sync output
    echo "=== Sync Command Output ==="
    echo "$sync_output"
    echo "=========================="
    echo ""

    # Validate the results
    print_status "INFO" "Validating sync results..."
    validate_sync_results $sync_exit_code "$sync_output"

    # Additional validation: check that repositories are on main/master branch
    print_status "INFO" "Validating repository states..."
    local validation_passed=true

    for repo_name in "${REPO_NAMES[@]}"; do
        if [ -d "$repo_name" ]; then
            cd "$repo_name"
            local current_branch=$(git branch --show-current 2>/dev/null || echo "unknown")
            if [[ "$current_branch" =~ ^(main|master)$ ]]; then
                echo "  ✅ $repo_name is on branch: $current_branch"
            else
                echo -e "  ${YELLOW}⚠️  $repo_name is on branch: $current_branch (not main/master)${NC}"
            fi
            cd ..
        fi
    done

    print_status "PASS" "Repository state validation completed"

    # Test summary
    echo ""
    echo -e "${GREEN}🎉 Integration Tests Summary${NC}"
    echo "  Repositories cloned: $clone_count"
    echo "  Sync exit code: $sync_exit_code"
    echo "  Test directory: $TEST_DIR"
    echo ""
    print_status "PASS" "All integration tests completed successfully!"
}

# Run the main function
main "$@"