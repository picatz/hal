package chat

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/picatz/openai"

	"golang.org/x/text/language"
	"golang.org/x/text/search"
)

// Thread is a "chat thread" that is used to store the chat history and
// metadata for a chat session. It implements the list.Item interface
// so that it can shown in a list in the UI.
type Thread struct {
	// Name (title) of the thread.
	Name string `json:"name"`

	// Summary (description) of the thread.
	Summary string `json:"summary"`

	// Created is the date the thread was created.
	Created time.Time `json:"date"`

	// The prompt is the initial text that the model will generate from.
	// This is the text that the user will see when they first open the
	// chat window.
	Prompt string `json:"prompt"`

	// The chat history is the list of messages that have been sent and
	// received in the chat session.
	ChatHistory []openai.ChatMessage `json:"chat_history"`

	// Tokens is the last reported number of tokens used in the chat session.
	Tokens int `json:"tokens"`
}

// Implement the list.Item interface.
func (ct *Thread) Title() string       { return ct.Name }
func (ct *Thread) Description() string { return ct.Summary }
func (ct *Thread) FilterValue() string { return ct.Name }

// Threads is a collection of chat threads.
type Threads []*Thread

// Implement the sort.Interface interface.
func (cts Threads) Len() int           { return len(cts) }
func (cts Threads) Less(i, j int) bool { return cts[i].Created.Before(cts[j].Created) }
func (cts Threads) Swap(i, j int)      { cts[i], cts[j] = cts[j], cts[i] }

// ListItems returns a list of list.Item's from the Threads.
func (cts Threads) ListItems() []list.Item {
	chatThreadListItems := []list.Item{}

	for _, ct := range cts {
		chatThread := ct
		chatThreadListItems = append(chatThreadListItems, chatThread)
	}

	return chatThreadListItems
}

// Summarize returns a summrized version of the chat history.
func (ct *Thread) Summarize(ctx context.Context, client *openai.Client) (string, error) {
	// Create a new thread with a new system prompt to summarize conversation.
	chatHistory := []openai.ChatMessage{
		{
			Role:    openai.ChatRoleSystem,
			Content: "Answer as concisely as possible to summarize a conversation, capturing the most important points to continue the conversation.",
		},
		{
			Role: openai.ChatRoleUser,
			Content: func() string {
				var b strings.Builder

				b.WriteString("Please summarize the following conversation:\n\n")

				for _, m := range ct.ChatHistory {
					if m.Role == openai.ChatRoleSystem {
						continue
					}
					b.WriteString(fmt.Sprintf("%s: %s", m.Role, m.Content))
					b.WriteString("\n")
				}

				return b.String()
			}(),
		},
	}

	// create a summary of the chat history
	summary, err := client.CreateChat(ctx, &openai.CreateChatRequest{
		Model:    openai.ModelGPT35Turbo,
		Messages: chatHistory,
	})

	if err != nil {
		return "", fmt.Errorf("failed to create summary of chat thread %q: %w", ct.Name, err)
	}

	return summary.Choices[0].Message.Content, nil
}

// SearchResult is a search result for a chat thread.
type SearchResult struct {
	// The message that matched the search query.
	Message *openai.ChatMessage

	// MessageIndex is the index of the message in the chat history.
	MessageIndex int

	// MatchStart is the index of the start of the match in the message.
	StartIndex int

	// MatchEnd is the index of the end of the match in the message.
	EndIndex int
}

// Search retruns the messages that match the search query. Searching
// happens locally on the host.
//
// TODO: consider returning index information in the search results so
// that we can highlight the matches in the UI.
func (ct *Thread) Search(ctx context.Context, query string) ([]*SearchResult, error) {
	matcher := search.New(language.AmericanEnglish, search.IgnoreCase)

	pattern := matcher.CompileString(query)

	results := []*SearchResult{}

	for i, m := range ct.ChatHistory {
		msg := m
		if start, end := pattern.IndexString(msg.Content); start != -1 && end != -1 {
			results = append(results, &SearchResult{
				Message:      &msg,
				MessageIndex: i,
				StartIndex:   start,
				EndIndex:     end,
			})
		}
	}

	return results, nil
}
