# YNAB Export Tool

A simple, beginner-friendly tool to export your YNAB (You Need A Budget) budget data for import into [Actual Budget][actual-budget].

This tool simplifies the export process described in the [Actual Budget migration guide][actual-migration] by providing an interactive terminal interface that guides you through each step.

## Features

- 🎯 **Beginner-friendly**: Interactive prompts guide you through the entire process
- �️ **Dual UI modes**: Graphical interface (GUI) or Terminal interface (TUI)
- �🔒 **Secure**: Your API token is masked as you type and never stored
- 🚀 **Fast**: Direct API integration with YNAB
- 💾 **Automatic saving**: Exports are automatically saved to your Downloads folder
- 🌐 **Cross-platform**: Works on Windows, macOS, and Linux
- 📦 **Zero dependencies**: Single binary with no installation required

## Installation

### Download Pre-built Binary

Download the latest release for your operating system from the [Releases page][releases].

Choose the appropriate binary:

- **Windows**: `ynab-export_*_windows_amd64.exe`
- **macOS (Intel)**: `ynab-export_*_darwin_amd64`
- **macOS (Apple Silicon)**: `ynab-export_*_darwin_arm64`
- **Linux (64-bit)**: `ynab-export_*_linux_amd64`
- **Linux (ARM64)**: `ynab-export_*_linux_arm64`

The binaries are ready to run - no extraction or installation needed!

### Build from Source

If you have Go installed (version 1.25 or later), you can build from source:

> **Note:** This project uses Go's experimental `encoding/json/v2` package. You'll need to set the `GOEXPERIMENT` environment variable.

```bash
GOEXPERIMENT=jsonv2 go install github.com/StephenBrown2/ynab-export@latest
```

Or clone and build:

```bash
git clone https://github.com/StephenBrown2/ynab-export.git
cd ynab-export
GOEXPERIMENT=jsonv2 go build
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

Simply run the downloaded binary:

**Linux/macOS:**

```bash
# Make it executable (first time only)
chmod +x ynab-export_*_linux_amd64  # or darwin_amd64/darwin_arm64

# Run it (GUI launches automatically when double-clicked)
./ynab-export_*_linux_amd64

# Or force terminal mode:
./ynab-export_*_linux_amd64 --tui
```

**Windows:**

Double-click the `.exe` file for GUI mode, or run from Command Prompt:

```cmd
ynab-export_0.0.1_windows_amd64.exe

# Or force terminal mode:
ynab-export_0.0.1_windows_amd64.exe --tui
```

> **Tip:** You can rename the binary to simply `ynab-export` (or `ynab-export.exe` on Windows) for easier use.

#### UI Mode Selection

The tool automatically chooses the appropriate interface:

- **GUI mode** (default): Launches when you double-click the executable or run it outside a terminal
- **TUI mode** (terminal): Launches automatically when run from a terminal/console

You can explicitly select a mode with command-line flags:

```bash
# Force graphical interface
./ynab-export --gui

# Force terminal interface
./ynab-export --tui
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

### 1. Welcome Screen (Token Entry)

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

### 2. Token Validation

```text
┌────────────────────────────────────────────────────┐
│ YNAB Budget Exporter                               │
│                                                    │
│ ✓ Token validated successfully                     │
│                                                    │
│ Fetching your budgets...                          │
└────────────────────────────────────────────────────┘
```

### 3. Budget Selection

```text
┌────────────────────────────────────────────────────┐
│ YNAB Budget Exporter                               │
│                                                    │
│ Select a budget to export:                        │
│                                                    │
│ > Personal Budget (Last Modified: 2025-10-14)     │
│   f1a2b3c4-d5e6-7f8g-9h0i-1j2k3l4m5n6o             │
│                                                    │
│   Family Budget (Last Modified: 2025-10-10)       │
│   a1b2c3d4-e5f6-7g8h-9i0j-1k2l3m4n5o6p             │
│                                                    │
│   Business Budget (Last Modified: 2025-09-28)     │
│   z9y8x7w6-v5u4-t3s2-r1q0-p9o8n7m6l5k4             │
│                                                    │
│ Use ↑/↓ to navigate • / to search • Enter to      │
│ select • Esc to go back • q/Ctrl+C to quit        │
└────────────────────────────────────────────────────┘
```

### 4. Export in Progress

```text
┌────────────────────────────────────────────────────┐
│ YNAB Budget Exporter                               │
│                                                    │
│ Exporting Budget...                               │
│                                                    │
│ Downloading budget: Personal Budget               │
│ Please wait...                                    │
└────────────────────────────────────────────────────┘
```

### 5. Export Complete with Budget Summary

```text
┌────────────────────────────────────────────────────┐
│ ✓ Export Complete!                                │
│                                                    │
│ Budget: Personal Budget                           │
│ Saved to: ~/Downloads/ynab-export-personal-       │
│           budget-20251015-143022.json             │
│                                                    │
│ Budget Summary:                                   │
│   Currency:     USD ($)                           │
│   Accounts:     8 (plus 2 closed)                 │
│   Categories:   24 (plus 3 hidden, 1 deleted)     │
│   Payees:       142                               │
│   Transactions: 1,847                             │
│   Date Range:   Jan 2023 to Oct 2025              │
│   File Size:    2.3 MB                            │
│                                                    │
│ You can now import this file into Actual Budget:  │
│   1. Open Actual Budget                           │
│   2. If a budget is already open, select the      │
│      dropdown menu and 'Close File'               │
│   3. Select 'Import file'                         │
│   4. Choose 'nYNAB'                               │
│   5. Select the exported JSON file                │
│   6. Once imported, review your budget and        │
│      follow cleanup steps at                      │
│      actualbudget.org/docs/migration/nynab#cleanup│
│                                                    │
│ Press Enter or q to quit                          │
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
# Replace with your actual binary name, e.g., ynab-export_0.0.1_darwin_arm64
xattr -d com.apple.quarantine ynab-export_*_darwin_*

# Or make it executable
chmod +x ynab-export_*_darwin_*
```

Or right-click the file, select "Open", and click "Open" in the security dialog.

## Development

### Prerequisites

- Go 1.25 or later
- (Optional) [Just][just] for build automation
- **Windows users:** Git for Windows is required for build automation (Just uses Git Bash)

### Technical Details

- **GUI Framework**: [Gio UI][gio] - Pure Go, no CGO dependencies for easy cross-compilation
- **TUI Framework**: [Bubble Tea][bubbletea] - Elegant terminal user interface
- **JSON Handling**: Go's experimental `encoding/json/v2` for improved performance

The GUI uses Gio UI, which provides excellent cross-platform support without requiring CGO. This means:

- Simple cross-compilation (just set GOOS/GOARCH)
- No platform-specific C dependencies
- Smaller binary sizes
- Easier to build and distribute

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
- [Charm Bracelet][charm] for the beautiful TUI libraries (Bubble Tea)
- [Gio UI][gio] for the pure Go GUI framework

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
[gio]: https://gioui.org/
[bubbletea]: https://github.com/charmbracelet/bubbletea
