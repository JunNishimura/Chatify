package style

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

const (
	White              = "#ffffff"
	Gray               = "#777777"
	Red                = "#ff4444"
	HighlightColor     = "#1DB954"
	SemiHighlightColor = "#146542"
	BgColor            = "#191414"
)

func RecommendationNormal(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		BorderStyle(lipgloss.HiddenBorder())
}

func RecommendationFocused(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(HighlightColor))
}

func ChatNormal(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		BorderStyle(lipgloss.HiddenBorder()).
		Background(lipgloss.AdaptiveColor{Dark: BgColor, Light: BgColor})
}

func ChatFocused(width, height int) lipgloss.Style {
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

func List() (s list.DefaultItemStyles) {
	s.NormalTitle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: White, Dark: White}).
		Padding(0, 0, 0, 4)

	s.NormalDesc = s.NormalTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: Gray, Dark: Gray})

	s.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.AdaptiveColor{Light: HighlightColor, Dark: HighlightColor}).
		Foreground(lipgloss.AdaptiveColor{Light: HighlightColor, Dark: HighlightColor}).
		Padding(0, 0, 0, 1)

	s.SelectedDesc = s.SelectedTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: SemiHighlightColor, Dark: SemiHighlightColor})

	s.DimmedTitle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: HighlightColor, Dark: HighlightColor}).
		Padding(0, 0, 0, 2)

	s.DimmedDesc = s.DimmedTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: SemiHighlightColor, Dark: SemiHighlightColor})

	s.FilterMatch = lipgloss.NewStyle().Underline(true)

	return s
}

func AsciiArt() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(HighlightColor))
}

const (
	DefaultErrorMessage = "unexpected error happens.\n\nplease report at\n\nhttps://github.com/JunNishimura/Chatify/issues"
	errorDisplayWidth   = 60
	errorDisplayHeight  = 10
)

func ErrorView(message string, windowWidth, windowHeight int) string {
	return lipgloss.Place(windowWidth, windowHeight, lipgloss.Center, lipgloss.Center,
		lipgloss.NewStyle().
			Width(errorDisplayWidth).
			Height(errorDisplayHeight).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color(HighlightColor)).
			Background(lipgloss.AdaptiveColor{Dark: BgColor, Light: BgColor}).
			Render(lipgloss.Place(errorDisplayWidth, errorDisplayHeight, lipgloss.Center, lipgloss.Center,
				lipgloss.NewStyle().
					Width(errorDisplayWidth).
					Align(lipgloss.Center).
					Foreground(lipgloss.Color(Red)).
					Background(lipgloss.Color(BgColor)).
					Render(message))))
}
