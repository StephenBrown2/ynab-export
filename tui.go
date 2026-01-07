package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

const (
	// YNAB API tokens are 43 characters long.
	ynabTokenLength = 43
)

var (
	titleStyle   = lipgloss.NewStyle().MarginLeft(2).Bold(true).Foreground(lipgloss.Color("205"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Bold(true)
	warningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	validStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
	fieldStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("46")) // Green for field names like nushell
)

type state int

const (
	stateValidatingToken state = iota
	stateToken
	stateFetchingBudgets
	stateBudgetSelect
	stateExporting
	stateDone
	stateError
)

// humanizeFileSize converts bytes to a human-readable format (KB, MB, GB).
func humanizeFileSize(bytes int64) string {
	const (
		kb = 1024
		mb = kb * 1024
		gb = mb * 1024
	)

	switch {
	case bytes >= gb:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(gb))
	case bytes >= mb:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(mb))
	case bytes >= kb:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(kb))
	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}

// formatMonthYear converts a date string (YYYY-MM-DD) to "Mon YYYY" format.
func formatMonthYear(dateStr string) string {
	t, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		return dateStr
	}
	return t.Format("Jan 2006")
}

// createBudgetTable creates a Nushell-style table from JSON data.budget object.
func createBudgetTable(jsonData []byte) string {
	// Extract keys and values in their original order
	keys, budget, err := extractBudgetKeysAndValues(jsonData)
	if err != nil {
		return fmt.Sprintf("Error extracting budget data: %v", err)
	}

	// Create rows for the table
	rows := make([][]string, 0, len(keys))
	for _, key := range keys {
		value := budget[key]
		inspected := inspectJSONValue(value)
		rows = append(rows, []string{key, inspected})
	}

	// Create lipgloss table with Nushell-style borders
	t := table.New().
		Border(lipgloss.RoundedBorder()).
		StyleFunc(func(_, col int) lipgloss.Style {
			// First column (field names) in green
			if col == 0 {
				return fieldStyle
			}
			// Second column (values) in default color
			return lipgloss.NewStyle()
		}).
		Rows(rows...)

	return t.Render()
}

type model struct {
	budgetList         list.Model
	err                error
	selectedBudget     budget
	token              string
	exportPath         string
	tokenValidationErr string
	budgets            []budget
	tokenInput         textinput.Model
	budgetTable        string
	summary            budgetSummary
	state              state
	tokenLengthValid   bool
	tokenSource        TokenSource
}

type budgetsFetchedMsg struct {
	err     error
	budgets []budget
}

type exportDoneMsg struct {
	err      error
	path     string
	jsonData []byte
	summary  budgetSummary
}

type tokenValidatedMsg struct {
	err   error
	token string
}

func validateTokenAsync(token string) tea.Cmd {
	return func() tea.Msg {
		if err := validateToken(token); err != nil {
			return tokenValidatedMsg{err: err}
		}
		return tokenValidatedMsg{token: token}
	}
}

func initialModel(token string, source TokenSource) model {
	ti := textinput.New()
	ti.Placeholder = "Enter your YNAB API token..."
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 70
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '•'

	// If we have a token from flag, env, or cache, validate it
	if token != "" && source != TokenSourceNone {
		return model{
			state:       stateValidatingToken,
			token:       token,
			tokenSource: source,
			tokenInput:  ti,
		}
	}

	return model{
		state:      stateToken,
		tokenInput: ti,
	}
}

func (m model) Init() tea.Cmd {
	// If we're starting in validating state, validate the token
	if m.state == stateValidatingToken {
		return validateTokenAsync(m.token)
	}
	return textinput.Blink
}

// handleKeyPress processes keyboard input based on current state.
func (m model) handleKeyPress(key string) (model, tea.Cmd) {
	switch key {
	case "ctrl+c":
		// Always allow Ctrl+C to quit
		return m, tea.Quit
	case "q":
		// Allow 'q' to quit only in done/error states
		if m.state == stateDone || m.state == stateError {
			return m, tea.Quit
		}
	case "esc":
		return m.handleEscapeKey()
	case "enter":
		return m.handleEnterKey()
	}
	return m, nil
}

// handleEscapeKey handles Esc key press.
func (m model) handleEscapeKey() (model, tea.Cmd) {
	if m.state == stateBudgetSelect {
		// If filtering is active, let the list handle Esc to clear filter
		if m.budgetList.FilterState() == list.Filtering {
			return m, nil // Let the list handle it in updateInputs
		}
		// Otherwise, go back to token entry
		m.state = stateToken
		m.tokenInput.SetValue("")
		m.tokenInput.Focus()
		return m, textinput.Blink
	}
	return m, nil
}

