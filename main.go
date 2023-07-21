/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/JunNishimura/Chatify/internal/hey"
	"github.com/spf13/cobra"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const (
	redirectURI = "http://localhost:8888/callback"
)

var (
	auth          *spotifyauth.Authenticator
	clientChannel = make(chan *spotify.Client, 1)
)

func main() {
	err := initEnv()
	if err != ErrConfigNotFound && err != nil {
		log.Fatal(err)
	}

	auth = spotifyauth.New(
		spotifyauth.WithClientID(clientViper.GetString(SpotifyIDKeyName)),
		spotifyauth.WithClientSecret(clientViper.GetString(SpotifySecretKeyName)),
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate),
	)

	// check if token viper is set
	if !isClientInfoSet() {
		if err := askClientInfo(); err != nil {
			log.Fatal(err)
		}

		authorize()
	} else {
		ctx := context.Background()
		token := getToken()

		newToken, err := auth.RefreshToken(ctx, token)
		if err != nil {
			log.Fatalf("fail to get a new access token: %v", err)
		}

		// update an access token if it has expired
		if token.AccessToken != newToken.AccessToken {
			if err := saveToken(newToken); err != nil {
				log.Fatalf("fail to save token: %v", err)
			}
		}

		spotifyClient := spotify.New(auth.Client(ctx, newToken))

		clientChannel <- spotifyClient
	}

	client := <-clientChannel

	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatalf("fail to get current user: %v\n", err)
	}
	fmt.Printf("logged in as: %s\n", user.DisplayName)

	rootCmd := &cobra.Command{
		Use:   "chatify",
		Short: "chatify is a CLI tool that suggests music recommendations for you",
		Long:  "chatify is a CLI tool that suggests music recommendations for you",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	heyCommand := hey.NewCommand(clientViper.GetString(OpenAIApiKeyName))

	rootCmd.AddCommand(heyCommand)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
