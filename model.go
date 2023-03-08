package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/picatz/openai"

	"github.com/picatz/hal/pkg/chat"
	"github.com/picatz/hal/pkg/editor"
	"github.com/picatz/hal/pkg/statusbar"
)

var welcomeToHAL = lipgloss.JoinHorizontal(
	lipgloss.Left,
	"Welcome to",
	lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true).Render(" HAL"),
	"!",
)

type model struct {
	// Mode is the current mode of the application. TODO.
	mode Mode

	// General dimensions.
	width  int
	height int

	// Shared base style for HAL.
	halStyle lipgloss.Style

	// Error message to display, usually from OpenAI.
	err error

	editor textarea.Model

	// Status bar.
	statusbar *statusbar.Model

	// OpenAI API client and chat history.
	client            *openai.Client
	chatSystemMessage openai.ChatMessage
	chatThreadList    list.Model
	chatThreads       chat.Threads
	currnetThread     *chat.Thread
}

func newModel() model {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY environment variable is not set")
		os.Exit(1)
	}

	// Statusbar is a shown at the bottom of the screen, and is used to display
	// information about the current thread, message count, and token count.
	//
	// In the future it will be used to display other information, such as
	// the current mode, and other things. It's a work in progress.
	statusbar := statusbar.New()

	// TODO: Keep track of threads over time, allow switching between them.
	//       There's so much to do here, but this is a start.
	chatThreads := chat.Threads{
		{
			Name:    "Get to know HAL",
			Summary: "Learn how to work together.",
			Created: time.Now(),
			ChatHistory: []openai.ChatMessage{
				chat.SystemMessage,
			},
		},
	}

	// Setup chat thread list.
	chatThreadList := list.New(chatThreads.ListItems(), list.NewDefaultDelegate(), 80, 10)
	chatThreadList.SetWidth(80)
	chatThreadList.SetHeight(15)
	chatThreadList.Styles.FilterCursor = halStyleColor

	chatThreadList.Title = "Threads"
	chatThreadList.Styles.TitleBar = halStyleColor
	chatThreadList.Styles.FilterCursor = halStyleColor
	chatThreadList.Styles.FilterPrompt = halStyleColor
	chatThreadList.Styles.DefaultFilterCharacterMatch = halStyleColor
	chatThreadList.SetShowHelp(false) // true?

	// Setup text area for user input.
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

	return model{
		// Started in chat thread list mode by default (if not file selected in args?)
		mode: ModeChatThreadList,

		editor: editor,

		halStyle: halStyleColor,
		err:      nil,

		chatThreads:    chatThreads,
		chatThreadList: chatThreadList,

		client:            openai.NewClient(apiKey),
		chatSystemMessage: chat.SystemMessage,

		statusbar: statusbar,
	}
}

// Init implements tea.Model, it just starts the blinking cursor.
func (m model) Init() tea.Cmd {
	return textarea.Blink
}

