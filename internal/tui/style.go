package tui

import "github.com/charmbracelet/lipgloss"

var Accent = lipgloss.Color("#9fef00")

var Border = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	Padding(0, 1)

var Panel = lipgloss.NewStyle().
	Padding(0, 1)

var ActiveTab = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(Accent).
	Bold(true).
	Padding(0, 1)

var InactiveTab = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("240")).
	Bold(false).
	Padding(0, 1)

var Title = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("39"))

var Warning = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#ff5f5f"))

var Footer = lipgloss.NewStyle().
	Foreground(lipgloss.Color("240"))

var footerStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("240")).
	Padding(1, 1).
	BorderTop(true).
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("238")).
	Width(0)
