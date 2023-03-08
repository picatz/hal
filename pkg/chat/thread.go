package chat

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/picatz/openai"
)

type Thread struct {
	// Name (title) of the thread.
	Name string `json:"title"`

	// Summary (description) of the thread.
	Summary string `json:"description"`

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

func (ct *Thread) Title() string       { return ct.Name }
func (ct *Thread) Description() string { return ct.Summary }
func (ct *Thread) FilterValue() string { return ct.Name }

type Threads []*Thread

func (cts Threads) Len() int           { return len(cts) }
func (cts Threads) Less(i, j int) bool { return cts[i].Created.Before(cts[j].Created) }
func (cts Threads) Swap(i, j int)      { cts[i], cts[j] = cts[j], cts[i] }

func (cts Threads) ListItems() []list.Item {
	chatThreadListItems := []list.Item{}

	for _, ct := range cts {
		chatThread := ct
		chatThreadListItems = append(chatThreadListItems, chatThread)
	}

	return chatThreadListItems
}
