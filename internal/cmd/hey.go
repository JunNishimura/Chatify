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
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"github.com/zmb3/spotify/v2"
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
					Role: openai.ChatMessageRoleSystem,
					Content: `Below is a conversation with an AI chatbot.

The bot suggests recommended tracks to the interlocutor.
Through a question-and-answer session, the bot will suggest recommended tracks based on information obtained. 
However, please observe the following rules.

[Rule]
1. Ask questions of the genre that the interlocutor wants to listen to.
2. Ask only one question at a time.
3. Output <end> at the end of a sentence when recommending tracks

Then, please ask a question.`,
				},
			}

			availableGenres, err := client.GetAvailableGenreSeeds(ctx)
			if err != nil {
				return fmt.Errorf("fail to get available genres: %v", err)
			}

			scanner := bufio.NewScanner(os.Stdin)
			for {
				resp, err := openAIClient.CreateChatCompletion(
					ctx,
					openai.ChatCompletionRequest{
						Model:        openai.GPT3Dot5Turbo,
						Messages:     messages,
						Functions:    functions.GetTemplate(availableGenres),
						FunctionCall: "auto",
					},
				)
				if err != nil {
					return fmt.Errorf("chat completion error: %v", err)
				}

				// show AI response
				functionCall := resp.Choices[0].Message.FunctionCall
				if functionCall != nil {
					if functionCall.Name == "recommend" {
						args := make(map[string]string)
						if err := json.Unmarshal([]byte(functionCall.Arguments), &args); err != nil {
							return fmt.Errorf("fail to unmarshal json: %v", err)
						}

						recommendations, err := functions.Recommend(ctx, client, args["genres"])
						if err != nil {
							return err
						}

						messages = append(messages, openai.ChatCompletionMessage{
							Name:    functionCall.Name,
							Role:    openai.ChatMessageRoleFunction,
							Content: recommendations,
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
			}

			return nil
		},
	}
}
