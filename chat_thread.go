package main

import (
	"time"

	"github.com/picatz/openai"
)

type chatThread struct {
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
}

func (ct *chatThread) Title() string       { return ct.Name }
func (ct *chatThread) Description() string { return ct.Summary }
func (ct *chatThread) FilterValue() string { return ct.Name }
