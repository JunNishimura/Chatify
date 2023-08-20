package hey

import (
	"encoding/json"
	"strings"

	"github.com/JunNishimura/Chatify/ai/functions"
	"github.com/JunNishimura/Chatify/ai/model"
	"github.com/JunNishimura/Chatify/ai/prompt"
	"github.com/JunNishimura/Chatify/auth"
	"github.com/JunNishimura/Chatify/config"
	"github.com/JunNishimura/spotify/v2"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sashabaranov/go-openai"
)

type errMsg struct{ err error }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		inputCmd tea.Cmd
		listCmd  tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.window.UpdateSize()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.state == chatView {
				m.state = recommendationView
			} else {
				m.state = chatView
			}
		case "enter":
			if m.state == chatView {
				answer := m.textInput.Value()
				m.chatCompMessages = append(m.chatCompMessages, openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleUser,
					Content: answer,
				})
				m.conversation = append(m.conversation, answer)

				m.textInput.Reset()

				return m, m.generate
			} else {
				if selectedItem, ok := m.list.SelectedItem().(Item); ok {
					m.selectedItem = selectedItem
					return m, m.playMusic
				}
			}
		}
	case loadConfigMsg:
		m.cfg = msg.cfg
		return m, m.setupSpotify
	case spotifyMsg:
		m.spotifyClient = msg.client
		m.availableGenres = msg.availableGenres
		m.user = msg.user
		return m, m.setupOpenAI
	case openaiMsg:
		m.openaiClient = msg.client

		m.chatCompMessages = []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: prompt.Base,
			},
		}

		m.functions = msg.functions

		return m, m.generate
	case responseMsg:
		functionCall := msg.resp.Choices[0].Message.FunctionCall
		if functionCall != nil {
			return m, m.handleFunctionCall(functionCall)
		} else {
			content := msg.resp.Choices[0].Message.Content
			m.chatCompMessages = append(m.chatCompMessages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: content,
			})
			m.conversation = append(m.conversation, content)
		}
	case chatCompMsg:
		m.chatCompMessages = append(m.chatCompMessages, msg.msg)

		return m, tea.Batch(m.generate, m.recommend)
	case recommendMsg:
		m.recommendItems = msg.items
		m.list = list.New(m.recommendItems, list.NewDefaultDelegate(), 100, 40)
	}

	m.textInput, inputCmd = m.textInput.Update(msg)
	m.list, listCmd = m.list.Update(msg)

	return m, tea.Batch(inputCmd, listCmd)
}

type loadConfigMsg struct {
	cfg *config.Config
}

func (m Model) loadConfig() tea.Msg {
	cfg, err := config.New()
	if err != nil {
		return errMsg{err}
	}

	if err := cfg.Load(); err != nil {
		return errMsg{err}
	}

	return loadConfigMsg{cfg: cfg}
}

type spotifyMsg struct {
	client          *spotify.Client
	availableGenres []string
	user            *model.User
}

func (m Model) setupSpotify() tea.Msg {
	a := auth.NewAuth(m.cfg)

	token := m.cfg.GetToken()

	newToken, err := a.RefreshToken(m.ctx, token)
	if err != nil {
		return errMsg{err}
	}

	// update an access token if it has expired
	if token.AccessToken != newToken.AccessToken {
		if err := m.cfg.SetToken(token); err != nil {
			return errMsg{err}
		}
	}

	client := spotify.New(a.Client(m.ctx, newToken))

	availableGenres, err := client.GetAvailableGenreSeeds(m.ctx)
	if err != nil {
		return errMsg{err}
	}

	curUser, err := client.CurrentUser(m.ctx)
	if err != nil {
		return errMsg{err}
	}

	user := model.NewUser(curUser.DisplayName)

	return spotifyMsg{client, availableGenres, user}
}

type openaiMsg struct {
	client    *openai.Client
	functions []openai.FunctionDefinition
}

func (m Model) setupOpenAI() tea.Msg {
	openaiAPIKey := m.cfg.GetClientValue(config.OpenAIAPIKey)

	return openaiMsg{
		client:    openai.NewClient(openaiAPIKey),
		functions: functions.GetFunctionDefinitions(m.availableGenres),
	}
}

type responseMsg struct{ resp openai.ChatCompletionResponse }

func (m Model) generate() tea.Msg {
	resp, err := m.openaiClient.CreateChatCompletion(
		m.ctx,
		openai.ChatCompletionRequest{
			Model:        openai.GPT3Dot5Turbo,
			Messages:     m.chatCompMessages,
			Functions:    m.functions,
			FunctionCall: "auto",
		},
	)
	if err != nil {
		return errMsg{err}
	}

	return responseMsg{resp}
}

type chatCompMsg struct{ msg openai.ChatCompletionMessage }

