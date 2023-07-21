/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package hey

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

func NewCommand(openaiApiKey string) *cobra.Command {
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

The bot suggests recommended artists to the interlocutor.
Through a question-and-answer session, the bot asks the interlocutor what kind of music he or she likes and what kind of mood he or she is in. Based on the information obtained, we will suggest recommended artists. However, please observe the following rules.

[Rule]
1. Three questions will be asked.
2. Ask only one question at a time.
3. Recommend three artists.
4. Recommend artists based on three answers.
5. Output <end> at the end of a sentence when recommending artists

Then, please ask a question.`,
				},
			}

			scanner := bufio.NewScanner(os.Stdin)
			for {
				resp, err := openAIClient.CreateChatCompletion(
					context.Background(),
					openai.ChatCompletionRequest{
						Model:    openai.GPT3Dot5Turbo,
						Messages: messages,
					},
				)
				if err != nil {
					return fmt.Errorf("chat completion error: %v", err)
				}

				// show AI response
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

			return nil
		},
	}
}
