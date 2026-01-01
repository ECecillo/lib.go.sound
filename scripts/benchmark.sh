#!/bin/bash
# Benchmark runner and comparison script

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
BENCHMARKS_DIR="$REPO_ROOT/benchmarks"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Create benchmarks directory if it doesn't exist
mkdir -p "$BENCHMARKS_DIR"

# Function to print colored output
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to run benchmarks
run_benchmarks() {
    local output_file="$1"
    local benchtime="${2:-1s}"

    print_info "Running benchmarks (benchtime=$benchtime)..."
    cd "$REPO_ROOT"

    go test ./... -bench=. -benchmem -benchtime="$benchtime" | tee "$output_file"

    print_info "Benchmarks saved to: $output_file"
}

# Function to save baseline
save_baseline() {
    local timestamp=$(date +%Y%m%d_%H%M%S)
    local baseline_file="$BENCHMARKS_DIR/baseline_$timestamp.txt"

    print_info "Saving baseline benchmarks..."
    run_benchmarks "$baseline_file" "2s"

    # Create symlink to latest baseline
    ln -sf "baseline_$timestamp.txt" "$BENCHMARKS_DIR/baseline_latest.txt"

    print_info "Baseline saved as: baseline_$timestamp.txt"
    print_info "Latest baseline: baseline_latest.txt"
}

# Function to compare with baseline
compare_with_baseline() {
    local baseline_file="$BENCHMARKS_DIR/baseline_latest.txt"

    if [ ! -f "$baseline_file" ]; then
        print_error "No baseline found. Run: $0 baseline"
        exit 1
    fi

    local new_file="$BENCHMARKS_DIR/comparison_$(date +%Y%m%d_%H%M%S).txt"

    print_info "Running benchmarks for comparison..."
    run_benchmarks "$new_file" "2s"

    # Check if benchstat is installed
    if ! command -v benchstat &> /dev/null; then
        print_warning "benchstat not found. Install with:"
        print_warning "  go install golang.org/x/perf/cmd/benchstat@latest"
        print_info ""
        print_info "Manual comparison:"
        print_info "  Baseline: $baseline_file"
        print_info "  New:      $new_file"
        exit 0
    fi

    print_info ""
    print_info "Comparison Results:"
    print_info "===================="
    benchstat "$baseline_file" "$new_file"
}

# Function to run quick benchmarks (for CI/development)
quick_bench() {
    print_info "Running quick benchmarks (100ms each)..."
    cd "$REPO_ROOT"
    go test ./... -bench=. -benchmem -benchtime=100ms
}

# Function to run comprehensive benchmarks
comprehensive_bench() {
    local timestamp=$(date +%Y%m%d_%H%M%S)
    local output_file="$BENCHMARKS_DIR/comprehensive_$timestamp.txt"

    print_info "Running comprehensive benchmarks (5s each, this will take a while)..."
    run_benchmarks "$output_file" "5s"
}

# Function to show usage
show_usage() {
    cat << EOF
Benchmark Runner and Comparison Tool

Usage: $0 <command>

Commands:
    baseline        Save current performance as baseline (2s per benchmark)
    compare         Compare current performance against baseline
    quick           Quick benchmark check (100ms per benchmark)
    comprehensive   Comprehensive benchmark suite (5s per benchmark)
    help            Show this help message

Examples:
    # Save baseline before making changes
    $0 baseline

    # After making changes, compare performance
    $0 compare

    # Quick check during development
    $0 quick

    # Comprehensive benchmarking for release
    $0 comprehensive

Benchmark files are saved to: $BENCHMARKS_DIR

Note: Install benchstat for statistical comparison:
    go install golang.org/x/perf/cmd/benchstat@latest
EOF
}

# Main script
main() {
    local command="${1:-help}"

    case "$command" in
        baseline)
            save_baseline
            ;;
        compare)
            compare_with_baseline
            ;;
        quick)
            quick_bench
            ;;
        comprehensive)
            comprehensive_bench
            ;;
        help|--help|-h)
            show_usage
            ;;
        *)
            print_error "Unknown command: $command"
            echo ""
            show_usage
            exit 1
            ;;
    esac
}

main "$@"