func (m Model) handleFunctionCall(functionCall *openai.FunctionCall) tea.Cmd {
	switch functionCall.Name {
	case functions.SetGenresFunctionName:
		return func() tea.Msg {
			result := &struct {
				QualitativeValue string `json:"qualitative_value"`
			}{}
			if err := json.Unmarshal([]byte(functionCall.Arguments), result); err != nil {
				return func() tea.Msg {
					return errMsg{err}
				}
			}

			cleanGenres := strings.TrimSpace(result.QualitativeValue)
			splitGenres := strings.Split(cleanGenres, ",")

			functions.SetGenres(&m.user.MusicOrientation.Genres, splitGenres)

			return chatCompMsg{
				msg: openai.ChatCompletionMessage{
					Name:    functionCall.Name,
					Role:    openai.ChatMessageRoleFunction,
					Content: result.QualitativeValue,
				},
			}
		}
	case functions.SetDanceabilityFunctionName:
		return func() tea.Msg {
			result := &struct {
				QualitativeValue  string  `json:"qualitative_value"`
				QuantitativeValue float64 `json:"quantitative_value"`
			}{}
			if err := json.Unmarshal([]byte(functionCall.Arguments), result); err != nil {
				return func() tea.Msg {
					return errMsg{err}
				}
			}

			functions.SetDanceability(&m.user.MusicOrientation.Danceability, result.QuantitativeValue)

			return chatCompMsg{
				msg: openai.ChatCompletionMessage{
					Name:    functionCall.Name,
					Role:    openai.ChatMessageRoleFunction,
					Content: result.QualitativeValue,
				},
			}
		}
	case functions.SetValenceFunctionName:
		return func() tea.Msg {
			result := &struct {
				QualitativeValue  string  `json:"qualitative_value"`
				QuantitativeValue float64 `json:"quantitative_value"`
			}{}
			if err := json.Unmarshal([]byte(functionCall.Arguments), result); err != nil {
				return func() tea.Msg {
					return errMsg{err}
				}
			}

			functions.SetValence(&m.user.MusicOrientation.Valence, result.QuantitativeValue)

			return chatCompMsg{
				msg: openai.ChatCompletionMessage{
					Name:    functionCall.Name,
					Role:    openai.ChatMessageRoleFunction,
					Content: result.QualitativeValue,
				},
			}
		}
	case functions.SetPopularityFunctionName:
		return func() tea.Msg {
			result := &struct {
				QualitativeValue  string  `json:"qualitative_value"`
				QuantitativeValue float64 `json:"quantitative_value"`
			}{}
			if err := json.Unmarshal([]byte(functionCall.Arguments), result); err != nil {
				return func() tea.Msg {
					return errMsg{err}
				}
			}

			functions.SetPopularity(&m.user.MusicOrientation.Popularity, int(result.QuantitativeValue))

			return chatCompMsg{
				msg: openai.ChatCompletionMessage{
					Name:    functionCall.Name,
					Role:    openai.ChatMessageRoleFunction,
					Content: result.QualitativeValue,
				},
			}
		}
	}

	return nil
}

const RecommendCount = 5

type recommendMsg struct{ items []list.Item }

func (m Model) recommend() tea.Msg {
	if !m.user.MusicOrientation.Genres.HasChanged {
		return nil
	}

	genres := m.user.MusicOrientation.Genres.Value
	if len(genres) > 5 {
		genres = genres[:5]
	}

	seeds := spotify.Seeds{
		Genres: genres,
	}

	trackAttrib := spotify.NewTrackAttributes()
	if m.user.MusicOrientation.Danceability.HasChanged {
		trackAttrib.TargetDanceability(m.user.MusicOrientation.Danceability.Value)
	}
	if m.user.MusicOrientation.Valence.HasChanged {
		trackAttrib.TargetValence(m.user.MusicOrientation.Valence.Value)
	}
	if m.user.MusicOrientation.Popularity.HasChanged {
		trackAttrib.TargetPopularity(m.user.MusicOrientation.Popularity.Value)
	}

	recommendations, err := m.spotifyClient.GetRecommendations(m.ctx, seeds, trackAttrib, spotify.Limit(RecommendCount))
	if err != nil {
		return errMsg{err}
	}

	items := make([]list.Item, 0)
	for _, track := range recommendations.Tracks {
		var artists []string
		for _, artist := range track.Artists {
			artists = append(artists, artist.Name)
		}

		item := Item{
			album:   track.Album.Name,
			artists: artists,
			uri:     track.URI,
		}
		items = append(items, item)
	}

	return recommendMsg{items}
}

func (m Model) playMusic() tea.Msg {
	if err := m.spotifyClient.PlayOpt(m.ctx, &spotify.PlayOptions{
		URIs: []spotify.URI{m.selectedItem.uri},
	}); err != nil {
		return errMsg{err}
	}
	return nil
}
