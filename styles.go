package main

import "github.com/charmbracelet/lipgloss"

var titleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#7D56F4"))

var navStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#9D56F4"))

var errStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#eb5449"))
