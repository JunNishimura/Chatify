package hey

import (
	"strings"

	"github.com/JunNishimura/Chatify/ui/hey/style"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	var s string
	if m.state == chatView {
		s += lipgloss.JoinHorizontal(
			lipgloss.Top,
			style.Focused.Render(m.chatView()),
			style.Normal.Render(m.recommendationView()),
		)
	} else {
		s += lipgloss.JoinHorizontal(
			lipgloss.Top,
			style.Normal.Render(m.chatView()),
			style.Focused.Render(m.recommendationView()),
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
	return style.List.Render(m.list.View())
}
