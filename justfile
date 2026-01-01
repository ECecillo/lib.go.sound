# lib.go.sound - Modern command runner with Just
# Run `just` or `just --list` to see all available commands

# Variables
data_path := "data"
ffmpeg := "~/Downloads/ffmpeg"
ffplay := "~/Downloads/ffplay"

# Default recipe (runs when you type `just`)
default:
    @just --list

# ==============================================================================
# Application Commands
# ==============================================================================

# Run the application
run:
    @echo "Running application..."
    go run cmd/main.go

# Encode output.bin to output.wav
encode:
    @echo "Encoding output.bin to output.wav..."
    {{ffmpeg}} -f s16le -ar 44100 -ac 1 -i {{data_path}}/output.bin {{data_path}}/output.wav

# Play the generated audio file
play:
    {{ffplay}} -i {{data_path}}/output.wav

# Play with waveform visualization
play-with-wave:
    {{ffplay}} -showmode 1 {{data_path}}/output.wav

# Generate audio, encode to WAV, and play
all: run encode play

# ==============================================================================
# Testing Commands
# ==============================================================================

# Run all tests (unit + golden file tests)
test:
    @echo "Running all tests..."
    go test ./... -v

# Run tests with coverage report
test-coverage:
    @echo "Running tests with coverage..."
    go test ./... -cover
    @echo ""
    @echo "Detailed coverage by package:"
    go test ./pkg/format/... -cover
    go test ./pkg/sine/... -cover

# Run unit tests only (exclude golden and fuzz)
test-unit:
    @echo "Running unit tests..."
    go test ./... -v -run "Test[^G]" -short

# Run golden file tests
test-golden:
    @echo "Running golden file tests..."
    go test ./pkg/sine/... -v -run TestGoldenFiles

# Update golden reference files (use after intentional changes)
test-golden-update:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "⚠️  WARNING: This will update golden reference files!"
    echo "Only do this if you've verified the changes are correct."
    read -p "Continue? [y/N] " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "Updating golden files..."
        go test ./pkg/sine/... -run TestGoldenFiles -update-golden -v
        echo "✓ Golden files updated"
    else
        echo "Cancelled."
    fi

# Run fuzz tests (30s each)
test-fuzz:
    @echo "Running fuzz tests (30s each)..."
    go test ./pkg/... -fuzz=. -fuzztime=20s

# Run all tests including coverage, golden files, and fuzz
test-all: test-coverage test-golden

# ==============================================================================
# Benchmark Commands
# ==============================================================================

# Run all benchmarks
bench:
    @echo "Running all benchmarks..."
    go test ./... -bench=. -benchmem

# Quick benchmark check (100ms per benchmark)
bench-quick:
    @echo "Running quick benchmarks..."
    ./scripts/benchmark.sh quick

# Save current performance as baseline
bench-baseline:
    @echo "Saving performance baseline..."
    ./scripts/benchmark.sh baseline

# Compare current performance with baseline
bench-compare:
    @echo "Comparing with baseline..."
    ./scripts/benchmark.sh compare

# Run comprehensive benchmark suite (5s each)
bench-comprehensive:
    @echo "Running comprehensive benchmarks..."
    ./scripts/benchmark.sh comprehensive

# Save benchmark results with timestamp
bench-save:
    #!/usr/bin/env bash
    set -euo pipefail
    mkdir -p benchmarks
    echo "Saving benchmark results..."
    go test ./... -bench=. -benchmem | tee benchmarks/bench_$(date +%Y%m%d_%H%M%S).txt
    echo "Results saved to benchmarks/"

# ==============================================================================
# Pre-commit / CI Commands
# ==============================================================================

# Run pre-commit checks (tests + coverage + quick fuzz + quick bench)
pre-commit:
    @echo "Running pre-commit checks..."
    @echo ""
    @echo "1/4 Running tests..."
    go test ./...
    @echo ""
    @echo "2/4 Checking coverage..."
    go test ./... -cover
    @echo ""
    @echo "3/4 Quick fuzz tests..."
    go test ./pkg/format/... -fuzz=FuzzPCM16_ConvertSample -fuzztime=5s
    go test ./pkg/sine/... -fuzz=FuzzSineGeneration -fuzztime=5s
    @echo ""
    @echo "4/4 Quick benchmarks..."
    ./scripts/benchmark.sh quick
    @echo ""
    @echo "✓ All pre-commit checks passed!"

# Fast CI pipeline (< 1 minute)
ci-fast:
    @echo "Running fast CI pipeline..."
    go test ./... -cover -race

# Comprehensive CI pipeline (~10 minutes)
ci-comprehensive:
    @echo "Running comprehensive CI pipeline..."
    @echo "This will take about 10 minutes..."
    go test ./... -bench=. -benchmem -benchtime=2s
    go test ./pkg/format/... -fuzz=. -fuzztime=30s
    go test ./pkg/sine/... -fuzz=. -fuzztime=30s

# ==============================================================================
# Development Commands
# ==============================================================================

# Format code with gofmt
fmt:
    @echo "Formatting code..."
    gofmt -s -w .
    @echo "✓ Code formatted"

# Run linter
lint:
    @echo "Running linter..."
    golangci-lint run

# Tidy dependencies
tidy:
    @echo "Tidying dependencies..."
    go mod tidy
    @echo "✓ Dependencies tidied"

# Verify dependencies
verify:
    @echo "Verifying dependencies..."
    go mod verify
    @echo "✓ Dependencies verified"

# Download dependencies
deps:
    @echo "Downloading dependencies..."
    go mod download
    @echo "✓ Dependencies downloaded"

# Build the application
build:
    @echo "Building application..."
    go build -v ./cmd/...
    @echo "✓ Build successful"

# Clean build artifacts and test cache
clean:
    @echo "Cleaning build artifacts and test cache..."
    go clean -cache -testcache -modcache
    rm -rf {{data_path}}/*.wav {{data_path}}/*.bin
    @echo "✓ Cleaned"

# ==============================================================================
# Installation Commands
# ==============================================================================

# Install required tools
install-tools:
    @echo "Installing required tools..."
    go install golang.org/x/perf/cmd/benchstat@latest
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    @echo "✓ Tools installed"

# Setup the project (install deps and tools)
setup: deps install-tools
    @echo "✓ Project setup complete"

# ==============================================================================
# Documentation Commands
# ==============================================================================

# Generate documentation
docs:
    @echo "Generating documentation..."
    go doc ./...

# Serve documentation locally
docs-serve:
    @echo "Serving documentation at http://localhost:6060"
    godoc -http=:6060

# ==============================================================================
# Quick Workflows
# ==============================================================================

# Quick development check (test + lint + build)
check: test lint build
    @echo "✓ All checks passed!"

# Full local verification (same as pre-commit)
verify-all: pre-commit
    @echo "✓ Full verification complete!"

# Watch tests (requires entr: brew install entr)
watch-test:
    @echo "Watching for changes..."
    find . -name "*.go" | entr -c just test

# Watch and run application
watch-run:
    @echo "Watching for changes..."
    find . -name "*.go" | entr -c just run
