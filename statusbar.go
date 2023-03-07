package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// statusBarMsg is a message sent to the status bar.
type statusBarMsg struct {
	Key    string
	Value  string
	Info   string
	Branch string
}

// statusBar is a status bar model.
type statusBar struct {
	width int
	style lipgloss.Style
	msg   statusBarMsg
}

// New creates a new status bar component.
func newStatusBar() *statusBar {
	return &statusBar{
		style: lipgloss.NewStyle().
			// Padding(0, 1).
			Background(lipgloss.Color("69")),
	}
}

// Init implements tea.Model.
func (s *statusBar) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (s *statusBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statusBarMsg:
		s.msg = msg
		// handle window resize
	case tea.WindowSizeMsg:
		s.width = msg.Width
	}
	return s, nil
}

// View implements tea.Model.
func (s *statusBar) View() string {
	// v := truncate.StringWithTail(s.msg.Value, uint(maxWidth-st.StatusBarValue.GetHorizontalFrameSize()), "…")

	// TODO: add a way to set the status bar style, and stuff inside it.
	return s.style.Render(strings.Repeat(" ", s.width))
}
