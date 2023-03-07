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
	width  int
	height int
	style  lipgloss.Style
	msg    statusBarMsg
}

// New creates a new status bar component.
func newStatusBar() *statusBar {
	s := &statusBar{
		// width: w,
		// height: h,

		style: lipgloss.NewStyle().
			Padding(0, 1).
			Background(lipgloss.Color("235")).
			Foreground(lipgloss.Color("243")),
	}
	return s
}

// SetSize implements common.Component.
func (s *statusBar) SetSize(width, height int) {
	s.width = width
	s.height = height
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
		s.height = msg.Height
	}
	return s, nil
}

// View implements tea.Model.
func (s *statusBar) View() string {
	// st := s.style
	// w := lipgloss.Width
	// maxWidth := s.width // - w(key) - w(info) - w(branch) - w(help)
	// v := truncate.StringWithTail(s.msg.Value, uint(maxWidth-st.StatusBarValue.GetHorizontalFrameSize()), "…")
	// value := st.StatusBarValue.
	// 	Width(maxWidth).
	// 	Render(v)

	return lipgloss.NewStyle().
		MaxWidth(s.width).
		MaxHeight(s.height).
		Background(lipgloss.Color("69")).
		Render(
			strings.Repeat(" ", s.width),
		)
}
