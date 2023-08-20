package hey

import (
	"context"
	"strings"

	"github.com/JunNishimura/Chatify/ai/model"
	"github.com/JunNishimura/Chatify/config"
	"github.com/JunNishimura/Chatify/ui/hey/style"
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
	chatCompMessages []openai.ChatCompletionMessage
	conversation     []*Message
	functions        []openai.FunctionDefinition
	availableGenres  []string
	recommendItems   []list.Item
}

func NewModel() Model {
	return Model{
		ctx:       context.Background(),
		window:    utils.NewWindow(),
		textInput: newTextInput(),
		list:      newListModel([]list.Item{}, 0, 0),
	}
}

const ListTitle = "Chatify's recommendation"

func newListModel(items []list.Item, width, height int) list.Model {
	newList := list.New(items, list.NewDefaultDelegate(), width, height)
	newList.Title = ListTitle
	newList.Styles.Title.Background(lipgloss.Color(style.HighlightColor))
	return newList
}

func newTextInput() textinput.Model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 100

	return ti
}
