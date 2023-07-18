/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package hey

import (
	"context"
	"fmt"

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
			resp, err := openAIClient.CreateChatCompletion(
				context.Background(),
				openai.ChatCompletionRequest{
					Model: openai.GPT3Dot5Turbo,
					Messages: []openai.ChatCompletionMessage{
						{
							Role:    openai.ChatMessageRoleUser,
							Content: "Hello",
						},
					},
				},
			)
			if err != nil {
				return fmt.Errorf("chat completion error: %v", err)
			}

			// show response
			fmt.Println(resp.Choices[0].Message.Content)

			return nil
		},
	}
}
