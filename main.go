/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package chatify

import (
	"context"
	"fmt"
	"log"

	"github.com/JunNishimura/Chatify/cmd"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const (
	redirectURI = "http://localhost:8888/callback"
)

var (
	clientChannel = make(chan *spotify.Client, 1)
)

func main() {
	err := initEnv()
	if err != ErrConfigNotFound && err != nil {
		log.Fatal(err)
	}

	// check if token viper is set
	if !isClientInfoSet() {
		if err := askClientInfo(); err != nil {
			log.Fatal(err)
		}

		authorize()
	} else {
		ctx := context.Background()
		token := getToken()
		httpClient := spotifyauth.New().Client(ctx, token)
		clientChannel <- spotify.New(httpClient)
	}

	client := <-clientChannel

	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatalf("fail to get current user: %v\n", err)
	}
	fmt.Printf("logged in as: %s\n", user.DisplayName)

	cmd.Execute()
}
