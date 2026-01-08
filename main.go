package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/StephenBrown2/ynab-export/internal/mockserver"
	tea "github.com/charmbracelet/bubbletea"
)

// Version information (set by goreleaser).
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

const envTrue = "true"

func main() {
	// Define command-line flags
	showVersion := flag.Bool("version", false, "show version information")
	tokenFlag := flag.String("token", "", "YNAB API token (overrides environment variable and cached token)")

	// Short flag aliases
	flag.BoolVar(showVersion, "v", false, "show version information (shorthand)")
	flag.StringVar(tokenFlag, "t", "", "YNAB API token (shorthand)")

	flag.Parse()

	// Check for version flag
	if *showVersion {
		fmt.Fprintf(os.Stderr, "ynab-export version %s\n", version)
		fmt.Fprintf(os.Stderr, "commit: %s\n", commit)
		fmt.Fprintf(os.Stderr, "built: %s\n", date)
		os.Exit(0)
	}

	// Check for demo mode
	var shutdownMock func()
	if os.Getenv("YNAB_DEMO_MODE") == envTrue {
		serverURL, shutdown := mockserver.StartMockServer()
		shutdownMock = shutdown
		ynabAPIBase = serverURL + "/v1"
		if os.Getenv("YNAB_DEMO_QUIET") != envTrue {
			fmt.Fprintf(os.Stderr, "Demo mode enabled. Using mock YNAB API at %s\n", serverURL)
		}
	}
	fmt.Fprintf(os.Stderr, "\n") // Separate TUI output from prompt

	// Determine token and its source (priority: flag > env > cached)
	token, source := resolveToken(*tokenFlag)

	// Launch TUI and run
	exitCode := runTUI(token, source)

	// Cleanup mock server if running
	if shutdownMock != nil {
		shutdownMock()
	}

	if exitCode != 0 {
		os.Exit(exitCode)
	}
}

// runTUI launches the terminal UI and returns exit code.
func runTUI(token string, source TokenSource) int {
	p := tea.NewProgram(initialModel(token, source))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}
	return 0
}

// resolveToken determines the token to use and its source.
// Priority: flag > environment variable > cached token.
func resolveToken(flagToken string) (string, TokenSource) {
	// Priority 1: Command-line flag
	if flagToken != "" {
		fmt.Fprintf(os.Stderr, "Using API token from command-line flag.\n\n")
		return flagToken, TokenSourceFlag
	}

	// Priority 2: Environment variable
	if envToken := os.Getenv("YNAB_API_TOKEN"); envToken != "" {
		fmt.Fprintf(os.Stderr, "Using API token from environment variable.\n\n")
		return envToken, TokenSourceEnv
	}

	// Priority 3: Cached token (unless caching is disabled)
	if os.Getenv("YNAB_NO_CACHE") != envTrue {
		cachedToken, err := LoadCachedToken()
		if err != nil {
			// Warn user about cache read failure
			fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
			fmt.Fprintf(os.Stderr, "You will need to enter your token manually.\n\n")
		} else if cachedToken != "" {
			fmt.Fprintf(os.Stderr, "Using cached API token from %s.\n\n", GetTokenCacheLocation())
			return cachedToken, TokenSourceCached
		}
	}

	return "", TokenSourceNone
}
