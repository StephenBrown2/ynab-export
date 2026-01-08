# VHS Tape for recording ynab-export main demo GIF
# See: https://github.com/charmbracelet/vhs
#
# This tape demonstrates the full workflow with manual token entry.
# For other token entry methods, see:
#   - 01-manual-token.tape
#   - 02-env-token.tape
#   - 03-flag-token.tape

Set PlaybackSpeed 0.8
Source demo/config.tape
Output demo/demo.gif
Output demo/demo.txt

# Run program
Type "./ynab-export"
Enter

# Wait for token entry screen
Sleep 1500ms

Screenshot demo/demo-1-token.png

# Type demo token (43 characters)
Type "demo-ynab-token"
Screenshot demo/demo-2-invalid-token.png

Type "-for-recording-1234567890abc"
Sleep 1s
Screenshot demo/demo-3-valid-token.png

Enter

# Wait for token validation and budget list
Sleep 2s

# Navigate to second budget
Down
Sleep 1s

Screenshot demo/demo-4-selected-budget.png

Up
Sleep 1s

Down
Sleep 1s

Right
Sleep 1s

Left
Sleep 1s

# Select budget
Enter

# Wait for export
Sleep 5s

Screenshot demo/demo-5-export.png
Sleep 1s
