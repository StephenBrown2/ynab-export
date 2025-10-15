# YNAB Export Tool

A simple, beginner-friendly tool to export your YNAB (You Need A Budget) budget data for import into [Actual Budget][actual-budget].

This tool simplifies the export process described in the [Actual Budget migration guide][actual-migration] by providing an interactive terminal interface that guides you through each step.

## Features

- 🎯 **Beginner-friendly**: Interactive prompts guide you through the entire process
- 🔒 **Secure**: Your API token is masked as you type and never stored
- 🚀 **Fast**: Direct API integration with YNAB
- 💾 **Automatic saving**: Exports are automatically saved to your Downloads folder
- 🖥️ **Cross-platform**: Works on Windows, macOS, and Linux
- 📦 **Zero dependencies**: Single binary with no installation required

## Installation

### Download Pre-built Binary

Download the latest release for your operating system from the [Releases page][releases].

Choose the appropriate file:

- **Windows**: `ynab-export_*_Windows_x86_64.zip`
- **macOS (Intel)**: `ynab-export_*_Darwin_x86_64.tar.gz`
- **macOS (Apple Silicon)**: `ynab-export_*_Darwin_arm64.tar.gz`
- **Linux (64-bit)**: `ynab-export_*_Linux_x86_64.tar.gz`
- **Linux (ARM64)**: `ynab-export_*_Linux_arm64.tar.gz`

Extract the archive and run the executable.

### Build from Source

If you have Go installed (version 1.25 or later), you can build from source:

```bash
go install github.com/StephenBrown2/ynab-export@latest
```

Or clone and build:

```bash
git clone https://github.com/StephenBrown2/ynab-export.git
cd ynab-export
go build
```

### Using Just (Build Tool)

If you have [Just][just] installed:

```bash
# Build for current platform
just build

# Build for all platforms
just build-all

# Run directly
just run

# See all available commands
just
```

## Usage

### Step 1: Get Your YNAB API Token

Before running the tool, you'll need a YNAB Personal Access Token:

1. Sign in to the [YNAB web app][ynab-app]
2. Go to **Account Settings** → **Developer Settings**
3. Under "Personal Access Tokens", click **New Token**
4. Enter your password and click **Generate**
5. **Important**: Copy the FULL token from the top of the page (under "New Personal Access Token:")
   - Do NOT use the partially obfuscated token shown in the table (e.g., `XXXXXXXXXX-Wax0q8`)
   - The token is only shown once, so copy it immediately!

Direct link: [YNAB Developer Settings][ynab-developer]

### Step 2: Run the Tool

Simply run the executable:

```bash
./ynab-export
```

On Windows, double-click `ynab-export.exe` or run from Command Prompt:

```cmd
ynab-export.exe
```

### Step 3: Follow the Prompts

The tool will guide you through:

1. **Enter your API token** (the token you generated in Step 1)
   - **Tip**: Set the `YNAB_API_TOKEN` environment variable to skip entering your token each time
2. **Select your budget** from the list of budgets in your YNAB account
3. **Wait for export** - the tool downloads your budget data
4. **Done!** Your budget is saved to `~/Downloads/ynab-export-budget-name-YYYYMMDD-HHMMSS.json`

#### Using Environment Variable for Token (Optional)

To skip entering your token each time, set the `YNAB_API_TOKEN` environment variable:

```bash
# Linux/macOS
export YNAB_API_TOKEN="your-token-here"
./ynab-export

# Windows (PowerShell)
$env:YNAB_API_TOKEN="your-token-here"
.\ynab-export.exe

# Windows (Command Prompt)
set YNAB_API_TOKEN=your-token-here
ynab-export.exe
```

The tool will automatically validate the token from the environment variable and skip the token entry screen if valid.

### Step 4: Import into Actual Budget

Now that you have your exported JSON file:

1. Open **Actual Budget**
2. Select the dropdown menu and choose **Close File**
3. Click **Import file**
4. Select **nYnab** as the import type
5. Choose the exported JSON file from your Downloads folder
6. Follow any cleanup steps mentioned in the [Actual Budget migration guide][actual-migration-cleanup]

## Screenshots

```text
┌────────────────────────────────────────────────────┐
│ YNAB Budget Exporter                               │
│                                                    │
│ This tool will help you export your YNAB budget   │
│ for import into Actual Budget.                    │
│                                                    │
│ To get your API token:                            │
│   1. Sign in to the YNAB web app                  │
│   2. Go to Account Settings → Developer Settings  │
│   3. Under 'Personal Access Tokens', click 'New'  │
│   4. Enter your password and click 'Generate'     │
│                                                    │
│ Enter your YNAB API token: ••••••••••••••••••••   │
│                                                    │
│ Press Enter to continue • Ctrl+C to quit          │
└────────────────────────────────────────────────────┘
```

## Keyboard Shortcuts

- **Arrow Keys** (↑/↓): Navigate through budget list
- **/** : Filter/search budgets
- **Enter**: Select/Confirm
- **Esc**: Go back to previous screen
- **Ctrl+C** or **q**: Quit the application

## Troubleshooting

### "API error: 401 Unauthorized"

Your API token is invalid or expired. Generate a new token from YNAB's Developer Settings.

### "No budgets found"

Make sure you have at least one budget in your YNAB account.

### "Permission denied" when saving file

Check that you have write permissions to your Downloads folder.

### Binary won't run on macOS

macOS may block the binary because it's not from an identified developer. To run it:

```bash
xattr -d com.apple.quarantine ynab-export
```

Or right-click the file, select "Open", and click "Open" in the security dialog.

## Development

### Prerequisites

- Go 1.25 or later
- (Optional) [Just][just] for build automation
- **Windows users:** Git for Windows is required for build automation (Just uses Git Bash)

### Building

```bash
# Install dependencies
go mod download

# Build
go build -o ynab-export

# Run tests (when added)
go test -v ./...

# Cross-compile for all platforms
just build-all
```

### Project Structure

```text
ynab-export/
├── main.go              # Main application code
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
├── Justfile            # Build automation recipes
├── .goreleaser.yml     # GoReleaser configuration
├── .github/
│   └── workflows/
│       └── release.yml # GitHub Actions release workflow
└── README.md           # This file
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [YNAB][ynab] for their excellent budgeting software and API
- [Actual Budget][actual-budget] for the open-source alternative
- [Charm Bracelet][charm] for the beautiful TUI libraries

## Related Links

- [YNAB API Documentation][ynab-api]
- [Actual Budget Migration Guide][actual-migration]
- [Actual Budget Website][actual-budget]

## Support

If you encounter any issues or have questions:

1. Check the [Troubleshooting](#troubleshooting) section above
2. Review the [Issues][issues] page
3. Create a new issue if your problem isn't already listed

---

**Note**: This tool is not affiliated with YNAB or Actual Budget. It's a community project to help users migrate their data.

<!-- Link References -->
[actual-budget]: https://actualbudget.org/
[actual-migration]: https://actualbudget.org/docs/migration/nynab
[actual-migration-cleanup]: https://actualbudget.org/docs/migration/nynab#cleanup
[releases]: https://github.com/StephenBrown2/ynab-export/releases
[issues]: https://github.com/StephenBrown2/ynab-export/issues
[just]: https://github.com/casey/just
[ynab]: https://www.ynab.com/
[ynab-app]: https://app.ynab.com
[ynab-developer]: https://app.ynab.com/settings/developer
[ynab-api]: https://api.ynab.com/
[charm]: https://charm.sh/
