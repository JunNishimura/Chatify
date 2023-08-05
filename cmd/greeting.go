/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewGreetingCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "greeting",
		Short: "config setting for Chatify",
		Long:  "config setting for Chatify",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("greeting called")
		},
	}
}
