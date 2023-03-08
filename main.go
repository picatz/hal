package main

// A simple program demonstrating the text area component from the Bubbles
// component library.

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TODO: make these configurable, maybe HCL?
var (
	halStyleColor = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
)

func main() {
	p := tea.NewProgram(
		newModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
