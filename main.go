/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/JunNishimura/Chatify/cmd"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/pkg/browser"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const (
	envFileName = ".env"
	redirectURI = "http://localhost:8888/callback"
	port        = "8888"
)

var (
	auth          *spotifyauth.Authenticator
	clientChannel = make(chan *spotify.Client)
	server        *http.Server
	state         string
)

func main() {
	setEnv()

	if err := godotenv.Load(envFileName); err != nil {
		log.Fatalf("fail to load env file: %v\n", err)
	}
	auth = spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate),
	)

	server = &http.Server{Addr: fmt.Sprintf(":%s", port)}
	state = uuid.New().String()
	http.HandleFunc("/callback", completeAuth)
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	authURL := auth.AuthURL(state)
	if err := browser.OpenURL(authURL); err != nil {
		log.Fatalf("fail to open: %s\n", authURL)
	}

	client := <-clientChannel

	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatalf("fail to get current user: %v\n", err)
	}
	fmt.Printf("logged in as: %s\n", user.DisplayName)

	cmd.Execute()
}
