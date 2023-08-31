/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log"

	"github.com/JunNishimura/Chatify/cmd"
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "chatify <command> [flags]",
		Short: "Chatify",
		Long:  "chatify is a TUI tool that suggests music recommendations for you",
		Example: heredoc.Doc(`
			$ chatify greeting
			$ chatify hey
		`),
		Run: func(cmd *cobra.Command, args []string) {},
	}

	greetingCommand := cmd.NewGreetingCommand()
	heyCommand := cmd.NewHeyCommand()

	rootCmd.AddCommand(
		greetingCommand,
		heyCommand,
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
