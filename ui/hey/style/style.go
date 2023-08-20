package style

import "github.com/charmbracelet/lipgloss"

func GetNormal(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		BorderStyle(lipgloss.HiddenBorder())
}

func GetFocused(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("69"))
}
