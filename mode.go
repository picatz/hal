package main

// Mode of the editor (still figuring it out).
type Mode int

const (
	ModeChatThreadList Mode = iota
	ModeEditorInsert
	ModeShell
)
