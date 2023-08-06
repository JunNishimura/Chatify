/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log"

	"github.com/JunNishimura/Chatify/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "chatify",
		Short: "chatify is a CLI tool that suggests music recommendations for you",
		Long:  "chatify is a CLI tool that suggests music recommendations for you",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	// heyCommand := cmd.NewHeyCommand(context.Background(), client, clientViper.GetString(OpenAIApiKeyName))
	greetingCommand := cmd.NewGreetingCommand()

	rootCmd.AddCommand(
		// heyCommand,
		greetingCommand,
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
