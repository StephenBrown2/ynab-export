# VHS Tape for recording ynab-export demo: Manual Token Entry
# See: https://github.com/charmbracelet/vhs

Source demo/config.tape
Output demo/01-manual-token.gif

# Run program
Type "./ynab-export"
Enter
Sleep 1s

# Wait for token entry screen
Sleep 500ms

# Type demo token (43 characters)
Type@150ms "demo-ynab-token-for-recording-1234567890abc"
Sleep 1s
Enter

# Wait for token validation and budget list
Sleep 2s

# Navigate budgets
Down
Sleep 500ms
Down
Sleep 500ms
Up
Sleep 500ms

# Select budget
Enter

# Wait for export
Sleep 10s
