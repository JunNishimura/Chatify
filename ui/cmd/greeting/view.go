package greeting

import (
	"fmt"

	"github.com/JunNishimura/Chatify/ui/cmd/base"
	"github.com/JunNishimura/Chatify/ui/style"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	if m.err != nil {
		return style.ErrorView(style.DefaultErrorMessage, m.Window.Width, m.Window.Height)
	}
	return lipgloss.Place(m.Window.Width, m.Window.Height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			style.AsciiArt().Render(asciiArtView),
			style.ChatFocused(m.getViewWidth(), m.getViewHeight()).Render(m.chatView()),
		),
	)
}

func (m *Model) getViewWidth() int {
	halfWidth := m.Window.Width / 2
	return halfWidth
}

func (m *Model) getViewHeight() int {
	halfHeight := m.Window.Height * 2 / 3
	return halfHeight
}

var asciiArtView = `
  ___  _   _    __   ____  ____  ____  _  _
 / __)( )_( )  /__\ (_  _)(_  _)( ___)( \/ )
( (__  ) _ (  /(__)\  )(   _)(_  )__)  \  /
 \___)(_) (_)(__)(__)(__) (____)(__)   (__)
 `

func (m *Model) chatView() string {
	var s string
	switch m.phase {
	case completePhase:
		var ss string
		ss += lipgloss.NewStyle().Width(m.getViewWidth()).Align(lipgloss.Center, lipgloss.Center).Render(fmt.Sprintf("Nice to see you, %s!", m.user.DisplayName)) + "\n\n"
		ss += lipgloss.NewStyle().Width(m.getViewWidth()).Align(lipgloss.Center, lipgloss.Center).Render("If you want to talk to me, please type") + "\n\n"
		ss += lipgloss.NewStyle().Width(m.getViewWidth()).Align(lipgloss.Center, lipgloss.Center).Foreground(lipgloss.Color(style.HighlightColor)).Background(lipgloss.Color(style.BgColor)).Render("$ chatify hey") + "\n\n"
		ss += lipgloss.NewStyle().Width(m.getViewWidth()).Align(lipgloss.Center, lipgloss.Center).Render(`See you then!! (press to "enter" to exit)`)
		s = lipgloss.Place(m.getViewWidth(), m.getViewHeight(), lipgloss.Center, lipgloss.Center, ss)
	default:
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
		if s != "" && lastSpeaker == base.Bot && m.phase == questionPhase {
			s += style.TextInput().Render(m.TextInput.View())
		}
	}
	return s
}
