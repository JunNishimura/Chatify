package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/browser"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const (
	port = "8888"
)

var (
	auth   *spotifyauth.Authenticator
	state  string
	server *http.Server
)

func completeAuth(w http.ResponseWriter, r *http.Request) {
	token, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "fail to get token", http.StatusForbidden)
		log.Fatalf("fail to get token: %v\n", err)
	}
	if getState := r.FormValue("state"); getState != state {
		http.NotFound(w, r)
		log.Fatalf("state mismatch: got = %s, expected = %s\n", getState, state)
	}

	// set token config
	if err := saveToken(token); err != nil {
		log.Fatal(err)
	}

	client := spotify.New(auth.Client(r.Context(), token))
	clientChannel <- client

	go func() {
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("shutdown server")
			log.Fatalf("fail to shutdown server: %v\n", err)
		}
	}()
}

func authorize() {
	auth = spotifyauth.New(
		spotifyauth.WithClientID(clientViper.GetString(SpotifyIDKeyName)),
		spotifyauth.WithClientSecret(clientViper.GetString(SpotifySecretKeyName)),
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
}
