package main

// Mode of the editor.
type Mode int

const (
	ModeChat Mode = iota
	ModeInsert
	ModeShell
)
