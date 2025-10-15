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
3. Extract and run the executable

### Option B: Build from Source

```bash
git clone https://github.com/StephenBrown2/ynab-export.git
cd ynab-export
go build
./ynab-export
```

## 3. Follow the Prompts

1. Paste your API token
2. Select your budget from the list (use `/` to filter, `Esc` to clear filter)
3. Wait for export to complete
4. Find your file in `~/Downloads/ynab-export-budget-name-TIMESTAMP.json`

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

## Keyboard Shortcuts

- `↑/↓` - Navigate
- `/` - Search/Filter
- `Enter` - Select
- `Esc` - Go Back
- `q` or `Ctrl+C` - Quit

<!-- Link References -->
[releases]: https://github.com/StephenBrown2/ynab-export/releases
[issues]: https://github.com/StephenBrown2/ynab-export/issues
