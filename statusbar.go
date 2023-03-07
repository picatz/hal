package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/ansi"
)

// statusBarMsg is a message sent to the status bar.
type statusBarMsg struct {
	ChatThread *chatThread
}

// statusBar is a status bar model.
type statusBar struct {
	width      int
	style      lipgloss.Style
	chatThread *chatThread
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
	case tea.WindowSizeMsg:
		s.width = msg.Width
	case *statusBarMsg:
		if msg.ChatThread != nil {
			s.chatThread = msg.ChatThread
		}
	}
	return s, nil
}

// View implements tea.Model.
func (s *statusBar) View() string {
	// v := truncate.StringWithTail(s.msg.Value, uint(maxWidth-st.StatusBarValue.GetHorizontalFrameSize()), "…")

	var (
		currentThreadName string = "*"
		chatMessageCount  int
		chatTokensCount   int
	)

	if s.chatThread != nil {
		currentThreadName = s.chatThread.Name
		chatMessageCount = len(s.chatThread.ChatHistory)
		chatTokensCount = s.chatThread.Tokens
	}

	messageCountStatusBarBlock := lipgloss.NewStyle().Background(lipgloss.Color("63")).Render(fmt.Sprintf(" Messages: %d ", chatMessageCount))

	tokensCountStatusBarBlock := lipgloss.NewStyle().Background(lipgloss.Color("62")).Render(fmt.Sprintf(" Tokens: %d ", chatTokensCount))

	currentThreadNameBlock := lipgloss.NewStyle().Background(lipgloss.Color("69")).Bold(true).Render(fmt.Sprintf(" %s ", currentThreadName))

	// get printable characters (non ANSI escape codes)
	printableChars := ansi.PrintableRuneWidth(messageCountStatusBarBlock + tokensCountStatusBarBlock + currentThreadNameBlock)

	// build status bar including the current thread name on right hand side, filling the rest of the space with spaces
	statusText := strings.Repeat(" ", s.width-printableChars-2) + messageCountStatusBarBlock + tokensCountStatusBarBlock + currentThreadNameBlock + " "

	// TODO: add a way to set the status bar style, and stuff inside it.
	return s.style.Render(statusText)
}
