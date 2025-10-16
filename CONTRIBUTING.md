# Contributing to YNAB Export Tool

Thank you for your interest in contributing to the YNAB Export Tool!
This document provides guidelines and instructions for contributing.

## Getting Started

### Prerequisites

- Go 1.25 or later
- Git
- (Optional) [Just][just] for easier development workflow
- **Windows users:** Git for Windows is required for build automation
  (Just uses Git Bash)

### Technical Details

- **TUI Framework**: [Bubble Tea][bubbletea] for terminal user interface
- **JSON Handling**: Go's experimental `encoding/json/v2` for performance
  and order preservation
- **Styling**: [Lipgloss][lipgloss] for terminal styling and table rendering
- **Build Tool**: [Just][just] for task automation
- **Release Automation**: [GoReleaser][goreleaser] for cross-platform
  binary releases

### Setting Up Your Development Environment

1. Fork the repository on GitHub
2. Clone your fork:

   ```bash
   git clone https://github.com/YOUR_USERNAME/ynab-export.git
   cd ynab-export
   ```

3. Add the upstream repository:

   ```bash
   git remote add upstream https://github.com/StephenBrown2/ynab-export.git
   ```

4. Install dependencies:

   ```bash
   go mod download
   # or
   just deps
   ```

## Development Workflow

### Building

**Important:** This project requires `GOEXPERIMENT=jsonv2` to be set.
The Justfile handles this automatically, so always use `just` commands.

```bash
# Build for current platform
just build

# Build for all platforms (Linux, macOS, Windows - amd64 & arm64)
just build-all

# Clean build artifacts
just clean
```

If you must run `go` commands directly:

```bash
GOEXPERIMENT=jsonv2 go build
GOEXPERIMENT=jsonv2 go run .
GOEXPERIMENT=jsonv2 go test ./...
```

### Running

```bash
# Run with just (recommended)
just run

# Run with environment variable for token
export YNAB_API_TOKEN="your-token-here"
just run
```

### Code Style and Linting

**CRITICAL:** After every code change, always run:

```bash
just lint-fix
```

This will:

- Auto-format code with gofumpt and goimports
- Fix common linting issues automatically
- Ensure code style consistency

Before committing, run full checks:

```bash
just check
```

This runs formatting, linting, building, and tests.

### Making Changes

1. Create a new branch for your feature or fix:

   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes

3. Run `just lint-fix` after each change

4. Test your changes thoroughly

5. Commit your changes using [Conventional Commits][conventional-commits]:

   ```bash
   git commit -m "feat: add new feature"
   ```

   Commit types:

   - `feat:` for new features
   - `fix:` for bug fixes
   - `docs:` for documentation changes
   - `refactor:` for code refactoring
   - `test:` for test changes
   - `chore:` for maintenance tasks
   - `ci:` for CI/CD changes

   Examples:

   ```text
   feat: add budget structure table with order preservation
   fix: handle empty currency symbol gracefully
   docs: update README with Windows Terminal instructions
   chore: update dependencies
   ```

6. Run `just check` before committing

7. Push to your fork:

   ```bash
   git push origin feature/your-feature-name
   ```

8. Create a Pull Request on GitHub

## Project Structure

```text
ynab-export/
├── main.go              # Application entry point and CLI flag parsing
├── tui.go               # Terminal UI implementation (Bubble Tea)
├── ynab.go              # YNAB API integration and data handling
├── json.go              # Order-preserving JSON parsing utilities
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
├── Justfile             # Build automation recipes
├── .goreleaser.yml      # GoReleaser configuration for releases
├── .golangci.yml        # Linter configuration
├── .markdownlint.json   # Markdown linter configuration
├── .github/
│   ├── workflows/
│   │   └── release.yml  # GitHub Actions release workflow
│   └── ISSUE_TEMPLATE/  # Issue templates
├── README.md            # User-facing documentation
├── QUICKSTART.md        # Quick start guide
├── CONTRIBUTING.md      # This file
└── AGENTS.md            # AI agent instructions
```

### Key Files

#### main.go

- Entry point with CLI flag parsing
- Handles `--version` flag
- Checks for `YNAB_API_TOKEN` environment variable
- Launches Terminal UI (Bubble Tea)

