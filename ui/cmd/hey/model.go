package hey

import (
	"context"
	"strings"

	"github.com/JunNishimura/Chatify/ai/functions"
	"github.com/JunNishimura/Chatify/ai/model"
	"github.com/JunNishimura/Chatify/ai/prompt"
	"github.com/JunNishimura/Chatify/auth"
	"github.com/JunNishimura/Chatify/config"
	"github.com/JunNishimura/Chatify/ui/style"
	"github.com/JunNishimura/Chatify/utils"
	"github.com/JunNishimura/spotify/v2"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/sashabaranov/go-openai"
)

type sessionState uint

const (
	chatView sessionState = iota
	recommendationView
)

func (s sessionState) Change() {
	switch s {
	case chatView:
		s = recommendationView
	case recommendationView:
		s = chatView
	}
}

func (s sessionState) String() string {
	switch s {
	case chatView:
		return "chat"
	case recommendationView:
		return "recommendation"
	default:
		return ""
	}
}

type Item struct {
	album   string
	artists []string
	uri     spotify.URI
}

func (i Item) Title() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(style.White)).
		Render(i.album)
}
func (i Item) Description() string { return strings.Join(i.artists, ", ") }
func (i Item) FilterValue() string { return i.album }

type Speaker int

const (
	Bot Speaker = iota
	User
)

type Message struct {
	content string
	speaker Speaker
}

type Model struct {
	ctx              context.Context
	window           *utils.Window
	state            sessionState
	textInput        textinput.Model
	list             list.Model
	selectedItem     Item
	cfg              *config.Config
	user             *model.User
	spotifyClient    *spotify.Client
	openaiClient     *openai.Client
	questionIndex    int
	functionCall     any
	chatCompMessages []openai.ChatCompletionMessage
	conversation     []*Message
	functions        []openai.FunctionDefinition
	recommendItems   []list.Item
	err              error
}

func NewModel() (*Model, error) {
	ctx := context.Background()

	window := utils.NewWindow()

	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}

	openAIAPIkey := cfg.GetClientValue(config.OpenAIAPIKey)
	openAIAPIclient := openai.NewClient(openAIAPIkey)

	spotifyClient, err := getSpotifyClient(ctx, cfg)
	if err != nil {
		return nil, err
	}

	availableGenres, err := spotifyClient.GetAvailableGenreSeeds(ctx)
	if err != nil {
		return nil, err
	}

	user, err := getUser(ctx, spotifyClient)
	if err != nil {
		return nil, err
	}

	return &Model{
		ctx:           ctx,
		window:        window,
		textInput:     newTextInput(window.Width),
		list:          newListModel([]list.Item{}, 0, 0),
		cfg:           cfg,
		user:          user,
		spotifyClient: spotifyClient,
		openaiClient:  openAIAPIclient,
		questionIndex: 0,
		functionCall:  "auto",
		chatCompMessages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: prompt.Base,
			},
		},
		functions: functions.GetFunctionDefinitions(availableGenres),
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

func getSpotifyClient(ctx context.Context, cfg *config.Config) (*spotify.Client, error) {
	a := auth.NewAuth(cfg)

	token := cfg.GetToken()

	newToken, err := a.RefreshToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// update an access token if it has expired
	if token.AccessToken != newToken.AccessToken {
		if err := cfg.SetToken(token); err != nil {
			return nil, err
		}
	}

	client := spotify.New(a.Client(ctx, newToken))

	return client, nil
}

func getUser(ctx context.Context, client *spotify.Client) (*model.User, error) {
	curUser, err := client.CurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	user := model.NewUser(curUser.DisplayName)

	return user, nil
}

const listTitle = "Chatify's recommendation"

func newListModel(items []list.Item, width, height int) list.Model {
	newDelegate := list.NewDefaultDelegate()
	newDelegate.Styles = style.List()

	newList := list.New(items, newDelegate, width, height)
	newList.Title = listTitle
	newList.Styles.Title.Background(lipgloss.Color(style.HighlightColor))
	return newList
}

func newTextInput(width int) textinput.Model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = width

	return ti
}
