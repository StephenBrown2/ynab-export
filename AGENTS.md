# Agent Instructions for ynab-export

This document provides guidance for LLM agents working in this repository.

## Repository Overview

This is a Go 1.25+ project that exports YNAB budgets for import into Actual Budget. It uses:

- **Go Modules** for dependency management
- **Just** for build automation and task running
- **golangci-lint** for linting and formatting
- **GoReleaser** for automated releases
- **GitHub Actions** for CI/CD
- **Experimental Go features**: `GOEXPERIMENT=jsonv2` (required for building)

## Essential Commands

### Always Use `just` Commands

This project uses [Just](https://github.com/casey/just) as the task runner. **Always prefer `just` commands** over direct `go` commands, as they include necessary environment setup (like `GOEXPERIMENT=jsonv2`).

View all available commands:

```bash
just
```

### After Every Code Change

**CRITICAL**: After making any code changes, always run:

```bash
just lint-fix
```

This will:

- Auto-format code with gofumpt and goimports
- Fix common linting issues automatically
- Ensure code style consistency

### Common Development Workflow

1. **Make code changes** to `.go` files
2. **Run lint-fix** (ALWAYS):

   ```bash
   just lint-fix
   ```

3. **Build and test**:

   ```bash
   just build
   ```

4. **Run the application**:

   ```bash
   just run
   ```

5. **Run full checks** before committing:

   ```bash
   just check
   ```

### Build Commands

```bash
# Build for current platform
just build

# Build for all platforms (Linux, macOS, Windows - amd64 & arm64)
just build-all

# Clean build artifacts
just clean
```

### Linting Commands

```bash
# Auto-fix issues (USE THIS AFTER EVERY CHANGE)
just lint-fix

# Check formatting only
just fmt-check

# Run full linter (checks without fixing)
just lint

# Format code manually
just fmt
```

### Testing Commands

```bash
# Run tests
just test

# Run all checks (format, lint, test)
just check
```

## Important Environment Variables

### GOEXPERIMENT

This project **requires** `GOEXPERIMENT=jsonv2` to be set. The Justfile handles this automatically, which is why you should always use `just` commands.

If you must run `go` commands directly:

```bash
GOEXPERIMENT=jsonv2 go build
GOEXPERIMENT=jsonv2 go run .
GOEXPERIMENT=jsonv2 go test ./...
```

## Code Style Guidelines

### Formatting

- Use **gofumpt** (stricter than gofmt)
- Use **goimports** for import ordering
- **Run `just lint-fix` after every change**

### Linting

This project uses golangci-lint with comprehensive rules including:

- `govet`, `errcheck`, `staticcheck`
- `gosimple`, `ineffassign`, `unused`
- `gofumpt`, `goimports`, `misspell`
- `revive`, `stylecheck`, `unconvert`
- And many more (see `.golangci.yml`)

### Commit Messages

Use **Conventional Commits** format:

```text
type(scope): subject

body

footer
```

Types:

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `chore`: Maintenance tasks
- `refactor`: Code refactoring
- `test`: Test additions or changes
- `ci`: CI/CD changes

Examples:

```text
feat: add last modified date to budget list
fix: handle empty currency symbol gracefully
docs: update README with new screenshots
chore: update dependencies
```

### Git Configuration

This repository uses **Commit signing**: Enabled (SSH)

## Project Structure

```text
ynab-export/
├── main.go              # Application entry point with UI selection
├── tui.go               # Terminal UI implementation (Bubble Tea)
├── gui.go               # Graphical UI implementation (Fyne)
├── ynab.go              # YNAB API integration and data handling
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
├── Justfile             # Task runner recipes (USE THIS!)
├── .golangci.yml        # Linter configuration
├── .goreleaser.yml      # Release configuration
├── .github/
│   ├── workflows/
│   │   └── release.yml  # GitHub Actions release workflow
│   └── ISSUE_TEMPLATE/  # Issue templates
├── README.md            # User documentation
├── QUICKSTART.md        # Quick start guide
├── CONTRIBUTING.md      # Contribution guidelines
├── DEPLOYMENT.md        # Deployment instructions (internal)
└── AGENTS.md            # This file
```

## File-Specific Notes

### main.go

- Entry point with CLI flag parsing
- Handles `--version`, `--tui`, and `--gui` flags
- Checks for `YNAB_API_TOKEN` environment variable
- Detects terminal context and selects appropriate UI
- Initializes TUI or GUI based on context/flags

### tui.go

- Implements Bubble Tea state machine for terminal interface
- States: token validation, budget selection, exporting, done, error
- Handles all user interaction and display in terminal mode

### gui.go

- Implements Fyne graphical interface
- Mirrors TUI flow: token entry, budget selection, export, completion
- Uses synchronous wrapper functions from ynab.go
- Suitable for double-click execution

### ynab.go

- YNAB API client implementation
- Budget fetching and exporting
- Data structure definitions
- Budget summary calculations
- Provides both async (Bubble Tea messages) and sync (return values) functions

## Release Process

Releases are automated via GitHub Actions when tags are pushed:

```bash
# Create and push a release tag
just release v0.0.2

# Or manually:
git tag -a v0.0.2 -m "Release v0.0.2"
git push origin v0.0.2
```

This triggers:

1. GoReleaser builds binaries for all platforms
2. GitHub Release is created
3. Binaries are uploaded as release assets

## Common Pitfalls

### ❌ Don't Do This

```bash
# Don't run go commands directly (missing GOEXPERIMENT)
go build
go run .

# Don't skip lint-fix
# (Your code won't match project style)

# Don't commit without running checks
git commit -am "changes"
```

### ✅ Do This Instead

```bash
# Use just commands
just build
just run

# Always run lint-fix after changes
just lint-fix

# Run checks before committing
just check
git add .
git commit -m "feat: descriptive message"
```

## Debugging

### Build Failures

If you see `encoding/json/v2` errors:

```text
build constraints exclude all Go files in .../encoding/json/v2
```

Solution: Use `just build` (not `go build`)

### Linting Failures

Run auto-fix first:

```bash
just lint-fix
```

Then check what remains:

```bash
just lint
```

### Format Check Failures

```bash
# Auto-fix formatting
just fmt

# Or use lint-fix which does both
just lint-fix
```

## Testing

### Testing the UI Modes

This project has two UI modes that should be tested:

**Terminal UI (TUI)**:

```bash
# Run in terminal
just run

# Or force TUI mode
./ynab-export --tui
```

**Graphical UI (GUI)**:

```bash
# Force GUI mode (in terminal)
./ynab-export --gui

# Or double-click the binary (auto-detects GUI)
# On Linux/macOS: chmod +x ynab-export first
```

**Auto-detection**:

- When run from terminal: defaults to TUI
- When double-clicked: defaults to GUI
- Use `--tui` or `--gui` flags to override

### GUI-Specific Testing Notes

The GUI implementation uses Fyne v2.6.3 and should be tested on multiple platforms:

- **Linux**: Requires X11/Wayland display server
- **macOS**: Should work natively
- **Windows**: Should work natively

Common GUI issues to check:

- Window sizing and layout
- Token entry field behavior
- Budget list selection
- Progress dialogs during export
- Completion screen display

## Best Practices for Agents

1. **Always use `just` commands** - They include necessary environment setup
2. **Run `just lint-fix` after every code change** - Non-negotiable
3. **Run `just check` before suggesting commits** - Ensures quality
4. **Use conventional commit messages** - Keeps history clean
5. **Test build after changes** - Run `just build` to verify
6. **Test both UI modes** - Verify TUI and GUI work correctly
7. **Update documentation** - Keep README.md in sync with code changes
8. **Check for errors** - Use VS Code's problem panel or `just lint`

## Quick Reference Card

```bash
# Make changes → Lint → Build → Test → Commit
just lint-fix  # After EVERY change
just build     # Verify it compiles
just test      # Run tests (if any)
just check     # Full validation
git add .
git commit -m "type: message"

# View all commands
just

# Get help on a specific recipe
just --show <recipe>
```

## Questions or Issues?

- Check the [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines
- Review the [README.md](README.md) for project documentation
- Check the [Justfile](Justfile) for available commands
- Run `just` to see all available recipes

## Summary

**Remember the golden rule**: After every code change, run `just lint-fix` before doing anything else. This ensures your code matches the project's style and catches common issues automatically.

Happy coding! 🚀
