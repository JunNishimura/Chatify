package hey

import (
	"github.com/JunNishimura/Chatify/ui/cmd/base"
	"github.com/JunNishimura/Chatify/ui/style"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	if m.err != nil {
		return style.ErrorView(m.Window.Width, m.Window.Height)
	}

	var s string
	if m.state == chatView {
		s += lipgloss.Place(m.Window.Width, m.Window.Height, lipgloss.Center, lipgloss.Center,
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				style.ChatFocused(m.getViewWidth(), m.getViewHeight()).Render(m.chatView()),
				style.RecommendationNormal(m.getViewWidth(), m.getViewHeight()).Render(m.recommendationView()),
			),
		)
	} else {
		s += lipgloss.Place(m.Window.Width, m.Window.Height, lipgloss.Center, lipgloss.Center,
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				style.ChatNormal(m.getViewWidth(), m.getViewHeight()).Render(m.chatView()),
				style.RecommendationFocused(m.getViewWidth(), m.getViewHeight()).Render(m.recommendationView()),
			),
		)
	}
	return s
}

func (m *Model) getViewWidth() int {
	wholeWidth := m.Window.Width - 4
	halfWidth := wholeWidth / 2
	return halfWidth
}

func (m *Model) getViewHeight() int {
	height := m.Window.Height - 5
	return height
}

func (m *Model) chatView() string {
	var s string
	var lastSpeaker base.Speaker
	for _, message := range m.Conversation {
		if message.Speaker == base.Bot {
			s += style.BotChat(m.getViewWidth()).Render(message.Content) + "\n\n"
			lastSpeaker = base.Bot
		} else if message.Speaker == base.User {
			s += style.UserChat().Render(message.Content) + "\n\n"
			lastSpeaker = base.User
		}
	}
	if s != "" && lastSpeaker == base.Bot {
		s += style.TextInput().Render(m.TextInput.View())
	}
	return s
}

func (m *Model) recommendationView() string {
	return m.list.View()
}
