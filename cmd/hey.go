/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"

	heyUI "github.com/JunNishimura/Chatify/ui/cmd/hey"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	recommendNum int
)

func NewHeyCommand() *cobra.Command {
	heyCmd := &cobra.Command{
		Use:   "hey",
		Short: "start conversation with chatify",
		Long:  "start conversation with chatify",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			model, err := heyUI.NewModel(context.Background(), recommendNum)
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

	heyCmd.Flags().IntVarP(&recommendNum, "num", "n", 25, "the number of recommendations(maximum is 100)")

	return heyCmd
}
