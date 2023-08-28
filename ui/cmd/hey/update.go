package hey

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/JunNishimura/Chatify/ai/functions"
	"github.com/JunNishimura/spotify/v2"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sashabaranov/go-openai"
)

type errMsg struct{ err error }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			m.state = m.state.Switch()
		case "enter":
			switch m.state {
			case chatView:
				answer := m.textInput.Value()
				m.chatCompMessages = append(m.chatCompMessages, openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleUser,
					Content: answer,
				})
				m.conversation = append(m.conversation, &Message{content: fmt.Sprintf("> %s", answer), speaker: User})

				m.textInput.Reset()

				if m.questionIndex == 0 {
					// genres don't have to be converted into quantitative values
					m.functionCall = functions.Call{
						Name: string(functions.List[m.questionIndex]),
					}
					return m, m.generate
				}
				return m, m.guessQuantitativeValue
			case recommendationView:
				if selectedItem, ok := m.list.SelectedItem().(Item); ok {
					m.selectedItem = selectedItem
					return m, m.playMusic
				}
			}
		}
	case generationMsg:
		functionCall := msg.resp.Choices[0].Message.FunctionCall
		if functionCall != nil {
			return m, m.handleFunctionCall(functionCall)
		} else {
			content := msg.resp.Choices[0].Message.Content
			m.chatCompMessages = append(m.chatCompMessages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: content,
			})
			m.conversation = append(m.conversation, &Message{content: content, speaker: Bot})
		}
	case guessedMsg:
		content := msg.resp.Choices[0].Message.Content
		m.functionCall = functions.Call{
			Name: string(functions.List[m.questionIndex]),
		}
		m.chatCompMessages = append(m.chatCompMessages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: content,
		})
		return m, m.generate
	case functionCallMsg:
		m.chatCompMessages = append(m.chatCompMessages, msg.msg)
		m.functionCall = "auto"
		m.questionIndex++
		return m, tea.Batch(m.generate, m.recommend)
	case recommendMsg:
		m.recommendItems = msg.items
		m.list = newListModel(m.recommendItems, m.getViewWidth(), m.getViewHeight())
	case errMsg:
		m.err = msg.err
	}

	m.textInput, inputCmd = m.textInput.Update(msg)
	m.list, listCmd = m.list.Update(msg)

	return m, tea.Batch(inputCmd, listCmd)
}

type generationMsg struct{ resp openai.ChatCompletionResponse }

func (m *Model) generate() tea.Msg {
	resp, err := m.openaiClient.CreateChatCompletion(
		m.ctx,
		openai.ChatCompletionRequest{
			Model:        openai.GPT3Dot5Turbo16K0613,
			Messages:     m.chatCompMessages,
			Functions:    m.functions,
			FunctionCall: m.functionCall,
		},
	)
	if err != nil {
		return errMsg{err}
	}

	return generationMsg{resp}
}

type guessedMsg struct{ resp openai.ChatCompletionResponse }

// convert user input into quantitative value
func (m *Model) guessQuantitativeValue() tea.Msg {
	resp, err := m.openaiClient.CreateChatCompletion(
		m.ctx,
		openai.ChatCompletionRequest{
			Model:        openai.GPT3Dot5Turbo16K0613,
			Messages:     m.chatCompMessages,
			Functions:    m.functions,
			FunctionCall: "none",
		},
	)
	if err != nil {
		return errMsg{err}
	}

	return guessedMsg{resp}
}

type functionCallMsg struct{ msg openai.ChatCompletionMessage }

