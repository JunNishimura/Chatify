package style

import "github.com/charmbracelet/lipgloss"

const (
	White          = "#ffffff"
	HighlightColor = "#1DB954"
	BgColor        = "#191414"
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
		BorderForeground(lipgloss.Color(HighlightColor))
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
		BorderForeground(lipgloss.Color(HighlightColor)).
		Background(lipgloss.AdaptiveColor{Dark: BgColor, Light: BgColor})
}

func BotChat(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Right).
		Width(width)
}

func UserChat() lipgloss.Style {
	return lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Left).
		Foreground(lipgloss.Color(HighlightColor))
}

func TextInput() lipgloss.Style {
	return lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Left)
}
