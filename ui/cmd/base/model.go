package base

import (
	"github.com/JunNishimura/Chatify/config"
	"github.com/JunNishimura/Chatify/utils"
	"github.com/charmbracelet/bubbles/textinput"
)

type Speaker int

const (
	Bot Speaker = iota
	User
)

type Message struct {
	Content string
	Speaker Speaker
}

type Base struct {
	Window       *utils.Window
	Cfg          *config.Config
	TextInput    textinput.Model
	Conversation []*Message
}

func New() (*Base, error) {
	window := utils.NewWindow()

	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}

	return &Base{
		Window:    window,
		Cfg:       cfg,
		TextInput: newTextInput(window.Width),
	}, nil
}

func loadConfig() (*config.Config, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	if err := cfg.Load(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func newTextInput(width int) textinput.Model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = width

	return ti
}