// Update implements tea.Model, it handles all user input and updates the
// model accordingly.
//
// It might switch the application into a different mode, or update state
// in the current mode.
//
// Updates for any mode are batched and applied at the end of this function.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		statusbarCmd      tea.Cmd
		textareaCmd       tea.Cmd
		chatThreadListCmd tea.Cmd
	)

	// Handle status bar updates, always show status bar.
	m.statusbar, statusbarCmd = m.statusbar.Update(msg)

	// Handle update based on current mode.
	switch m.mode {
	case ModeChatThreadList:
		m.chatThreadList, chatThreadListCmd = m.chatThreadList.Update(msg)
	case ModeEditorInsert:
		m.editor, textareaCmd = m.editor.Update(msg)
	default:
		// TODO: handle other modes.
	}

	// TODO: have mode-specific keybindings.
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC: // Quit the program.
			return m, tea.Quit
		case tea.KeyCtrlE: // Open editor with current text with textarea buffer.
			return m, editor.OpenExternal(m.editor.Value())
		// case tea.KeyCtrlO: // Open editor with current text with viewport buffer.
		// 	// Strip any ANSI color codes from the viewport.
		// 	vpView := Strip(m.chatOutput.View())
		// 	return m, openEditor(vpView, true)
		case tea.KeyCtrlT: // Truncate the previous chat history.
			if m.currnetThread != nil && len(m.currnetThread.ChatHistory) > 2 {
				lastMessage := m.currnetThread.ChatHistory[len(m.currnetThread.ChatHistory)-1]

				m.currnetThread.ChatHistory = []openai.ChatMessage{
					m.chatSystemMessage,
					lastMessage,
				}
			}
		// case tea.KeyCtrlL: // Clear the viewport.
		// 	m.chatOutput.SetContent("")
		case tea.KeyEscape:
			text := m.editor.Value()

			// m.chatOutput.GotoBottom()
			m.editor.Reset()
			m.editor.Placeholder = "..."

			m.statusbar.Spinning = true

			// send the message to the OpenAI chat API

			sendCmd := chat.Send(m.client, m.currnetThread.ChatHistory, text)

			return m, tea.Batch(sendCmd, statusbarCmd, textareaCmd, chatThreadListCmd, m.statusbar.Spinner.Tick)
		case tea.KeyEnter:
			if m.currnetThread == nil {
				// Select the thread.
				m.editor.SetValue("") // For some reason, the text area is not cleared when selecting a thread.
				m.currnetThread = m.chatThreadList.SelectedItem().(*chat.Thread)

				if len(m.currnetThread.ChatHistory) == 0 {
					m.currnetThread.ChatHistory = append(m.currnetThread.ChatHistory, m.chatSystemMessage)
				}

				// Update the status bar with the current thread.
				m.statusbar.Update(&statusbar.Model{
					ChatThread: m.currnetThread,
				})

				// Change the mode to editor mode.
				m.mode = ModeEditorInsert

				return m, nil
			}
		}
	case editor.ExternalFinishedMsg: // When the editor is finished, update the textarea with buffer.
		if msg.Err != nil {
			panic(msg.Err)
			// return m, tea.Quit
		}

		m.editor.SetValue(string(msg.Buffer))
	case chat.FinishedMsg:
		if msg.Err != nil {
			return m, tea.Quit
		}

		m.currnetThread.ChatHistory = msg.History
		m.currnetThread.Tokens = msg.Tokens

		m.statusbar.ChatThread = m.currnetThread

		m.statusbar.Spinning = false

		m.editor.SetValue(string(msg.Buffer))
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.editor.SetHeight(msg.Height - 4)
		m.editor.SetWidth(msg.Width)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.statusbar.Spinner, cmd = m.statusbar.Spinner.Update(msg)
		return m, cmd
	}

	return m, tea.Batch(statusbarCmd, textareaCmd, chatThreadListCmd)
}

func (m model) chooseThreadListView() string {
	return lipgloss.JoinVertical(
		lipgloss.Top,
		welcomeToHAL,
		"",
		m.chatThreadList.View(),
	)
}

func (m model) viewChatInput() string {
	return m.editor.View()
}

func (m model) View() string {
	var mainView string

	if m.currnetThread == nil {
		mainView = m.chooseThreadListView()
	} else {
		mainView = lipgloss.JoinVertical(
			lipgloss.Top,
			// m.viewChatOutput(),
			m.viewChatInput(),
		)
	}

	// The space between the main view and the statusbar (sticky footer).
	//
	// -2 for the statusbar, -1 for the newline.
	spaceBetween := m.height - strings.Count(mainView, "\n") - 3

	// If there's no space for the statusbar, just return the main view.
	if spaceBetween < 0 {
		return mainView
	}

	// Place the statusbar view at the bottom of the screen, fit the rest above it.
	return lipgloss.JoinVertical(
		lipgloss.Top,
		mainView,
		strings.Repeat("\n", spaceBetween),
		m.statusbar.View(),
	)
}
