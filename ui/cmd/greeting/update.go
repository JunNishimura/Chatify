package greeting

import (
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
			switch m.phase {
			case questionPhase:
				m.qaList[m.questionIndex].answer = m.base.TextInput.Value()
				m.base.Conversation = append(m.base.Conversation, &base.Message{
					Speaker: base.User,
					Content: m.base.TextInput.Value(),
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
			m.base.Conversation = append(m.base.Conversation, &base.Message{
				Speaker: base.Bot,
				Content: "Please press enter to authorize",
			})
		}
	case spotifyMsg:
		m.user = msg.user
		m.spotifyClient = msg.client
		m.phase = devicePhase
		m.base.Conversation = append(m.base.Conversation, &base.Message{
			Speaker: base.Bot,
			Content: "Please open Spotify app to get device ID and press enter",
		})
	case deviceMsg:
		m.phase = completePhase
	case errMsg:
		m.err = msg.err
	}

	m.base.TextInput, inputCmd = m.base.TextInput.Update(msg)

	return m, inputCmd
}

type questionCompMsg struct {
	isDone bool
}

func (m *Model) setConfig() tea.Msg {
	if err := m.base.Cfg.Set(confKeyList[m.questionIndex], m.qaList[m.questionIndex].answer); err != nil {
		return errMsg{err}
	}

	m.questionIndex++
	m.base.TextInput.Reset()
	if m.questionIndex == len(m.qaList) {
		return questionCompMsg{
			isDone: true,
		}
	}
	m.base.Conversation = append(m.base.Conversation, &base.Message{
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
	authClient := auth.NewClient(m.base.Cfg)

	authClient.Authorize()

	spotifyClient := <-authClient.SpotifyChannel

	user, err := spotifyClient.CurrentUser(m.ctx)
	if err != nil {
		return errMsg{err: err}
	}

	return spotifyMsg{
		user:   user,
		client: spotifyClient,
	}
}

type deviceMsg struct{ deviceID string }

func (m *Model) getDevice() tea.Msg {
	devices, err := m.spotifyClient.PlayerDevices(m.ctx)
	if err != nil {
		return errMsg{err}
	}
	if len(devices) == 0 {
		return errMsg{
			err: errors.New("fail to get device"),
		}
	}

	deviceID := devices[0].ID.String()
	if err := m.base.Cfg.Set(config.DeviceID, deviceID); err != nil {
		return errMsg{err}
	}

	return deviceMsg{deviceID}
}
