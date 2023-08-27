package greeting

import (
	"context"

	"github.com/JunNishimura/Chatify/config"
	"github.com/JunNishimura/Chatify/utils"
	"github.com/JunNishimura/spotify/v2"
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/charmbracelet/bubbles/textinput"
)

var (
	qaListTemplate = []*QA{
		{
			question:    `May I ask for your "Spotify ID"?`,
			answer:      "",
			placeholder: `Please enter your "Spotify ID".`,
		},
		{
			question:    `Next, may I ask for your "Spotify Secret"?`,
			answer:      "",
			placeholder: `Please enter your "Spotify Secret".`,
		},
		{
			question:    `Finally, may I ask for your "OpenAI API key"?`,
			answer:      "",
			placeholder: `Please enter your "OpenAI API key".`,
		},
	}
	conversationTemplate = []*Message{
		{
			speaker: Bot,
			content: heredoc.Docf(`
				Hi there, I'm Chatify! I want to know about you.
				%s
			`, qaListTemplate[0].question),
		},
	}
	confKeyList = []config.ConfKey{
		config.SpotifyIDKey,
		config.SpotifySecretKey,
		config.OpenAIAPIKey,
	}
)

type QA struct {
	question    string
	answer      string
	placeholder string
}

type Speaker int

const (
	Bot Speaker = iota
	User
)

type Message struct {
	content string
	speaker Speaker
}

type Phase int

const (
	questionPhase Phase = iota
	authPhase
	devicePhase
	completePhase
)

type Model struct {
	ctx           context.Context
	window        *utils.Window
	textInput     textinput.Model
	cfg           *config.Config
	phase         Phase
	questionIndex int
	qaList        []*QA
	conversation  []*Message
	user          *spotify.PrivateUser
	spotifyClient *spotify.Client
}

func NewModel() *Model {
	window := utils.NewWindow()

	return &Model{
		ctx:           context.Background(),
		window:        window,
		textInput:     newTextInput(window.Width),
		phase:         questionPhase,
		questionIndex: 0,
		qaList:        qaListTemplate,
		conversation:  conversationTemplate,
	}
}

func newTextInput(width int) textinput.Model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = width
	ti.Placeholder = qaListTemplate[0].placeholder

	return ti
}