func (m *Model) handleFunctionCall(functionCall *openai.FunctionCall) tea.Cmd {
	switch functionCall.Name {
	case string(functions.SetGenresFunctionName):
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

			return functionCallMsg{
				msg: openai.ChatCompletionMessage{
					Name:    functionCall.Name,
					Role:    openai.ChatMessageRoleFunction,
					Content: result.QualitativeValue,
				},
			}
		}
	case string(functions.SetDanceabilityFunctionName):
		return func() tea.Msg {
			result := &struct {
				QuantitativeValue float64 `json:"quantitative_value"`
			}{}
			if err := json.Unmarshal([]byte(functionCall.Arguments), result); err != nil {
				return func() tea.Msg {
					return errMsg{err}
				}
			}

			functions.SetDanceability(&m.user.MusicOrientation.Danceability, result.QuantitativeValue)

			return functionCallMsg{
				msg: openai.ChatCompletionMessage{
					Name:    functionCall.Name,
					Role:    openai.ChatMessageRoleFunction,
					Content: fmt.Sprintf("Danceability: %f", result.QuantitativeValue),
				},
			}
		}
	case string(functions.SetValenceFunctionName):
		return func() tea.Msg {
			result := &struct {
				QuantitativeValue float64 `json:"quantitative_value"`
			}{}
			if err := json.Unmarshal([]byte(functionCall.Arguments), result); err != nil {
				return func() tea.Msg {
					return errMsg{err}
				}
			}

			functions.SetValence(&m.user.MusicOrientation.Valence, result.QuantitativeValue)

			return functionCallMsg{
				msg: openai.ChatCompletionMessage{
					Name:    functionCall.Name,
					Role:    openai.ChatMessageRoleFunction,
					Content: fmt.Sprintf("Valence: %f", result.QuantitativeValue),
				},
			}
		}
	case string(functions.SetPopularityFunctionName):
		return func() tea.Msg {
			result := &struct {
				QuantitativeValue float64 `json:"quantitative_value"`
			}{}
			if err := json.Unmarshal([]byte(functionCall.Arguments), result); err != nil {
				return func() tea.Msg {
					return errMsg{err}
				}
			}

			functions.SetPopularity(&m.user.MusicOrientation.Popularity, int(result.QuantitativeValue))

			return functionCallMsg{
				msg: openai.ChatCompletionMessage{
					Name:    functionCall.Name,
					Role:    openai.ChatMessageRoleFunction,
					Content: fmt.Sprintf("Popularity: %f", result.QuantitativeValue),
				},
			}
		}
	case string(functions.SetAcousticnessFunctionName):
		return func() tea.Msg {
			result := &struct {
				QuantitativeValue float64 `json:"quantitative_value"`
			}{}
			if err := json.Unmarshal([]byte(functionCall.Arguments), result); err != nil {
				return func() tea.Msg {
					return errMsg{err}
				}
			}

			functions.SetAcousticness(&m.user.MusicOrientation.Acousticness, result.QuantitativeValue)

			return functionCallMsg{
				msg: openai.ChatCompletionMessage{
					Name:    functionCall.Name,
					Role:    openai.ChatMessageRoleFunction,
					Content: fmt.Sprintf("Acousticness: %f", result.QuantitativeValue),
				},
			}
		}
	case string(functions.SetEnergyFunctionName):
		return func() tea.Msg {
			result := &struct {
				QuantitativeValue float64 `json:"quantitative_value"`
			}{}
			if err := json.Unmarshal([]byte(functionCall.Arguments), result); err != nil {
				return func() tea.Msg {
					return errMsg{err}
				}
			}

			functions.SetEnergy(&m.user.MusicOrientation.Energy, result.QuantitativeValue)

			return functionCallMsg{
				msg: openai.ChatCompletionMessage{
					Name:    functionCall.Name,
					Role:    openai.ChatMessageRoleFunction,
					Content: fmt.Sprintf("Energy: %f", result.QuantitativeValue),
				},
			}
		}
	case string(functions.SetInstrumentalnessFunctionName):
		return func() tea.Msg {
			result := &struct {
				QuantitativeValue float64 `json:"quantitative_value"`
			}{}
			if err := json.Unmarshal([]byte(functionCall.Arguments), result); err != nil {
				return func() tea.Msg {
					return errMsg{err}
				}
			}

			functions.SetInstrumentalness(&m.user.MusicOrientation.Instrumentalness, result.QuantitativeValue)

			return functionCallMsg{
				msg: openai.ChatCompletionMessage{
					Name:    functionCall.Name,
					Role:    openai.ChatMessageRoleFunction,
					Content: fmt.Sprintf("Instrumentalness: %f", result.QuantitativeValue),
				},
			}
		}
	case string(functions.SetLivenessFunctionaName):
		return func() tea.Msg {
			result := &struct {
				QuantitativeValue float64 `json:"quantitative_value"`
			}{}
			if err := json.Unmarshal([]byte(functionCall.Arguments), result); err != nil {
				return func() tea.Msg {
					return errMsg{err}
				}
			}

			functions.SetLiveness(&m.user.MusicOrientation.Liveness, result.QuantitativeValue)

			return functionCallMsg{
				msg: openai.ChatCompletionMessage{
					Name:    functionCall.Name,
					Role:    openai.ChatMessageRoleFunction,
					Content: fmt.Sprintf("Liveness: %f", result.QuantitativeValue),
				},
			}
		}
	case string(functions.SetSpeechinessFunctionName):
		return func() tea.Msg {
			result := &struct {
				QuantitativeValue float64 `json:"quantitative_value"`
			}{}
			if err := json.Unmarshal([]byte(functionCall.Arguments), result); err != nil {
				return func() tea.Msg {
					return errMsg{err}
				}
			}

			functions.SetSpeechiness(&m.user.MusicOrientation.Speechiness, result.QuantitativeValue)

			return functionCallMsg{
				msg: openai.ChatCompletionMessage{
					Name:    functionCall.Name,
					Role:    openai.ChatMessageRoleFunction,
					Content: fmt.Sprintf("Speechiness: %f", result.QuantitativeValue),
				},
			}
		}
	}

	return nil
}

const RecommendCount = 25

type recommendMsg struct{ items []list.Item }

func (m *Model) recommend() tea.Msg {
	if !m.user.MusicOrientation.HasOneChanged() {
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
	if m.user.MusicOrientation.Acousticness.HasChanged {
		trackAttrib.TargetAcousticness(m.user.MusicOrientation.Acousticness.Value)
	}
	if m.user.MusicOrientation.Energy.HasChanged {
		trackAttrib.TargetEnergy(m.user.MusicOrientation.Energy.Value)
	}
	if m.user.MusicOrientation.Instrumentalness.HasChanged {
		trackAttrib.TargetInstrumentalness(m.user.MusicOrientation.Instrumentalness.Value)
	}
	if m.user.MusicOrientation.Liveness.HasChanged {
		trackAttrib.TargetLiveness(m.user.MusicOrientation.Liveness.Value)
	}
	if m.user.MusicOrientation.Speechiness.HasChanged {
		trackAttrib.TargetSpeechiness(m.user.MusicOrientation.Speechiness.Value)
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

func (m *Model) playMusic() tea.Msg {
	if err := m.spotifyClient.PlayOpt(m.ctx, &spotify.PlayOptions{
		URIs: []spotify.URI{m.selectedItem.uri},
	}); err != nil {
		return errMsg{err}
	}
	return nil
}
