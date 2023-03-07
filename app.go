package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/picatz/openai"
)

var welcomeToHAL = lipgloss.JoinHorizontal(
	lipgloss.Left,
	"Welcome to",
	lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true).Render(" HAL"),
	"!",
)

type model struct {
	// General dimensions.
	width  int
	height int

	// Shared base style for HAL.
	halStyle lipgloss.Style

	// Error message to display, usually from OpenAI.
	err error

	// Spinner is the loading indicator for OpenAI requests,
	// or any other long running process impacting the viewport.
	chatSpinner spinner.Model
	chatLoading bool

	// chatOutput is the current chat request (input).
	chatInput textarea.Model

	// chatOutput is the current chat response (output).
	chatOutput viewport.Model

	// Status bar.
	statusbar *statusBar

	// OpenAI API client and chat history.
	client            *openai.Client
	chatTokens        int
	chatSystemMessage openai.ChatMessage
	chatThreadList    list.Model
	chatThreads       []chatThread
	currnetThread     *chatThread
}

func newModel() model {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY environment variable is not set")
		os.Exit(1)
	}

	statusbar := newStatusBar()

	systemMessage := openai.ChatMessage{
		Role:    openai.ChatRoleSystem,
		Content: "You are HAL, a powerful code and text editor controlled by natural language. Answer as concisely as possible.",
	}

	chatThreadListItems := []list.Item{}

	// Add starting chat threads.
	chatThreads := []chatThread{
		{
			Name:    "Get to know HAL",
			Summary: "Learn how to work together.",
			Created: time.Now(),
			ChatHistory: []openai.ChatMessage{
				systemMessage,
			},
		},
		// TODO: keep track of threads over time, allow switching between them.
	}

	for _, ct := range chatThreads {
		chatThread := ct
		chatThreadListItems = append(chatThreadListItems, &chatThread)
	}

	// Setup chat thread list.
	chatThreadList := list.New(chatThreadListItems, list.NewDefaultDelegate(), 80, 10)
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
	ta := textarea.New()
	ta.Placeholder = "What do you want to do?"
	ta.Prompt = halStyleColor.Bold(true).Render("│ ")
	ta.CharLimit = 4096
	ta.SetWidth(80)
	ta.Focus()
	ta.Focused()

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false

	vp := viewport.New(-1, 1)
	// vp.SetContent(welcomeToHAL)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = halStyleColor

	return model{
		chatInput:  ta,
		chatOutput: vp,
		halStyle:   halStyleColor,
		err:        nil,

		chatThreads:    chatThreads,
		chatThreadList: chatThreadList,

		chatSpinner: s,
		chatLoading: false,

		client:            openai.NewClient(apiKey),
		chatSystemMessage: systemMessage,

		statusbar: statusbar,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	statusbarModel, statusbarCmd := m.statusbar.Update(msg)

	m.statusbar = statusbarModel.(*statusBar)

	var (
		textareaCmd      tea.Cmd
		viewportCmd      tea.Cmd
		hatThreadListCmd tea.Cmd
	)

	m.chatInput, textareaCmd = m.chatInput.Update(msg)
	m.chatOutput, viewportCmd = m.chatOutput.Update(msg)

	// We haven't selected a threas yet, so allow for selection of the thread.
	if m.currnetThread == nil {
		m.chatThreadList, hatThreadListCmd = m.chatThreadList.Update(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC: // Quit the program.
			return m, tea.Quit
		case tea.KeyCtrlE: // Open editor with current text with textarea buffer.
			return m, openEditor(m.chatInput.Value(), false)
		case tea.KeyCtrlO: // Open editor with current text with viewport buffer.
			// Strip any ANSI color codes from the viewport.
			vpView := Strip(m.chatOutput.View())

			return m, openEditor(vpView, true)
		case tea.KeyCtrlT: // Truncate the previous chat history.
			if m.currnetThread != nil && len(m.currnetThread.ChatHistory) > 2 {
				lastMessage := m.currnetThread.ChatHistory[len(m.currnetThread.ChatHistory)-1]

				m.currnetThread.ChatHistory = []openai.ChatMessage{
					m.chatSystemMessage,
					lastMessage,
				}
			}
		case tea.KeyCtrlL: // Clear the viewport.
			m.chatOutput.SetContent("")
		case tea.KeyEscape:
			text := m.chatInput.Value()

			m.chatLoading = true
			m.chatOutput.GotoBottom()
			m.chatInput.Reset()

			// send the message to the OpenAI chat API
			return m, sendChatRequest(m.client, m.currnetThread.ChatHistory, text)
		case tea.KeyEnter:
			if m.currnetThread == nil {
				// Select the thread.
				m.chatInput.SetValue("") // For some reason, the text area is not cleared when selecting a thread.
				m.currnetThread = m.chatThreadList.SelectedItem().(*chatThread)

				if len(m.currnetThread.ChatHistory) == 0 {
					m.currnetThread.ChatHistory = append(m.currnetThread.ChatHistory, m.chatSystemMessage)
				}

				m.chatOutput.SetContent("» " + m.currnetThread.Name)
				return m, nil
			}
		}
	case editorFinishedMsg: // When the editor is finished, update the textarea with buffer.
		if msg.err != nil {
			panic(msg.err)
			// return m, tea.Quit
		}

		if msg.viewport {
			// Update the viewport with the editor buffer.
			styledMsgBuffer, _, err := renderMarkdown(msg.buffer, 80)
			if err != nil {
				panic(err)
			}

			msg.buffer = styledMsgBuffer

			m.chatOutput.SetContent(string(msg.buffer))

			// Update the last message in the chat history with the editor buffer.
			lastMessage := m.currnetThread.ChatHistory[len(m.currnetThread.ChatHistory)-1]
			lastMessage.Content = string(msg.buffer)
			m.currnetThread.ChatHistory[len(m.currnetThread.ChatHistory)-1] = lastMessage
		} else {
			m.chatInput.SetValue(string(msg.buffer))
		}
	// We handle errors just like any other message
	case errMsg:
		m.err = fmt.Errorf("error: %s", msg)
		return m, nil
	case chatFinishedMsg:
		m.chatLoading = false

		if msg.err != nil {
			return m, tea.Quit
		}

		m.currnetThread.ChatHistory = msg.history
		m.chatTokens = msg.tokens

		styledMsgBuffer, height, err := renderMarkdown(msg.buffer, 80)
		if err != nil {
			panic(err)
		}

		msg.buffer = styledMsgBuffer

		m.chatOutput.SetContent(
			lipgloss.JoinVertical(
				lipgloss.Left,
				string(msg.buffer),
			),
		)

		m.chatOutput.Height = height + 2
		m.chatOutput.GotoTop()
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.chatOutput.Style.Width(msg.Width)
		m.chatOutput.Style.Height(msg.Height - 2)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.chatSpinner, cmd = m.chatSpinner.Update(msg)
		return m, cmd
	}

	return m, tea.Batch(statusbarCmd, textareaCmd, viewportCmd, hatThreadListCmd, m.chatSpinner.Tick)
}

