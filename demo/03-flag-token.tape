# VHS Tape for recording ynab-export demo: Command-line Flag Token
# See: https://github.com/charmbracelet/vhs

Source demo/config.tape
Output demo/03-flag-token.gif

# Run program with token flag
Type "./ynab-export --token demo-ynab-token-for-recording-1234567890abc"
Enter

# Wait for automatic token detection and budget list
Sleep 2s

# Navigate budgets
Down
Sleep 500ms
Down
Sleep 500ms
Down
Sleep 500ms

# Select budget
Enter

# Wait for export
Sleep 10s
