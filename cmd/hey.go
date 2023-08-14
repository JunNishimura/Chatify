/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	heyUI "github.com/JunNishimura/Chatify/ui/hey"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type stringInfoResponse struct {
	Value    string `json:"value"`
	Subvalue string `json:"subvalue,omitempty"`
}

type numberInfoResponse struct {
	Value    float64 `json:"value"`
	Subvalue string  `json:"subvalue,omitempty"`
}

func NewHeyCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "hey",
		Short: "start conversation with chatify",
		Long:  "start conversation with chatify",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := tea.NewProgram(heyUI.NewModel(), tea.WithAltScreen())

			if _, err := p.Run(); err != nil {
				return err
			}

			return nil
		},
	}
}
