package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/picatz/hal/pkg/chat"
)

func ChatThreadList(chatThreads chat.Threads) list.Model {
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

	return chatThreadList
}
