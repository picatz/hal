package statusbar

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/ansi"
	"github.com/picatz/hal/pkg/chat"
)

var (
	messageCountStatusBarBlockStyle = lipgloss.NewStyle().Background(lipgloss.Color("63"))

	tokensCountStatusBarBlockStyle = lipgloss.NewStyle().Background(lipgloss.Color("62"))

	currentThreadNameBlockStyle = lipgloss.NewStyle().Background(lipgloss.Color("69")).Bold(true)
)

// ChatThreadMsg is a message sent to the status bar.
type ChatThreadMsg struct {
	ChatThread *chat.Thread
}

// SpinMsg is a message sent to the status bar.
type SpinMsg struct {
	Spinning bool
}

// Model is a status bar model.
type Model struct {
	Width int
	Style lipgloss.Style

	Spinner  spinner.Model
	Spinning bool

	ChatThread *chat.Thread
}

// New creates a new status bar component.
func New() *Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	// s.Style = lipgloss.NewStyle().
	// 	Background(lipgloss.Color("69")).
	// 	Bold(true)
	return &Model{
		Spinner: s,
		Style: lipgloss.NewStyle().
			// Padding(0, 1).
			Background(lipgloss.Color("69")),
	}
}

// Init implements tea.Model, but does nothing currently.
func (s *Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model, handles window size and status bar messages.
func (s *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.Width = msg.Width
	case *ChatThreadMsg:
		if msg.ChatThread != nil {
			s.ChatThread = msg.ChatThread
		}
	}

	if s.Spinning {
		return s, s.Spinner.Tick
	}

	return s, nil
}

// View implements tea.Model, which actually renders the status bar.
//
// This status bar is still early, but hopefully it'll be a good start.
//
// It should be much more customizable, but for now it's just a simple
// status bar with the current thread name, message count, and token count.
func (s *Model) View() string {
	var (
		currentThreadName string = "*"
		chatMessageCount  int
		chatTokensCount   int
	)

	// v := truncate.StringWithTail(s.msg.Value, uint(maxWidth-st.StatusBarValue.GetHorizontalFrameSize()), "…")

	if s.ChatThread != nil {
		currentThreadName = s.ChatThread.Name
		chatMessageCount = len(s.ChatThread.ChatHistory)
		chatTokensCount = s.ChatThread.Tokens
	}

	var spinnerBlock string
	if s.Spinning {
		spinnerBlock = s.Spinner.View()
	} else {
		spinnerBlock = "»"
	}

	var (
		messageCountStatusBarBlock = messageCountStatusBarBlockStyle.Render(fmt.Sprintf(" Messages: %d ", chatMessageCount))
		tokensCountStatusBarBlock  = tokensCountStatusBarBlockStyle.Render(fmt.Sprintf(" Tokens: %d ", chatTokensCount))
		currentThreadNameBlock     = currentThreadNameBlockStyle.Render(fmt.Sprintf(" %s ", currentThreadName))
	)

	// get printable characters (non ANSI escape codes)
	printableChars := ansi.PrintableRuneWidth(messageCountStatusBarBlock + tokensCountStatusBarBlock + currentThreadNameBlock)

	// build status bar including the current thread name on right hand side, filling the rest of the space with spaces
	statusText := " " + spinnerBlock + strings.Repeat(" ", s.Width-printableChars-6) + messageCountStatusBarBlock + tokensCountStatusBarBlock + currentThreadNameBlock + " "

	// TODO: add a way to set the status bar style, and stuff inside it.
	return s.Style.Render(statusText)
}
