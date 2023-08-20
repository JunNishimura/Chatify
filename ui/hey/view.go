package hey

import (
	"strings"

	"github.com/JunNishimura/Chatify/ui/hey/style"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	var s string

	// window size adjustmen
	wholeWidth := m.window.Width - 4
	halfWidth := wholeWidth / 2
	height := m.window.Height - 5
	if m.state == chatView {
		s += lipgloss.Place(m.window.Width, m.window.Height, lipgloss.Center, lipgloss.Center,
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				style.GetFocused(halfWidth, height).Render(m.chatView()),
				style.GetNormal(halfWidth, height).Render(m.recommendationView()),
			),
		)
	} else {
		s += lipgloss.Place(m.window.Width, m.window.Height, lipgloss.Center, lipgloss.Center,
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				style.GetNormal(halfWidth, height).Render(m.chatView()),
				style.GetFocused(halfWidth, height).Render(m.recommendationView()),
			),
		)
	}
	return s
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