// handleEnterKey handles Enter key press based on state.
func (m model) handleEnterKey() (model, tea.Cmd) {
	switch m.state {
	case stateToken:
		// Only proceed if token is non-empty and has valid length
		if strings.TrimSpace(m.tokenInput.Value()) != "" && m.tokenLengthValid {
			m.token = strings.TrimSpace(m.tokenInput.Value())
			m.tokenSource = TokenSourceManual
			m.state = stateValidatingToken
			return m, validateTokenAsync(m.token)
		}
	case stateBudgetSelect:
		if selected, ok := m.budgetList.SelectedItem().(budget); ok {
			m.selectedBudget = selected
			m.state = stateExporting
			return m, func() tea.Msg { return exportBudget(m.token, selected.ID, selected.Name) }
		}
	case stateValidatingToken, stateFetchingBudgets, stateExporting:
		// No action needed for these states
	case stateDone, stateError:
		return m, tea.Quit
	}
	return m, nil
}

// handleTokenValidated processes token validation message.
func (m model) handleTokenValidated(msg tokenValidatedMsg) (model, tea.Cmd) {
	if msg.err != nil {
		// Token validation failed, show token input screen with error
		// Include the source of the token in the error message
		errorMsg := "Invalid token: " + msg.err.Error()
		if m.tokenSource != TokenSourceNone && m.tokenSource != TokenSourceManual {
			errorMsg = fmt.Sprintf("Token from %s is no longer valid: %s", m.tokenSource, msg.err.Error())
			// If the invalid token was cached, delete it (best effort, ignore errors)
			if m.tokenSource == TokenSourceCached {
				_ = DeleteCachedToken() //nolint:errcheck // Best effort cleanup, don't block on failure
			}
		}
		m.state = stateToken
		m.token = ""
		m.tokenLengthValid = false
		m.tokenSource = TokenSourceNone
		m.tokenValidationErr = errorMsg
		return m, textinput.Blink
	}

	// Token is valid, save it to cache for future use (best effort, ignore errors)
	_ = SaveCachedToken(msg.token) //nolint:errcheck // Best effort caching, don't block on failure

	// Token is valid, proceed to fetch budgets
	m.token = msg.token
	m.tokenLengthValid = false
	m.tokenValidationErr = ""
	m.state = stateFetchingBudgets
	return m, func() tea.Msg { return fetchBudgets(m.token) }
}

// handleBudgetsFetched processes budgets fetched message.
//
//nolint:unparam // Cmd is always nil, but keeping consistent signature with other handlers
func (m model) handleBudgetsFetched(msg budgetsFetchedMsg) (model, tea.Cmd) {
	if msg.err != nil {
		m.err = msg.err
		m.state = stateError
		return m, nil
	}

	m.budgets = msg.budgets
	items := make([]list.Item, len(msg.budgets))
	for i, b := range msg.budgets {
		items[i] = b
	}

	delegate := list.NewDefaultDelegate()
	m.budgetList = list.New(items, delegate, 80, 20)
	m.budgetList.Title = "Select a Budget"
	m.budgetList.SetShowStatusBar(false)
	m.budgetList.SetFilteringEnabled(true)
	m.budgetList.Styles.Title = titleStyle
	m.state = stateBudgetSelect
	return m, nil
}

// handleExportDone processes export done message.
func (m model) handleExportDone(msg exportDoneMsg) (model, tea.Cmd) {
	if msg.err != nil {
		m.err = msg.err
		m.state = stateError
		return m, nil
	}

	m.exportPath = msg.path
	m.summary = msg.summary

	// Create budget structure table
	m.budgetTable = createBudgetTable(msg.jsonData)

	m.state = stateDone
	return m, tea.Quit
}