func (m model) chooseThreadListView() string {
	return lipgloss.JoinVertical(
		lipgloss.Top,
		welcomeToHAL,
		"",
		m.chatThreadList.View(),
		"",
	)
}

func (m model) viewChatInput() string {
	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.chatInput.View(),
		"",
		m.viewChatInputDiagnostics(),
		"",
		m.viewChatInputHelp(),
	)
}

func (m model) viewChatInputDiagnostics() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true).Render("Bytes: "),
		lipgloss.NewStyle().Render(fmt.Sprint(len(m.chatInput.Value()))),
		" ",
		lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true).Render("~Tokens: "),
		lipgloss.NewStyle().Render(fmt.Sprint(countTokens(m.chatInput.Value()))),
		" ",
		lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true).Render("Chat Tokens: "),
		lipgloss.NewStyle().Render(fmt.Sprint(m.chatTokens)),
		" ",
		lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true).Render("Chat Messages: "),
		lipgloss.NewStyle().Render(fmt.Sprint(len(m.currnetThread.ChatHistory))),
	)
}

func (m model) viewChatOutput() string {
	if m.chatLoading {
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			m.chatSpinner.View(),
			"Running...",
		)
	}

	return m.chatOutput.View()
}

var (
	halHelpFaint = lipgloss.NewStyle().Faint(true)
	halHelpBold  = lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true)

	// TODO: make configurable and deeply consider the result.
	viewChatInputHelpCached = lipgloss.JoinHorizontal(
		lipgloss.Left,
		halHelpFaint.Render("Press"),
		halHelpBold.Render(" ESC "), halHelpFaint.Render("to run,"),
		halHelpBold.Render(" CTRL+E "), halHelpFaint.Render("to edit input,"),
		halHelpBold.Render(" CTRL+O "), halHelpFaint.Render("to edit output,"),
		halHelpBold.Render(" CTRL+C "), halHelpFaint.Render("to quit."),
	)
)

func (m model) viewChatInputHelp() string {
	return viewChatInputHelpCached
}

func (m model) View() string {
	if m.currnetThread == nil {
		return m.chooseThreadListView()
	}

	mainView := lipgloss.JoinVertical(
		lipgloss.Top,
		m.viewChatOutput(),
		"",
		m.viewChatInput(),
	)

	// The space between the main view and the statusbar (sticky footer).
	//
	// -2 for the statusbar, -1 for the newline.
	spaceBetween := m.height - strings.Count(mainView, "\n") - 2 - 1

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
