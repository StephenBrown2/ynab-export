# VHS Tape for recording ynab-export demo: Environment Variable Token
# See: https://github.com/charmbracelet/vhs

Source demo/config.tape
Output demo/02-env-token.gif

# Show how to set the environment variable
Type "export YNAB_API_TOKEN='your-ynab-api-token-here...'"
Sleep 1s
Enter
Sleep 500ms

# Reset token to valid 43 character string
Env YNAB_API_TOKEN "demo-ynab-token-for-recording-1234567890abc"

# Run program
Type "./ynab-export"
Enter

# Wait for automatic token detection and budget list
Sleep 2s

# Navigate budgets
Down
Sleep 1s
Up
Sleep 1s

# Select budget
Enter

# Wait for export
Sleep 10s
