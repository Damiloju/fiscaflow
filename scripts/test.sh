#!/bin/bash

# Test runner script for FiscaFlow
# This script runs different types of tests with coverage reporting

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    print_error "go.mod not found. Please run this script from the project root."
    exit 1
fi

# Check if required tools are installed
if ! command_exists go; then
    print_error "Go is not installed or not in PATH"
    exit 1
fi

# Function to run tests with coverage
run_tests_with_coverage() {
    local test_type=$1
    local coverage_file="coverage_${test_type}.out"
    local coverage_html="coverage_${test_type}.html"

    print_status "Running ${test_type} tests with coverage..."

    # Run tests with coverage
    go test -v -coverprofile="$coverage_file" -covermode=atomic ./...

    if [ $? -eq 0 ]; then
        print_success "${test_type} tests passed"

        # Generate coverage report
        if command_exists go tool cover; then
            go tool cover -html="$coverage_file" -o "$coverage_html"
            print_status "Coverage report generated: $coverage_html"

            # Show coverage summary
            go tool cover -func="$coverage_file" | tail -1
        fi
    else
        print_error "${test_type} tests failed"
        return 1
    fi
}

# Function to run specific test packages
run_package_tests() {
    local package=$1
    local test_type=$2

    print_status "Running ${test_type} tests for package: $package"

    go test -v ./$package

    if [ $? -eq 0 ]; then
        print_success "${test_type} tests for $package passed"
    else
        print_error "${test_type} tests for $package failed"
        return 1
    fi
}

# Function to run benchmarks
run_benchmarks() {
    print_status "Running benchmarks..."

    go test -bench=. -benchmem ./...

    if [ $? -eq 0 ]; then
        print_success "Benchmarks completed"
    else
        print_error "Benchmarks failed"
        return 1
    fi
}

# Function to run race detection
run_race_detection() {
    print_status "Running tests with race detection..."

    go test -race ./...

    if [ $? -eq 0 ]; then
        print_success "Race detection tests passed"
    else
        print_error "Race detection tests failed"
        return 1
    fi
}

# Function to run linting
run_linting() {
    print_status "Running linting..."

    if command_exists golangci-lint; then
        golangci-lint run
        if [ $? -eq 0 ]; then
            print_success "Linting passed"
        else
            print_error "Linting failed"
            return 1
        fi
    else
        print_warning "golangci-lint not found, skipping linting"
        print_status "Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
    fi
}

# Function to run security scanning
run_security_scan() {
    print_status "Running security scan..."

    if command_exists gosec; then
        gosec ./...
        if [ $? -eq 0 ]; then
            print_success "Security scan passed"
        else
            print_warning "Security scan found issues (check output above)"
        fi
    else
        print_warning "gosec not found, skipping security scan"
        print_status "Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"
    fi
}

# Function to clean up coverage files
cleanup_coverage() {
    print_status "Cleaning up coverage files..."
    rm -f coverage_*.out coverage_*.html
    print_success "Coverage files cleaned up"
}

# Main execution
main() {
    local test_type=${1:-all}

    print_status "Starting test suite for FiscaFlow..."
    print_status "Test type: $test_type"

    # Clean up any existing coverage files
    cleanup_coverage

    case $test_type in
    "unit")
        print_status "Running unit tests only..."
        run_package_tests "internal/domain/user" "unit"
        run_package_tests "internal/api/handlers" "unit"
        ;;
    "integration")
        print_status "Running integration tests only..."
        run_package_tests "tests/integration" "integration"
        ;;
    "coverage")
        print_status "Running tests with coverage..."
        run_tests_with_coverage "all"
        ;;
    "benchmark")
        print_status "Running benchmarks..."
        run_benchmarks
        ;;
    "race")
        print_status "Running race detection..."
        run_race_detection
        ;;
    "lint")
        print_status "Running linting..."
        run_linting
        ;;
    "security")
        print_status "Running security scan..."
        run_security_scan
        ;;
    "all")
        print_status "Running all tests and checks..."

        # Run linting first
        run_linting

        # Run security scan
        run_security_scan

        # Run unit tests
        print_status "Running unit tests..."
        run_package_tests "internal/domain/user" "unit"
        run_package_tests "internal/api/handlers" "unit"

        # Run integration tests
        print_status "Running integration tests..."
        run_package_tests "tests/integration" "integration"

        # Run race detection
        run_race_detection

        # Run benchmarks
        run_benchmarks

        # Run coverage
        run_tests_with_coverage "all"
        ;;
    *)
        print_error "Unknown test type: $test_type"
        echo "Available test types:"
        echo "  unit      - Run unit tests only"
        echo "  integration - Run integration tests only"
        echo "  coverage  - Run tests with coverage"
        echo "  benchmark - Run benchmarks"
        echo "  race      - Run race detection"
        echo "  lint      - Run linting"
        echo "  security  - Run security scan"
        echo "  all       - Run all tests and checks (default)"
        exit 1
        ;;
    esac

    print_success "Test suite completed successfully!"
}

# Handle script arguments
if [ "$1" = "--help" ] || [ "$1" = "-h" ]; then
    echo "FiscaFlow Test Runner"
    echo ""
    echo "Usage: $0 [test_type]"
    echo ""
    echo "Test types:"
    echo "  unit         - Run unit tests only"
    echo "  integration  - Run integration tests only"
    echo "  coverage     - Run tests with coverage"
    echo "  benchmark    - Run benchmarks"
    echo "  race         - Run race detection"
    echo "  lint         - Run linting"
    echo "  security     - Run security scan"
    echo "  all          - Run all tests and checks (default)"
    echo ""
    echo "Examples:"
    echo "  $0              # Run all tests"
    echo "  $0 unit         # Run unit tests only"
    echo "  $0 coverage     # Run tests with coverage"
    echo "  $0 --help       # Show this help"
    exit 0
fi

# Run main function with all arguments
main "$@"
