package greeting

import (
	"context"
	"strings"

	"github.com/JunNishimura/Chatify/config"
	"github.com/JunNishimura/spotify/v2"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

const (
	TextAreaWidth     = 200
	TextAreaHeight    = 3
	TextAreaCharLimit = 200
	ViewportWidth     = 100
	ViewportHeight    = 8
)

var (
	senderStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
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
)

type Model struct {
	ctx             context.Context
	index           int
	writeIndex      int
	cfg             *config.Config
	textarea        textarea.Model
	confKeyList     []config.ConfKey
	displayMessages []string
	qaList          []*QA
	viewport        viewport.Model
	user            *spotify.PrivateUser
	spotifyClient   *spotify.Client
	senderStyle     lipgloss.Style
	qaDone          bool
	setConfigDone   bool
	greetingDone    bool
	err             error
}

type QA struct {
	question    string
	answer      string
	placeholder string
}

func NewModel() Model {
	greetings := []string{
		senderStyle.Render("Chatify: ") + "Hi there, I'm Chatify!",
		"         I want to know three things.",
		"         " + qaListTemplate[0].question,
	}

	return Model{
		ctx:             context.Background(),
		index:           0,
		writeIndex:      0,
		confKeyList:     []config.ConfKey{config.SpotifyIDKey, config.SpotifySecretKey, config.OpenAIAPIKey},
		displayMessages: greetings,
		qaList:          qaListTemplate,
		textarea:        newTextArea(),
		viewport:        newViewport(greetings),
		senderStyle:     senderStyle,
		err:             nil,
	}
}

func newTextArea() textarea.Model {
	ta := textarea.New()
	ta.Placeholder = qaListTemplate[0].placeholder
	ta.Focus()
	ta.CharLimit = TextAreaCharLimit
	ta.SetHeight(TextAreaHeight)
	ta.SetWidth(TextAreaWidth)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	return ta
}

func newViewport(greetings []string) viewport.Model {
	vp := viewport.New(ViewportWidth, ViewportHeight)
	vp.SetContent(strings.Join(greetings, "\n"))

	return vp
}
