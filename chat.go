package main

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/picatz/openai"
)

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
