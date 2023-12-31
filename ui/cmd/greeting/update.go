package greeting

import (
	"context"
	"errors"

	"github.com/JunNishimura/Chatify/auth"
	"github.com/JunNishimura/Chatify/config"
	"github.com/JunNishimura/Chatify/ui/cmd/base"
	"github.com/JunNishimura/spotify/v2"
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg struct{ err error }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var inputCmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.err != nil {
				return m, tea.Quit
			}

			switch m.phase {
			case questionPhase:
				m.qaList[m.questionIndex].answer = m.TextInput.Value()
				m.Conversation = append(m.Conversation, &base.Message{
					Speaker: base.User,
					Content: m.TextInput.Value(),
				})
				return m, tea.Batch(inputCmd, m.setConfig)
			case authPhase:
				return m, m.authorize
			case devicePhase:
				return m, m.getDevice
			case completePhase:
				return m, tea.Quit
			}
		}
	case questionCompMsg:
		if msg.isDone {
			m.phase = authPhase
			m.Conversation = append(m.Conversation, &base.Message{
				Speaker: base.Bot,
				Content: "Please press enter to authorize",
			})
		}
	case spotifyMsg:
		m.user = msg.user
		m.spotifyClient = msg.client
		m.phase = devicePhase
		m.Conversation = append(m.Conversation, &base.Message{
			Speaker: base.Bot,
			Content: "Please open Spotify app to get device ID and press enter",
		})
	case deviceMsg:
		m.phase = completePhase
	case errMsg:
		m.err = msg.err
	}

	m.TextInput, inputCmd = m.TextInput.Update(msg)

	return m, inputCmd
}

type questionCompMsg struct {
	isDone bool
}

func (m *Model) setConfig() tea.Msg {
	if err := m.Cfg.Set(confKeyList[m.questionIndex], m.qaList[m.questionIndex].answer); err != nil {
		return errMsg{err}
	}

	m.questionIndex++
	m.TextInput.Reset()
	if m.questionIndex == len(m.qaList) {
		return questionCompMsg{
			isDone: true,
		}
	}
	m.Conversation = append(m.Conversation, &base.Message{
		Speaker: base.Bot,
		Content: m.qaList[m.questionIndex].question,
	})
	return questionCompMsg{
		isDone: false,
	}
}

type spotifyMsg struct {
	user   *spotify.PrivateUser
	client *spotify.Client
}

func (m *Model) authorize() tea.Msg {
	authClient := auth.NewClient(m.Cfg)

	authClient.Authorize()

	spotifyClient := <-authClient.SpotifyChannel

	user, err := spotifyClient.CurrentUser(context.Background())
	if err != nil {
		return errMsg{err: err}
	}

	if err := m.Cfg.Set(config.UserIDKey, user.ID); err != nil {
		return errMsg{err: err}
	}

	return spotifyMsg{
		user:   user,
		client: spotifyClient,
	}
}

type deviceMsg struct{ deviceID string }

func (m *Model) getDevice() tea.Msg {
	devices, err := m.spotifyClient.PlayerDevices(context.Background())
	if err != nil {
		return errMsg{err}
	}
	if len(devices) == 0 {
		return errMsg{
			err: errors.New("fail to get device"),
		}
	}

	deviceID := devices[0].ID.String()
	if err := m.Cfg.Set(config.DeviceID, deviceID); err != nil {
		return errMsg{err}
	}

	return deviceMsg{deviceID}
}
