package main

// A simple program demonstrating the text area component from the Bubbles
// component library.

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	"github.com/picatz/openai"
)

func main() {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
		tea.WithMouseAllMotion(),
		// tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type editorFinishedMsg struct {
	err    error
	buffer []byte
}

func openEditor(buffer string) tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	// Write to a temp file and open it
	f, err := os.CreateTemp("", "hal-editor-*")
	if err != nil {
		panic(err)
	}

	f.WriteString(buffer)
	f.Sync()

	c := exec.Command(editor, f.Name())
	return tea.ExecProcess(c, func(err error) tea.Msg {
		defer func() {
			f.Close()
			os.Remove(f.Name())
		}()
		if err != nil {
			return editorFinishedMsg{err: err}
		}

		// Read the file back in
		f, err := os.Open(f.Name())
		if err != nil {
			return editorFinishedMsg{err: err}
		}

		buf := make([]byte, 4096)
		n, err := f.Read(buf)
		if err != nil {
			return editorFinishedMsg{err: err}
		}

		return editorFinishedMsg{
			err:    nil,
			buffer: buf[:n],
		}
	})
}

type chatFinishedMsg struct {
	err     error
	buffer  []byte
	history []openai.ChatMessage
	tokens  int
}

func sendChatRequest(client *openai.Client, chatHistory []openai.ChatMessage, text string) tea.Cmd {
	return func() tea.Msg {
		// send the message to the OpenAI chat API
		ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
		defer cancel()

		chatHistory = append(chatHistory, openai.ChatMessage{
			Role:    openai.ChatRoleUser,
			Content: text,
		})

		resp, err := client.CreateChat(ctx, &openai.CreateChatRequest{
			Model:    openai.ModelGPT35Turbo,
			Messages: chatHistory,
		})
		if err != nil {
			return editorFinishedMsg{err: err}
		}

		// Add response to chat history
		chatHistory = append(chatHistory, openai.ChatMessage{
			Role:    openai.ChatRoleSystem,
			Content: resp.Choices[0].Message.Content,
		})

		return chatFinishedMsg{
			err:     nil,
			buffer:  []byte(resp.Choices[0].Message.Content),
			history: chatHistory,
			tokens:  resp.Usage.TotalTokens,
		}
	}
}

type (
	errMsg error
)

type model struct {
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style

	spinner spinner.Model
	running bool
	err     error

	client      *openai.Client
	chatTokens  int
	chatHistory []openai.ChatMessage
}

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "What do you want to do?"

	ta.SetWidth(80)

	ta.Prompt = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("69")).Render("│ ")
	ta.CharLimit = 4096
	ta.Focus()
	ta.Focused()

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(-1, 1)
	vp.SetContent(
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			"Welcome to",
			lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true).Render(" HAL"),
			"!",
		),
	)

	ta.KeyMap.InsertNewline.SetEnabled(true)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY environment variable is not set")
		os.Exit(1)
	}

	return model{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("69")),
		err:         nil,

		spinner: s,
		running: false,

		client: openai.NewClient(apiKey),
		chatHistory: []openai.ChatMessage{
			// Configure the chat session. This should be included in every chat request.
			{
				Role:    openai.ChatRoleSystem,
				Content: "You are HAL, a code and text editor controlled by natural language. Answer as concisely as possible.",
			},
		},
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC: // Quit the program.
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyCtrlE: // Open editor with current text with textarea buffer.
			return m, openEditor(m.textarea.Value())
		case tea.KeyEscape:
			text := m.textarea.Value()

			m.running = true
			m.viewport.GotoBottom()
			m.textarea.Reset()

			// send the message to the OpenAI chat API
			return m, sendChatRequest(m.client, m.chatHistory, text)
		}
	case editorFinishedMsg: // When the editor is finished, update the textarea with buffer.
		if msg.err != nil {
			panic(msg.err)
			// return m, tea.Quit
		}
		m.textarea.SetValue(string(msg.buffer))
	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	case chatFinishedMsg:
		m.running = false

		if msg.err != nil {
			return m, tea.Quit
		}

		m.chatHistory = msg.history
		m.chatTokens = msg.tokens

		r, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(m.viewport.Width),
		)
		if err != nil {
			panic(err)
		}

		msg.buffer, err = r.RenderBytes(msg.buffer)
		if err != nil {
			panic(err)
		}

		msg.buffer = bytes.Trim(msg.buffer, "\r\n")

		m.viewport.SetContent(
			lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true).Render(fmt.Sprintf("HAL:%d:%d", len(m.chatHistory), m.chatTokens)),

				strings.TrimSpace(string(msg.buffer)),
			),
		)

		m.viewport.Height = strings.Count(string(msg.buffer), "\n") + 2
		m.viewport.GotoTop()

		// m.viewport.SetContent(strings.TrimSpace(string(msg.buffer)))
		// Window size is received when starting up and on every resize
	case tea.WindowSizeMsg:
		m.viewport.Style.Width(msg.Width)
		m.viewport.Style.Height(msg.Height - 2)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, tea.Batch(tiCmd, vpCmd, m.spinner.Tick)
}

func (m model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Top,
		func() string {
			if m.running {
				return lipgloss.JoinHorizontal(
					lipgloss.Left,
					m.spinner.View(),
					"Running...",
				)
			}

			return m.viewport.View()
		}(),
		"",
		m.textarea.View(),
		"",
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			lipgloss.NewStyle().Faint(true).Render("Press"),
			lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true).Render(" ESC "),
			lipgloss.NewStyle().Faint(true).Render("to run,"),
			lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true).Render(" CTRL+E "),
			lipgloss.NewStyle().Faint(true).Render("to edit,"),
			lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true).Render(" CTRL+C "),
			lipgloss.NewStyle().Faint(true).Render("to quit."),
		),
	)
}
