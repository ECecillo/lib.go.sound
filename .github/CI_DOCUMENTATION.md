# CI/CD Documentation

This document explains the GitHub Actions workflows and CI/CD setup for this repository.

## Workflows Overview

### 1. Main CI Pipeline (`ci.yml`)

**Triggers:** Push to main/master/develop, Pull Requests

**Jobs:**

#### Test & Coverage
- Runs all unit tests with race detection
- Generates coverage report
- Uploads to Codecov
- Fails if coverage drops below 90%

#### Golden File Tests
- Runs snapshot tests to ensure output consistency
- Detects any unintended changes to audio generation

#### Fuzz Tests
- Runs quick fuzz tests (10 seconds each)
- Tests format conversions and sine generation
- Finds edge cases automatically

#### Benchmarks
- Runs performance benchmarks
- Comments results on Pull Requests
- Stores results as artifacts

#### Lint & Format
- Runs golangci-lint with comprehensive checks
- Verifies code formatting with gofmt
- Ensures code quality standards

#### Build
- Verifies code compiles successfully
- Builds all packages and commands

#### Test Matrix
- Tests across Go versions 1.21, 1.22, 1.23
- Ensures backward compatibility

#### CI Success
- Summary job that requires all checks to pass
- Single status check for branch protection

**Duration:** ~5-10 minutes

### 2. Nightly Comprehensive Tests (`nightly.yml`)

**Triggers:** Daily at 2 AM UTC, Manual

**Jobs:**

#### Comprehensive Test Suite
- Extended test runs with race detection
- Full coverage analysis
- 2-second benchmark time

#### Extended Fuzz Testing
- Matrix of all 9 fuzz tests
- Runs for 2 minutes each
- Uploads fuzz corpus as artifacts

#### Comprehensive Benchmarks
- 5-second benchmark time for accuracy
- Compares with baseline if available
- Stores results for trend analysis

#### Memory Profiling
- Generates memory profiles
- Analyzes allocation patterns
- Uploads profiles for review

#### Race Detection
- Runs tests 10 times with race detector
- Catches concurrency issues

#### Notify on Failure
- Creates GitHub issue if any job fails
- Includes run details and links

**Duration:** ~30 minutes

### 3. Benchmark PR (`benchmark-pr.yml`)

**Triggers:** Pull Request open/sync/reopen

**Jobs:**

#### Benchmark Comparison
- Runs benchmarks on PR branch
- Runs benchmarks on base branch
- Uses benchstat for statistical comparison
- Comments comparison on PR
- Updates comment on new commits
- Warns about performance regressions

**Features:**
- Shows Δ% change for each benchmark
- Identifies statistically significant changes
- Provides legend for interpretation

**Duration:** ~3-5 minutes

### 4. Release (`release.yml`)

**Triggers:** Push tags matching `v*.*.*`

**Jobs:**

#### Create Release
- Runs full test suite
- Generates benchmark results
- Creates changelog from commits
- Generates release notes with:
  - Changes since last release
  - Test coverage
  - Performance benchmarks
- Creates GitHub Release
- Saves benchmark baseline
- Commits baseline to repository

**Duration:** ~5-10 minutes

## Configuration Files

### `.golangci.yml`

Configures linting rules:
- Enables 18 linters for code quality
- Customized rules for style and security
- Excludes some checks for test files
- Ensures consistent code style

### `.github/dependabot.yml`

Configures automatic dependency updates:
- Checks Go modules weekly
- Checks GitHub Actions weekly
- Opens PRs for updates
- Labels PRs appropriately

## Status Badges

Add these to your README.md:

