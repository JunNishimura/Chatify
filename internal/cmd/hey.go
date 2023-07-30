/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/JunNishimura/Chatify/internal/functions"
	"github.com/JunNishimura/Chatify/internal/object"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"github.com/zmb3/spotify/v2"
)

const (
	initPrompt = `Below is a conversation with an AI chatbot.

The bot analyzes the music the interlocutor is currently seeking through the conversation and suggests music recommendations based on the results of the analysis.

The bot analyzes the music orientation of the music the interlocutor is currently seeking by breaking it down into the following elements.
1. Genre
Music genres. For example, j-pop, techno, acoustic, folk
2. Danceability
Danceability describes how suitable a track is for dancing based on a combination of musical elements including tempo, rhythm stability, beat strength, and overall regularity. A value of 0.0 is least danceable and 1.0 is most danceable.
3. Valence
A measure from 0.0 to 1.0 describing the musical positiveness conveyed by a track. Tracks with high valence sound more positive (e.g. happy, cheerful, euphoric), while tracks with low valence sound more negative (e.g. sad, depressed, angry).
4. Popularity
A measure from 0 to 100 describing how much the track is popular. Tracks with high popularity is more popular.

Once all factors have been determined, the bot will suggest music recommendations to the interlocutor based on the information obtained.

There are some points to note when asking questions.
The possible values for the analysis elements Danceability, Valence, and Popularity are numerical values such as 0.6, 
but do not ask questions that force the interlocutor to directly answer with a numerical value, 
such as "How much is your danceability from 0 to 1?
Instead, ask a question to analyze how much daceability music the interlocutor is looking for,
such as "Do you want to be energetic?”. 
Then, guess the specific numerical value of the element from the interlocutor's answer.
For example, "I'm depressed and I want to get better" to which the response might be something like,
"I guess the daceability is 0.8”.
Also, limit the number of questions the bot asks the interlocutor in one conversation to one.

Please begin with the first question.`
)

func NewHeyCommand(ctx context.Context, client *spotify.Client, openaiApiKey string) *cobra.Command {
	return &cobra.Command{
		Use:   "hey",
		Short: "start conversation with chatify",
		Long:  "start conversation with chatify",
		RunE: func(cmd *cobra.Command, args []string) error {
			// get client
			openAIClient := openai.NewClient(openaiApiKey)

			messages := []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: initPrompt,
				},
			}

			availableGenres, err := client.GetAvailableGenreSeeds(ctx)
			if err != nil {
				return fmt.Errorf("fail to get available genres: %v", err)
			}

			musicOrientationInfo := object.NewMusicOrientationInfo()

			functionTemplate := functions.GetTemplate(availableGenres)

			scanner := bufio.NewScanner(os.Stdin)
			for {
				resp, err := openAIClient.CreateChatCompletion(
					ctx,
					openai.ChatCompletionRequest{
						Model:        openai.GPT3Dot5Turbo,
						Messages:     messages,
						Functions:    functionTemplate,
						FunctionCall: "auto",
					},
				)
				if err != nil {
					return fmt.Errorf("chat completion error: %v", err)
				}

				// show AI response
				functionCall := resp.Choices[0].Message.FunctionCall
				if functionCall != nil {
					switch functionCall.Name {
					case functions.RecommendFunction:
						args := make(map[string]any)
						if err := json.Unmarshal([]byte(functionCall.Arguments), &args); err != nil {
							return fmt.Errorf("fail to unmarshal json: %v", err)
						}

						genres := args["genres"].(string)
						danceability := args["danceability"].(float64)
						valence := args["valence"].(float64)
						popularity := args["popularity"].(int)

						recommendations, err := functions.Recommend(ctx, client, genres, danceability, valence, popularity)
						if err != nil {
							return err
						}

						messages = append(messages, openai.ChatCompletionMessage{
							Name:    functionCall.Name,
							Role:    openai.ChatMessageRoleFunction,
							Content: recommendations,
						})
					case functions.SetGenresFunction:
						args := make(map[string]string)
						if err := json.Unmarshal([]byte(functionCall.Arguments), &args); err != nil {
							return fmt.Errorf("fail to unmarshal json: %v", err)
						}
						genres := args["value"]
						functions.SetGenres(musicOrientationInfo, genres)

						messages = append(messages, openai.ChatCompletionMessage{
							Name:    functionCall.Name,
							Role:    openai.ChatMessageRoleFunction,
							Content: genres,
						})
					case functions.SetDanceabitliyFunction:
						args := make(map[string]any)
						if err := json.Unmarshal([]byte(functionCall.Arguments), &args); err != nil {
							return fmt.Errorf("fail to unmarshal json: %v", err)
						}

						danceability := args["value"].(float64)
						description := args["subvalue"].(string)

						functions.SetDanceability(musicOrientationInfo, danceability)

						messages = append(messages, openai.ChatCompletionMessage{
							Name:    functionCall.Name,
							Role:    openai.ChatMessageRoleFunction,
							Content: description,
						})
					case functions.SetValenceFunction:
						args := make(map[string]any)
						if err := json.Unmarshal([]byte(functionCall.Arguments), &args); err != nil {
							return fmt.Errorf("fail to unmarshal json: %v", err)
						}

						valence := args["value"].(float64)
						description := args["subvalue"].(string)

						functions.SetValence(musicOrientationInfo, valence)

						messages = append(messages, openai.ChatCompletionMessage{
							Name:    functionCall.Name,
							Role:    openai.ChatMessageRoleFunction,
							Content: description,
						})
					case functions.SetPopularityFunction:
						args := make(map[string]any)
						if err := json.Unmarshal([]byte(functionCall.Arguments), &args); err != nil {
							return fmt.Errorf("fail to unmarshal json: %v", err)
						}

						popularity := args["value"].(int)
						description := args["subvalue"].(string)

						functions.SetPopularity(musicOrientationInfo, popularity)

						messages = append(messages, openai.ChatCompletionMessage{
							Name:    functionCall.Name,
							Role:    openai.ChatMessageRoleFunction,
							Content: description,
						})
					}
				} else {
					respMessage := resp.Choices[0].Message.Content

					// termination check
					if strings.Contains(respMessage, "<end>") {
						sp := strings.Split(respMessage, "<end>")
						recommendMessage := strings.Join(sp, "")
						fmt.Println(recommendMessage)
						break
					}

					fmt.Println(respMessage)
					messages = append(messages, openai.ChatCompletionMessage{
						Role:    openai.ChatMessageRoleAssistant,
						Content: respMessage,
					})

					// get user answer
					fmt.Printf("> ")
					scanner.Scan()
					messages = append(messages, openai.ChatCompletionMessage{
						Role:    openai.ChatMessageRoleUser,
						Content: scanner.Text(),
					})
				}

				// when music orientation info is all set, add info to messages manually
				if musicOrientationInfo.IsSet() {
					functionTemplate = append(functionTemplate, functions.RecommendFunctionDefinition)
					messages = append(messages, openai.ChatCompletionMessage{
						Role:    openai.ChatMessageRoleUser,
						Content: musicOrientationInfo.String(),
					})
				}
			}

			return nil
		},
	}
}
