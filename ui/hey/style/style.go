package style

import "github.com/charmbracelet/lipgloss"

var (
	Normal = lipgloss.NewStyle().
		Width(100).
		Height(40).
		BorderStyle(lipgloss.HiddenBorder())
	Focused = lipgloss.NewStyle().
		Width(100).
		Height(40).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("69"))
	List = lipgloss.NewStyle().
		Margin(1, 2)
)
