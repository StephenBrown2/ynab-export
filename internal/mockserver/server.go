// Package mockserver provides a mock YNAB API server for demo/testing.
//
//nolint:revive // Method names must match the generated ServerInterface (e.g., GetBudgetById not GetBudgetByID)
package mockserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

// MockServer implements the YNAB API ServerInterface for demo/testing purposes.
type MockServer struct {
	generator *Generator
	delays    DelayConfig
}

// DelayConfig holds per-endpoint delay configuration.
type DelayConfig struct {
	User    time.Duration // Delay for /user endpoint
	Budgets time.Duration // Delay for /budgets endpoint
	Budget  time.Duration // Delay for /budgets/{id} endpoint
	Default time.Duration // Default delay for other endpoints
}

// LoadDelayConfig loads delay configuration from environment variables.
func LoadDelayConfig() DelayConfig {
	return DelayConfig{
		User:    parseDuration(os.Getenv("YNAB_MOCK_DELAY_USER"), parseDuration(os.Getenv("YNAB_MOCK_DELAY"), 0)),
		Budgets: parseDuration(os.Getenv("YNAB_MOCK_DELAY_BUDGETS"), parseDuration(os.Getenv("YNAB_MOCK_DELAY"), 0)),
		Budget:  parseDuration(os.Getenv("YNAB_MOCK_DELAY_BUDGET"), parseDuration(os.Getenv("YNAB_MOCK_DELAY"), 0)),
		Default: parseDuration(os.Getenv("YNAB_MOCK_DELAY"), 0),
	}
}

func parseDuration(s string, defaultVal time.Duration) time.Duration {
	if s == "" {
		return defaultVal
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return defaultVal
	}
	return d
}

// NewMockServer creates a new mock server with the given configuration.
func NewMockServer(config MockConfig, delays DelayConfig) *MockServer {
	gen := NewGenerator(config)
	// Pre-generate budgets so they're consistent across requests
	gen.GenerateBudgets()

	return &MockServer{
		generator: gen,
		delays:    delays,
	}
}

// StartMockServer starts the mock server on a random port and returns the URL and shutdown function.
func StartMockServer() (serverURL string, shutdown func()) {
	config := DefaultMockConfig()
	delays := LoadDelayConfig()
	mockServer := NewMockServer(config, delays)

	// Create handler with /v1 base URL to match YNAB API
	handler := HandlerFromMuxWithBaseURL(mockServer, http.NewServeMux(), "/v1")

	// Listen on random port using ListenConfig for proper context support
	var lc net.ListenConfig
	listener, err := lc.Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		panic(fmt.Sprintf("failed to start mock server: %v", err))
	}

	server := &http.Server{Handler: handler}
	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "Mock server error: %v\n", err)
		}
	}()

	serverURL = "http://" + listener.Addr().String()
	shutdown = func() {
		if err := server.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing mock server: %v\n", err)
		}
	}

	return serverURL, shutdown
}

// writeJSON writes a JSON response.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GetUser implements the /user endpoint.
func (s *MockServer) GetUser(w http.ResponseWriter, _ *http.Request) {
	time.Sleep(s.delays.User)

	user := s.generator.GenerateUser()
	resp := UserResponse{
		Data: struct {
			User User `json:"user"`
		}{
			User: user,
		},
	}
	writeJSON(w, http.StatusOK, resp)
}

// GetBudgets implements the /budgets endpoint.
func (s *MockServer) GetBudgets(w http.ResponseWriter, _ *http.Request, _ GetBudgetsParams) {
	time.Sleep(s.delays.Budgets)

	budgets := s.generator.GenerateBudgets()
	resp := map[string]any{
		"data": map[string]any{
			"budgets":        budgets,
			"default_budget": nil,
		},
	}
	writeJSON(w, http.StatusOK, resp)
}

// GetBudgetById implements the /budgets/{budget_id} endpoint.
func (s *MockServer) GetBudgetById(w http.ResponseWriter, _ *http.Request, budgetId string, _ GetBudgetByIdParams) {
	time.Sleep(s.delays.Budget)

	detail := s.generator.GenerateBudgetDetail(budgetId)
	if detail == nil {
		http.Error(w, "Budget not found", http.StatusNotFound)
		return
	}

	resp := BudgetDetailResponse{
		Data: struct {
			Budget          BudgetDetail `json:"budget"`
			ServerKnowledge int64        `json:"server_knowledge"`
		}{
			Budget:          *detail,
			ServerKnowledge: 1,
		},
	}
	writeJSON(w, http.StatusOK, resp)
}
