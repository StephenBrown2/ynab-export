package main

import (
	"context"
	"encoding/json/v2"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	ynabAPIBase = "https://api.ynab.com/v1"
)

type budget struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	LastModifiedOn string `json:"last_modified_on"`
}

func (b budget) Title() string {
	if b.LastModifiedOn != "" {
		t, err := time.Parse(time.RFC3339, b.LastModifiedOn)
		if err == nil {
			return fmt.Sprintf("%s (Last Modified: %s)", b.Name, t.Format("2006-01-02"))
		}
	}
	return b.Name
}
func (b budget) Description() string { return b.ID }
func (b budget) FilterValue() string { return b.Name }

type budgetsResponse struct {
	Data struct {
		Budgets []budget `json:"budgets"`
	} `json:"data"`
}

type budgetDetailResponse struct {
	Data struct {
		Budget budgetDetail `json:"budget"`
	} `json:"data"`
}

type budgetDetail struct {
	Name           string         `json:"name"`
	FirstMonth     string         `json:"first_month"`
	LastMonth      string         `json:"last_month"`
	CurrencyFormat currencyFormat `json:"currency_format"`
	Accounts       []account      `json:"accounts"`
	Payees         []payee        `json:"payees"`
	Categories     []category     `json:"categories"`
	Transactions   []transaction  `json:"transactions"`
}

type currencyFormat struct {
	ISOCode        string `json:"iso_code"`
	CurrencySymbol string `json:"currency_symbol"`
}

type account struct {
	Closed bool `json:"closed"`
}

type payee struct {
	ID string `json:"id"`
}

type category struct {
	Deleted bool `json:"deleted"`
	Hidden  bool `json:"hidden"`
}

type transaction struct {
	ID string `json:"id"`
}

type budgetSummary struct {
	Name                 string
	Currency             string
	FirstMonth           string
	LastMonth            string
	FileSize             int64
	AccountCount         int
	ClosedAccountCount   int
	TransactionCount     int
	CategoryCount        int
	HiddenCategoryCount  int
	DeletedCategoryCount int
	PayeeCount           int
}

// validateToken checks if a token is valid by calling the /user endpoint.
func validateToken(token string) error {
	client := &http.Client{Timeout: 10 * time.Second}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ynabAPIBase+"/user", http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to validate token: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid token: API returned %s", resp.Status)
	}

	return nil
}

func fetchBudgets(token string) tea.Msg {
	client := &http.Client{Timeout: 30 * time.Second}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ynabAPIBase+"/budgets", http.NoBody)
	if err != nil {
		return budgetsFetchedMsg{err: err}
	}

	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return budgetsFetchedMsg{err: err}
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return budgetsFetchedMsg{err: fmt.Errorf("API error: %s (failed to read body: %w)", resp.Status, readErr)}
		}
		return budgetsFetchedMsg{err: fmt.Errorf("API error: %s - %s", resp.Status, string(body))}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return budgetsFetchedMsg{err: err}
	}

	var budgetsResp budgetsResponse
	if err := json.Unmarshal(body, &budgetsResp); err != nil {
		return budgetsFetchedMsg{err: err}
	}

	// Sort budgets by last modified date (most recent first)
	budgets := budgetsResp.Data.Budgets
	for i := 0; i < len(budgets)-1; i++ {
		for j := i + 1; j < len(budgets); j++ {
			ti, errI := time.Parse(time.RFC3339, budgets[i].LastModifiedOn)
			tj, errJ := time.Parse(time.RFC3339, budgets[j].LastModifiedOn)
			if errI == nil && errJ == nil && tj.After(ti) {
				budgets[i], budgets[j] = budgets[j], budgets[i]
			}
		}
	}

	return budgetsFetchedMsg{budgets: budgets}
}

// createBudgetSummary extracts summary statistics from a budget.
func createBudgetSummary(budget budgetDetail, fileSize int64) budgetSummary {
	// Count categories (non-deleted, non-hidden, hidden, deleted)
	categoryCount := 0
	hiddenCategoryCount := 0
	deletedCategoryCount := 0
	for _, cat := range budget.Categories {
		switch {
		case cat.Deleted:
			deletedCategoryCount++
		case cat.Hidden:
			hiddenCategoryCount++
		default:
			categoryCount++
		}
	}

	// Count accounts (non-closed, closed)
	accountCount := 0
	closedAccountCount := 0
	for _, acc := range budget.Accounts {
		if acc.Closed {
			closedAccountCount++
		} else {
			accountCount++
		}
	}

	// Format currency with symbol if available
	currency := budget.CurrencyFormat.ISOCode
	if budget.CurrencyFormat.CurrencySymbol != "" {
		currency = fmt.Sprintf("%s (%s)", budget.CurrencyFormat.ISOCode, budget.CurrencyFormat.CurrencySymbol)
	}

	return budgetSummary{
		Name:                 budget.Name,
		Currency:             currency,
		FirstMonth:           budget.FirstMonth,
		LastMonth:            budget.LastMonth,
		FileSize:             fileSize,
		AccountCount:         accountCount,
		ClosedAccountCount:   closedAccountCount,
		TransactionCount:     len(budget.Transactions),
		CategoryCount:        categoryCount,
		HiddenCategoryCount:  hiddenCategoryCount,
		DeletedCategoryCount: deletedCategoryCount,
		PayeeCount:           len(budget.Payees),
	}
}

func exportBudget(token, budgetID, budgetName string) tea.Msg {
	client := &http.Client{Timeout: 30 * time.Second}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	url := fmt.Sprintf("%s/budgets/%s", ynabAPIBase, budgetID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return exportDoneMsg{err: err}
	}

	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return exportDoneMsg{err: err}
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return exportDoneMsg{err: fmt.Errorf("API error: %s (failed to read body: %w)", resp.Status, readErr)}
		}
		return exportDoneMsg{err: fmt.Errorf("API error: %s - %s", resp.Status, string(body))}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return exportDoneMsg{err: err}
	}

	// Parse the budget data to extract summary information
	var budgetResp budgetDetailResponse
	if unmarshalErr := json.Unmarshal(body, &budgetResp); unmarshalErr != nil {
		return exportDoneMsg{err: unmarshalErr}
	}

	budget := budgetResp.Data.Budget
	summary := createBudgetSummary(budget, int64(len(body)))

	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return exportDoneMsg{err: err}
	}

	// Create Downloads directory path (cross-platform)
	downloadsDir := filepath.Join(homeDir, "Downloads")

	// Ensure Downloads directory exists
	if err := os.MkdirAll(downloadsDir, 0o750); err != nil {
		return exportDoneMsg{err: err}
	}

	// Create filename with timestamp and budget name
	timestamp := time.Now().Format("20060102-150405")
	// Sanitize budget name: lowercase and replace spaces with dashes
	sanitizedName := strings.ToLower(budgetName)
	sanitizedName = strings.ReplaceAll(sanitizedName, " ", "-")
	filename := fmt.Sprintf("ynab-export-%s-%s.json", sanitizedName, timestamp)
	filePath := filepath.Join(downloadsDir, filename)

	// Write the JSON to file
	if err := os.WriteFile(filePath, body, 0o600); err != nil {
		return exportDoneMsg{err: err}
	}

	return exportDoneMsg{path: filePath, summary: summary, jsonData: body}
}
