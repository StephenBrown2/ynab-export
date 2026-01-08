# Common VHS configuration for ynab-export demo recordings
# See: https://github.com/charmbracelet/vhs

# Display settings
Set FontSize 14
Set Width 900
Set Height 600
Set Padding 20

# Environment for demo mode
Env YNAB_DEMO_MODE "true"
Env YNAB_DEMO_QUIET "true"
Env YNAB_NO_CACHE "true"
Env YNAB_MOCK_DELAY_USER "1000ms"
Env YNAB_MOCK_DELAY_BUDGETS "1000ms"
Env YNAB_MOCK_DELAY_BUDGET "1500ms"