#### tui.go

- Implements Bubble Tea state machine for terminal interface
- States: token validation, budget selection, exporting, done, error
- Handles all user interaction and display

#### ynab.go

- YNAB API client implementation
- Budget fetching and exporting
- Data structure definitions
- Budget summary calculations

#### json.go

- Order-preserving JSON parsing using `encoding/json/v2`
- `OrderedObject` and `ObjectMember` types
- `extractBudgetKeysAndValues()` preserves original JSON key order
- `inspectJSONValue()` provides Nushell-style type descriptions
- Date formatting utilities

## Testing

Currently, the project uses manual testing.
We welcome contributions to add automated tests!

### Manual Testing

```bash
# Run with environment variable
export YNAB_API_TOKEN="your-token-here"
just run

# Test specific platform builds
just build-linux
just build-darwin
just build-windows
```

### Testing Different Scenarios

1. **Token validation**: Test with invalid/expired tokens
2. **Budget selection**: Test filtering and navigation
3. **Export**: Verify JSON file is created correctly
4. **Error handling**: Test network failures, API errors

To manually test:

1. Build the application
2. Test with a real YNAB account (or create a test account)
3. Verify the export file is created correctly
4. Test on different platforms if possible

## Pull Request Guidelines

- Provide a clear description of the changes
- Reference any related issues
- Include screenshots or examples if applicable
- Ensure your code builds successfully with `just check`
- Update documentation if needed

## Release Process

Releases are automated via GitHub Actions when tags are pushed:

```bash
# Create and push a release tag
just release v0.0.5

# Or manually:
git tag -a v0.0.5 -m "Release v0.0.5"
git push origin v0.0.5
```

This triggers:

1. GoReleaser builds binaries for all platforms
2. GitHub Release is created with changelog
3. Binaries are uploaded as release assets

## Common Development Tasks

### Adding a New Feature

1. Create a feature branch: `git checkout -b feat/my-feature`
2. Make your changes
3. Run `just lint-fix` after each change
4. Run `just check` before committing
5. Commit with conventional commit message
6. Create a pull request

### Fixing a Bug

1. Create a bugfix branch: `git checkout -b fix/issue-description`
2. Add a test that reproduces the bug (if applicable)
3. Fix the bug
4. Run `just check`
5. Commit with conventional commit message
6. Create a pull request

### Updating Dependencies

```bash
# Update all dependencies
go get -u ./...
go mod tidy

# Update a specific dependency
go get -u github.com/charmbracelet/bubbletea

# Verify everything still works
just check
```

## Troubleshooting Development Issues

### Build Failures

If you see `encoding/json/v2` errors:

```text
build constraints exclude all Go files in .../encoding/json/v2
```

**Solution**: Use `just build` (not `go build`) or set `GOEXPERIMENT=jsonv2`

### Linting Failures

```bash
# Auto-fix formatting and common issues
just lint-fix

# Check what remains
just lint
```

## Areas for Contribution

Here are some areas where contributions would be especially welcome:

### Features

- Add configuration file support
- Add support for selecting multiple budgets
- Add progress bars for large exports
- Add GUI version (using Fyne, Wails, or similar)

### Test Coverage

- Add unit tests
- Add integration tests
- Add end-to-end tests

### Documentation

- Improve README
- Add more examples
- Create video tutorials
- Translate documentation

### Build/CI

- Improve GitHub Actions workflow
- Add code coverage reporting
- Add automated security scanning

## Questions or Problems?

- Check existing [Issues][issues]
- Create a new issue if needed
- Be respectful and constructive

## License

By contributing, you agree that your contributions will be licensed
under the MIT License.

## Thank You

Your contributions make this project better for everyone.
Thank you for taking the time to contribute! 🎉

<!-- Link References -->

[bubbletea]: https://github.com/charmbracelet/bubbletea
[lipgloss]: https://github.com/charmbracelet/lipgloss
[just]: https://github.com/casey/just
[goreleaser]: https://goreleaser.com/
[conventional-commits]: https://www.conventionalcommits.org/
[issues]: https://github.com/StephenBrown2/ynab-export/issues
