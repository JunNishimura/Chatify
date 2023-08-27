package hey

import (
	"github.com/JunNishimura/Chatify/ui/style"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	if m.err != nil {
		return style.ErrorView(m.window.Width, m.window.Height)
	}

	var s string
	if m.state == chatView {
		s += lipgloss.Place(m.window.Width, m.window.Height, lipgloss.Center, lipgloss.Center,
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				style.ChatFocused(m.getViewWidth(), m.getViewHeight()).Render(m.chatView()),
				style.RecommendationNormal(m.getViewWidth(), m.getViewHeight()).Render(m.recommendationView()),
			),
		)
	} else {
		s += lipgloss.Place(m.window.Width, m.window.Height, lipgloss.Center, lipgloss.Center,
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
	wholeWidth := m.window.Width - 4
	halfWidth := wholeWidth / 2
	return halfWidth
}

func (m *Model) getViewHeight() int {
	height := m.window.Height - 5
	return height
}

func (m *Model) chatView() string {
	var s string
	var lastSpeaker Speaker
	for _, message := range m.conversation {
		if message.speaker == Bot {
			s += style.BotChat(m.getViewWidth()).Render(message.content) + "\n\n"
			lastSpeaker = Bot
		} else if message.speaker == User {
			s += style.UserChat().Render(message.content) + "\n\n"
			lastSpeaker = User
		}
	}
	if s != "" && lastSpeaker == Bot {
		s += style.TextInput().Render(m.textInput.View())
	}
	return s
}

func (m *Model) recommendationView() string {
	return m.list.View()
}
