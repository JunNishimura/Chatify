package hey

import (
	"strings"

	"github.com/JunNishimura/Chatify/ui/hey/style"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	var s string

	// window size adjustmen
	if m.state == chatView {
		s += lipgloss.Place(m.window.Width, m.window.Height, lipgloss.Center, lipgloss.Center,
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				style.RecommendationFocused(m.getViewWidth(), m.getViewHeight()).Render(m.chatView()),
				style.ChatNomal(m.getViewWidth(), m.getViewHeight()).Render(m.recommendationView()),
			),
		)
	} else {
		s += lipgloss.Place(m.window.Width, m.window.Height, lipgloss.Center, lipgloss.Center,
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				style.RecommendationNormal(m.getViewWidth(), m.getViewHeight()).Render(m.chatView()),
				style.ChatFocused(m.getViewWidth(), m.getViewHeight()).Render(m.recommendationView()),
			),
		)
	}
	return s
}

func (m Model) getViewWidth() int {
	wholeWidth := m.window.Width - 4
	halfWidth := wholeWidth / 2
	return halfWidth
}

func (m Model) getViewHeight() int {
	height := m.window.Height - 5
	return height
}

func (m Model) chatView() string {
	var s string
	s += strings.Join(m.conversation, "\n\n")
	s += "\n\n" + m.textInput.View()
	return s
}

func (m Model) recommendationView() string {
	return m.list.View()
}
