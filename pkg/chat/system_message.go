package chat

import "github.com/picatz/openai"

var SystemMessage = openai.ChatMessage{
	Role:    openai.ChatRoleSystem,
	Content: "You are HAL, a powerful code and text editor controlled by natural language. Answer as concisely as possible.",
}
