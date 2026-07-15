package main

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39")).Padding(0, 1)
	headerStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")).Underline(true)
	cursorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true)
	descStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Italic(true)
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	errorStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196"))
	successStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
)
