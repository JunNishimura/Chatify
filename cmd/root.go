/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "chatify",
	Short: "chatify is a CLI tool that suggests music recommendations for you",
	Long:  "chatify is a CLI tool that suggests music recommendations for you",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("root called")
		// get client
		// client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
		// resp, err := client.CreateChatCompletion(
		// 	context.Background(),
		// 	openai.ChatCompletionRequest{
		// 		Model: openai.GPT3Dot5Turbo,
		// 		Messages: []openai.ChatCompletionMessage{
		// 			{
		// 				Role:    openai.ChatMessageRoleUser,
		// 				Content: "Hello",
		// 			},
		// 		},
		// 	},
		// )
		// if err != nil {
		// 	return fmt.Errorf("chat completion error: %v", err)
		// }

		// // show response
		// fmt.Println(resp.Choices[0].Message.Content)

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
