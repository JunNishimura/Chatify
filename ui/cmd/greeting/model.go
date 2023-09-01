package greeting

import (
	"context"

	"github.com/JunNishimura/Chatify/config"
	"github.com/JunNishimura/Chatify/ui/cmd/base"
	"github.com/JunNishimura/spotify/v2"
	"github.com/MakeNowJust/heredoc/v2"
)

var (
	qaListTemplate = []*QA{
		{
			question: `May I ask for your "Spotify ID"?`,
			answer:   "",
		},
		{
			question: `Next, may I ask for your "Spotify Secret"?`,
			answer:   "",
		},
		{
			question: `Finally, may I ask for your "OpenAI API key"?`,
			answer:   "",
		},
	}
	conversationTemplate = []*base.Message{
		{
			Speaker: base.Bot,
			Content: heredoc.Docf(`
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
	question string
	answer   string
}

type Phase int

const (
	questionPhase Phase = iota
	authPhase
	devicePhase
	completePhase
)

type Model struct {
	ctx context.Context
	*base.Base
	phase         Phase
	questionIndex int
	qaList        []*QA
	user          *spotify.PrivateUser
	spotifyClient *spotify.Client
	err           error
}

func NewModel(ctx context.Context, port string) (*Model, error) {
	base, err := base.New()
	if err != nil {
		return nil, err
	}

	base.Conversation = conversationTemplate

	if err := base.Cfg.Set(config.PortKey, port); err != nil {
		return nil, err
	}

	return &Model{
		ctx:           ctx,
		Base:          base,
		phase:         questionPhase,
		questionIndex: 0,
		qaList:        qaListTemplate,
	}, nil
}
