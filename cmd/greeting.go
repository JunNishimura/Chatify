/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	greetingUI "github.com/JunNishimura/Chatify/ui/cmd/greeting"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	port string
)

func NewGreetingCommand() *cobra.Command {
	greetingCmd := &cobra.Command{
		Use:   "greeting",
		Short: "config setting for Chatify",
		Long:  "config setting for Chatify",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			model, err := greetingUI.NewModel(port)
			if err != nil {
				return err
			}

			p := tea.NewProgram(model, tea.WithAltScreen())

			if _, err := p.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	greetingCmd.Flags().StringVarP(&port, "port", "p", "8888", "port nubmer for Spotify authorization")

	return greetingCmd
}
