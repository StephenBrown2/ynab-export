package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

// Version information (set by goreleaser).
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	// Define command-line flags
	var (
		showVersion = flag.Bool("version", false, "show version information")
		useTUI      = flag.Bool("tui", false, "force terminal UI mode")
		useGUI      = flag.Bool("gui", false, "force graphical UI mode")
	)

	// Short flag aliases
	flag.BoolVar(showVersion, "v", false, "show version information (shorthand)")

	flag.Parse()

	// Check for version flag
	if *showVersion {
		fmt.Fprintf(os.Stderr, "ynab-export version %s\n", version)
		fmt.Fprintf(os.Stderr, "commit: %s\n", commit)
		fmt.Fprintf(os.Stderr, "built: %s\n", date)
		os.Exit(0)
	}

	// Determine which UI to use
	if *useTUI && *useGUI {
		fmt.Fprintln(os.Stderr, "Error: Cannot specify both --tui and --gui")
		os.Exit(1)
	}

	var forceGUI bool
	switch {
	case *useGUI:
		forceGUI = true
	case *useTUI:
		forceGUI = false
	default:
		// Auto-detect: use GUI if not running in a terminal
		forceGUI = !term.IsTerminal(int(os.Stdin.Fd()))
	}

	// Launch appropriate UI
	if forceGUI {
		runGUI()
	} else {
		p := tea.NewProgram(initialModel())
		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}
}
