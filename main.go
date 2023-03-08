package main

// A simple program demonstrating the text area component from the Bubbles
// component library.

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// These variables are updated at compile time.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// TODO: make these configurable, maybe HCL?
var (
	halStyleColor = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
)

func main() {
	// Print version information if the user asks for it with "-v / --version / version".
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "-v" || arg == "--version" || arg == "version" {
			fmt.Println(version+"-"+commit, date)
			os.Exit(0)
		}
	}

	p := tea.NewProgram(
		newModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
