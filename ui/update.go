package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/JunNishimura/Chatify/auth"
	"github.com/JunNishimura/Chatify/config"
	"github.com/JunNishimura/spotify/v2"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type loadConfigMsg struct{ cfg *config.Config }
type spotifyUserMsg struct{ user *spotify.PrivateUser }
type errMsg struct{ err error }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			if m.greetingDone {
				return m, tea.Quit
			} else if m.setConfigDone {
				return m, m.Authorize
			}

			m.qaList[m.index].answer = m.textarea.Value()

			m.index++
			if m.index == len(m.qaList) {
				m.qaDone = true
				m.displayMessages = append(m.displayMessages, m.senderStyle.Render("Chatify: ")+"Thank you so much!")
				m.viewport.SetContent(strings.Join(m.displayMessages, "\n"))

				return m, tea.Batch(
					m.setClientConfig(m.confKeyList[m.writeIndex], qaListTemplate[m.writeIndex].answer),
				)
			}

			m.displayMessages = append(
				m.displayMessages,
				m.senderStyle.Render("You: ")+m.textarea.Value(),
				m.senderStyle.Render("Chatify: ")+m.qaList[m.index].question,
			)
			m.viewport.SetContent(strings.Join(m.displayMessages, "\n"))
			m.viewport.GotoBottom()

			m.textarea.Reset()
			m.textarea.Placeholder = m.qaList[m.index].placeholder
		}
	case writeClientConfigMsg:
		m.writeIndex++
		if m.writeIndex < len(m.confKeyList) {
			return m, m.setClientConfig(m.confKeyList[m.writeIndex], m.qaList[m.writeIndex].answer)
		}
		m.setConfigDone = true
	case loadConfigMsg:
		m.cfg = msg.cfg
		return m, textarea.Blink
	case spotifyUserMsg:
		m.user = msg.user
		m.greetingDone = true
		return m, tea.Batch(tiCmd, vpCmd)
	case errMsg:
		m.err = msg.err
		return m, tea.Quit
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m *Model) loadConfig() tea.Msg {
	cfg, err := config.New()
	if err != nil {
		return errMsg{err: err}
	}

	if err := cfg.Load(); err != nil {
		return errMsg{err: err}
	}

	return loadConfigMsg{cfg: cfg}
}

type writeClientConfigMsg string

func (m *Model) setClientConfig(key config.ConfKey, value any) tea.Cmd {
	start := time.Now()
	if err := m.cfg.Set(key, value); err != nil {
		return func() tea.Msg {
			return errMsg{err: err}
		}
	}
	elapsed := time.Since(start)

	return tea.Tick(elapsed, func(t time.Time) tea.Msg {
		return writeClientConfigMsg(key)
	})
}

func (m *Model) Authorize() tea.Msg {
	authClient := auth.New(m.cfg)

	authClient.Authorize()

	spotifyClient := <-authClient.SpotifyChannel

	user, err := spotifyClient.CurrentUser(context.Background())
	if err != nil {
		return errMsg{err: err}
	}

	return spotifyUserMsg{user: user}
}
