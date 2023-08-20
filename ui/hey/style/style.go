package style

import "github.com/charmbracelet/lipgloss"

const (
	White        = "#ffffff"
	FocusedColor = "#1DB954"
	BgColor      = "#191414"
)

func ChatNomal(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		BorderStyle(lipgloss.HiddenBorder())
}

func ChatFocused(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(FocusedColor))
}

func RecommendationNormal(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		BorderStyle(lipgloss.HiddenBorder()).
		Background(lipgloss.AdaptiveColor{Dark: BgColor, Light: BgColor})
}

func RecommendationFocused(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(FocusedColor)).
		Background(lipgloss.AdaptiveColor{Dark: BgColor, Light: BgColor})
}
