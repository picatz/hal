package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

func EditorTextArea() textarea.Model {
	editor := textarea.New()
	editor.Placeholder = "What do you want to do?"
	editor.Prompt = halStyleColor.Bold(true).Render("│")
	editor.CharLimit = 4096
	editor.SetWidth(80)
	editor.Focus()
	editor.Focused()
	editor.ShowLineNumbers = true
	editor.FocusedStyle.CursorLine = lipgloss.NewStyle().
		Background(lipgloss.Color("236")). // faint background with
		Foreground(lipgloss.Color("231"))  // extra bright text on cursor line

	return editor
}
