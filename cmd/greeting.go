/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/JunNishimura/Chatify/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func NewGreetingCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "greeting",
		Short: "config setting for Chatify",
		Long:  "config setting for Chatify",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := tea.NewProgram(ui.NewModel(), tea.WithAltScreen())

			if _, err := p.Run(); err != nil {
				return err
			}

			return nil
		},
	}
}
