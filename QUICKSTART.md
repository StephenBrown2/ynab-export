# Quick Start Guide

Get started with YNAB Export Tool in 5 minutes!

## 1. Get Your YNAB API Token

Visit: <https://app.ynab.com/settings/developer>

1. Click "New Token"
2. Enter your password
3. Copy the FULL token from the top (NOT the obfuscated one in the table!)

## 2. Download & Run

### Option A: Download Pre-built Binary

1. Go to [Releases][releases]
2. Download for your OS
3. Run the executable from a terminal

### Option B: Build from Source

```bash
git clone https://github.com/StephenBrown2/ynab-export.git
cd ynab-export
GOEXPERIMENT=jsonv2 go build
./ynab-export
```

## 3. Follow the Prompts

1. Paste your API token (it will be saved for future use)
2. Select your budget from the list (use `/` to filter, `Esc` to clear filter)
3. Wait for export to complete
4. Find your file in `~/Downloads/ynab-export-budget-name-TIMESTAMP.json`

> **Tip:** After the first run, your token is cached locally. Future runs will skip the token entry step!

## 4. Import to Actual Budget

1. Open Actual Budget
2. Close File → Import file
3. Select "nYnab"
4. Choose your exported JSON

## Done! 🎉

Your YNAB budget is now in Actual Budget.

## Need Help?

- Read the full [README](README.md)
- Check [Troubleshooting](README.md#troubleshooting)
- Open an [Issue][issues]

## Command-Line Options

```bash
./ynab-export [options]

  -t, --token    Provide API token directly (overrides cached/env token)
  -v, --version  Show version information
```

## Token Priority

The tool looks for your token in this order:

1. Command-line flag (`-t` or `--token`)
2. Environment variable (`YNAB_API_TOKEN`)
3. Cached token (from previous run)
4. Manual entry (prompted in the app)

## Keyboard Shortcuts

- `↑/↓` - Navigate
- `/` - Search/Filter
- `Enter` - Select
- `Esc` - Go Back / Clear Filter
- `q` or `Ctrl+C` - Quit

<!-- Link References -->
[releases]: https://github.com/StephenBrown2/ynab-export/releases
[issues]: https://github.com/StephenBrown2/ynab-export/issues
