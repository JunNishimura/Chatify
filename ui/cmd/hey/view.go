package hey

import (
	"errors"
	"fmt"

	"github.com/JunNishimura/Chatify/ui/cmd/base"
	"github.com/JunNishimura/Chatify/ui/style"
	"github.com/charmbracelet/lipgloss"
)

var (
	noDeviceMessage = "Fail to find active devices.\n\nPlease try again after opening Spotify App."
)

func (m *Model) View() string {
	if m.err != nil {
		if errors.Is(m.err, errNoDeviceFound) {
			return style.ErrorView(noDeviceMessage, m.Window.Width, m.Window.Height)
		}
		return style.ErrorView(style.DefaultErrorMessage, m.Window.Width, m.Window.Height)
	}

	if m.isQuit {
		var msg string
		if m.opts.playlist {
			msg = fmt.Sprintf("I made a playlist for you!!\n\n%s\n\nLet's talk again!!", m.playlist.ExternalURLs["spotify"])
		} else {
			msg = "Thanks for talking to me!\n\nLet's talk again!!"
		}
		return style.QuitView(msg, m.Window.Width, m.Window.Height)
	}

	if m.questionDone {
		return lipgloss.Place(m.Window.Width, m.Window.Height, lipgloss.Center, lipgloss.Center,
			style.RecommendationFocused(m.getViewWidth(), m.getViewHeight()).Render(m.recommendationView()),
		)
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
