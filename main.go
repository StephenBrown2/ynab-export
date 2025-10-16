package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// Version information (set by goreleaser).
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	// Define command-line flags

	showVersion := flag.Bool("version", false, "show version information")

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

	// Launch TUI
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
