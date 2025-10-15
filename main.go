package main

import (
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
	// Check for version flag
	if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		fmt.Fprintf(os.Stderr, "ynab-export version %s\n", version)
		fmt.Fprintf(os.Stderr, "commit: %s\n", commit)
		fmt.Fprintf(os.Stderr, "built: %s\n", date)
		os.Exit(0)
	}

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
