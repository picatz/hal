package main

import "github.com/charmbracelet/lipgloss"

var welcomeToHAL = lipgloss.JoinHorizontal(
	lipgloss.Left,
	"Welcome to",
	lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true).Render(" HAL"),
	"!",
)
