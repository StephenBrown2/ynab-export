# Contributing to YNAB Export Tool

Thank you for your interest in contributing to the YNAB Export Tool! This document provides guidelines and instructions for contributing.

## Getting Started

### Prerequisites

- Go 1.25 or later
- Git
- (Optional) [Just][just] for easier development workflow

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

```bash
# Build for current platform
go build

# Or using just
just build

# Build for all platforms
just build-all
```

### Running

```bash
# Run directly
go run .

# Or using just
just run
```

### Code Style

- Follow standard Go conventions
- Run `go fmt` before committing:

  ```bash
  go fmt ./...
  # or
  just fmt
  ```

- If you have `golangci-lint` installed:

  ```bash
  just lint
  ```

### Making Changes

1. Create a new branch for your feature or fix:

   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes

3. Test your changes thoroughly

4. Commit your changes with clear, descriptive commit messages:

   ```bash
   git commit -m "feat: add new feature"
   ```

   We follow [Conventional Commits][conventional-commits]:

   - `feat:` for new features
   - `fix:` for bug fixes
   - `docs:` for documentation changes
   - `refactor:` for code refactoring
   - `test:` for test changes
   - `chore:` for maintenance tasks

5. Push to your fork:

   ```bash
   git push origin feature/your-feature-name
   ```

6. Create a Pull Request on GitHub

## Testing

Currently, the project uses manual testing. We welcome contributions to add automated tests!

To manually test:

1. Build the application
2. Test with a real YNAB account (or create a test account)
3. Verify the export file is created correctly
4. Test on different platforms if possible

## Pull Request Guidelines

- Provide a clear description of the changes
- Reference any related issues
- Include screenshots or examples if applicable
- Ensure your code builds successfully
- Update documentation if needed

## Project Structure

```text
ynab-export/
├── main.go              # Main application with TUI logic
├── go.mod              # Go module definition
├── Justfile            # Build automation
├── .goreleaser.yml     # Release configuration
├── .github/
│   └── workflows/
│       └── release.yml # CI/CD for releases
└── README.md           # User documentation
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

By contributing, you agree that your contributions will be licensed under the MIT License.

## Thank You

Your contributions make this project better for everyone. Thank you for taking the time to contribute! 🎉

<!-- Link References -->
[just]: https://github.com/casey/just
[conventional-commits]: https://www.conventionalcommits.org/
[issues]: https://github.com/StephenBrown2/ynab-export/issues