```markdown
[![CI](https://github.com/ECecillo/lib.go.sound/actions/workflows/ci.yml/badge.svg)](https://github.com/ECecillo/lib.go.sound/actions/workflows/ci.yml)
[![Nightly](https://github.com/ECecillo/lib.go.sound/actions/workflows/nightly.yml/badge.svg)](https://github.com/ECecillo/lib.go.sound/actions/workflows/nightly.yml)
[![codecov](https://codecov.io/gh/ECecillo/lib.go.sound/branch/main/graph/badge.svg)](https://codecov.io/gh/ECecillo/lib.go.sound)
[![Go Report Card](https://goreportcard.com/badge/github.com/ECecillo/lib.go.sound)](https://goreportcard.com/report/github.com/ECecillo/lib.go.sound)
```

## Branch Protection Rules

Recommended settings for main branch:

1. **Require status checks to pass:**
   - CI Success
   - Test & Coverage
   - Golden File Tests
   - Fuzz Tests
   - Lint & Format

2. **Require branches to be up to date**

3. **Require pull request reviews:**
   - 1 approval required

4. **Dismiss stale reviews**

5. **Require linear history**

## Secrets Required

### Optional Secrets

1. **CODECOV_TOKEN**
   - For Codecov integration
   - Get from https://codecov.io
   - Add in repository settings → Secrets

## Workflow Permissions

The workflows use these permissions:

- **ci.yml:** read (default)
- **nightly.yml:** issues: write, contents: read
- **benchmark-pr.yml:** pull-requests: write, contents: read
- **release.yml:** contents: write

## Artifacts Stored

### CI Pipeline
- `benchmark-results` - Benchmark output from each run

### Nightly
- `fuzz-corpus-*` - Fuzz test corpus for each function
- `benchmark-comprehensive-*` - Extended benchmark results
- `memory-profile-*` - Memory profiling data

### Benchmark PR
- None (results posted as comments)

### Release
- `release_benchmarks.txt` - Attached to release

## Triggering Workflows Manually

### Nightly Tests
```bash
gh workflow run nightly.yml
```

### View Workflow Status
```bash
gh run list --workflow=ci.yml
gh run watch
```

### Download Artifacts
```bash
gh run download <run-id>
```

## Cost Optimization

The workflows are optimized for cost:

1. **Caching:** Go modules and build cache
2. **Timeouts:** All jobs have reasonable timeouts
3. **Conditional runs:** Some jobs only run on specific events
4. **Matrix strategy:** Parallel execution where possible

**Estimated GitHub Actions minutes per month:**
- CI Pipeline: ~200-300 minutes
- Nightly: ~900 minutes (30 min × 30 days)
- Benchmark PR: ~50-100 minutes
- **Total: ~1,200-1,400 minutes/month**

Free tier includes 2,000 minutes/month for public repos.

## Troubleshooting

### Tests Fail in CI but Pass Locally

1. Check Go version matches
2. Run with race detector: `go test -race ./...`
3. Check for non-deterministic tests
4. Review CI logs for specific errors

### Coverage Drops Below Threshold

1. Add tests for new code
2. Run `make test-coverage` locally
3. Identify untested code
4. Adjust threshold in `ci.yml` if justified

### Benchmark PR Comment Not Appearing

1. Check workflow permissions
2. Verify GITHUB_TOKEN has write access
3. Check workflow logs for errors

### Dependabot PRs Failing

1. Review the dependency update
2. Check if API changes affect code
3. Update code if needed
4. Run tests locally with new version

## Maintenance

### Weekly
- Review failed nightly runs
- Check Dependabot PRs
- Monitor coverage trends

### Monthly
- Review benchmark trends
- Update baseline if performance improved
- Clean up old artifacts

### Quarterly
- Review and update linter rules
- Update Go versions in matrix
- Review and optimize workflows

## Future Enhancements

Potential additions:

1. **Code coverage trends:** Track coverage over time
2. **Performance trends:** Graph benchmark results
3. **Security scanning:** Add CodeQL or Snyk
4. **Docker builds:** Build and publish containers
5. **Documentation generation:** Auto-generate docs
6. **Release automation:** Automated version bumping
