/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// heyCmd represents the hey command
var heyCmd = &cobra.Command{
	Use:   "hey",
	Short: "start conversation with chatify",
	Long:  "start conversation with chatify",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(heyCmd)
}
