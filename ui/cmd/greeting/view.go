package greeting

import (
	"fmt"

	"github.com/JunNishimura/Chatify/ui/style"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	return lipgloss.Place(m.window.Width, m.window.Height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			style.AsciiArt().Render(asciiArtView),
			style.ChatFocused(m.getViewWidth(), m.getViewHeight()).Render(m.chatView()),
		),
	)
}

func (m *Model) getViewWidth() int {
	halfWidth := m.window.Width / 2
	return halfWidth
}

func (m *Model) getViewHeight() int {
	halfHeight := m.window.Height * 2 / 3
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
		if s != "" && lastSpeaker == Bot && m.phase == questionPhase {
			s += style.TextInput().Render(m.textInput.View())
		}
	}
	return s
}