// updateInputs updates interactive components based on state.
func (m model) updateInputs(msg tea.Msg) (model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case stateToken:
		m.tokenInput, cmd = m.tokenInput.Update(msg)
		// Validate token length as user types
		tokenValue := strings.TrimSpace(m.tokenInput.Value())
		tokenLen := len(tokenValue)

		switch {
		case tokenLen == 0:
			m.tokenLengthValid = false
			m.tokenValidationErr = ""
		case tokenLen < ynabTokenLength:
			m.tokenLengthValid = false
			m.tokenValidationErr = fmt.Sprintf("Token too short (%d/%d characters)", tokenLen, ynabTokenLength)
		case tokenLen > ynabTokenLength:
			m.tokenLengthValid = false
			m.tokenValidationErr = fmt.Sprintf("Token too long (%d/%d characters)", tokenLen, ynabTokenLength)
		default:
			// Token is the correct length (43 characters)
			m.tokenLengthValid = true
			m.tokenValidationErr = ""
		}
	case stateBudgetSelect:
		m.budgetList, cmd = m.budgetList.Update(msg)
	case stateValidatingToken, stateFetchingBudgets, stateExporting, stateDone, stateError:
		// No interactive input in these states
	}
	return m, cmd
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		newModel, cmd := m.handleKeyPress(msg.String())
		if cmd != nil {
			return newModel, cmd
		}
		m = newModel
	case tokenValidatedMsg:
		return m.handleTokenValidated(msg)
	case budgetsFetchedMsg:
		return m.handleBudgetsFetched(msg)
	case exportDoneMsg:
		return m.handleExportDone(msg)
	}

	return m.updateInputs(msg)
}

func (m model) View() string {
	var b strings.Builder

	switch m.state {
	case stateValidatingToken:
		b.WriteString(titleStyle.Render("Validating Token...") + "\n\n")
		b.WriteString("Please wait while we validate your YNAB API token.\n")

	case stateToken:
		b.WriteString(titleStyle.Render("YNAB Budget Exporter") + "\n\n")
		b.WriteString("This tool will help you export your YNAB budget for import into Actual Budget.\n\n")
		b.WriteString("To get your API token:\n")
		b.WriteString("  1. Sign in to the YNAB web app\n")
		b.WriteString("  2. Go to Account Settings → Developer Settings\n")
		b.WriteString("  3. Under 'Personal Access Tokens', click 'New Token'\n")
		b.WriteString("  4. Enter your password and click 'Generate'\n")
		b.WriteString("  5. Copy the FULL token from the top (under 'New Personal Access Token:')\n")
		b.WriteString("     NOT the partially hidden one in the table below!\n\n")
		b.WriteString(helpStyle.Render("Visit: https://app.ynab.com/settings/developer") + "\n\n")
		b.WriteString(m.tokenInput.View() + "\n")

		// Show validation feedback
		if strings.TrimSpace(m.tokenInput.Value()) != "" {
			if m.tokenLengthValid {
				b.WriteString(validStyle.Render("✓ Token length valid") + "\n")
			} else if m.tokenValidationErr != "" {
				b.WriteString(warningStyle.Render("⚠ "+m.tokenValidationErr) + "\n")
			}
		}
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("Press Enter to continue • Ctrl+C to quit"))

	case stateFetchingBudgets:
		b.WriteString(titleStyle.Render("Fetching Budgets...") + "\n\n")
		b.WriteString("Please wait while we retrieve your budgets from YNAB.\n")

	case stateBudgetSelect:
		b.WriteString(m.budgetList.View())

	case stateExporting:
		b.WriteString(titleStyle.Render("Exporting Budget...") + "\n\n")
		b.WriteString(fmt.Sprintf("Downloading budget: %s\n", m.selectedBudget.Name))
		b.WriteString("Please wait...\n")

	case stateDone:
		b.WriteString(successStyle.Render("✓ Export Complete!") + "\n\n")
		b.WriteString(fmt.Sprintf("Budget: %s\n", m.selectedBudget.Name))
		b.WriteString(fmt.Sprintf("Saved to: %s\n", m.exportPath))
		b.WriteString(fmt.Sprintf("File Size: %s\n\n", humanizeFileSize(m.summary.FileSize)))

		// Display budget structure table
		b.WriteString(titleStyle.Render("Budget Structure (data.budget):") + "\n")
		b.WriteString(m.budgetTable + "\n\n")

		b.WriteString("You can now import this file into Actual Budget:\n")
		b.WriteString("  1. Open Actual Budget\n")
		b.WriteString("  2. If a budget is already open, select the dropdown menu and 'Close File'\n")
		b.WriteString("  3. Select 'Import file'\n")
		b.WriteString("  4. Choose 'nYNAB'\n")
		b.WriteString("  5. Select the exported JSON file\n")
		b.WriteString("  6. Once imported, review your budget and follow cleanup steps at\n")
		b.WriteString("     https://actualbudget.org/docs/migration/nynab#cleanup\n")

	case stateError:
		b.WriteString(errorStyle.Render("✗ Error") + "\n\n")
		b.WriteString(fmt.Sprintf("An error occurred: %v\n\n", m.err))
		b.WriteString(helpStyle.Render("Press Enter or q to quit"))
	}

	return b.String()
}
