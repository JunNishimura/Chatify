package main

import (
	"context"
	"log"
	"net/http"

	"github.com/zmb3/spotify/v2"
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

	client := spotify.New(auth.Client(r.Context(), token))
	clientChannel <- client

	go func() {
		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatalf("fail to shutdown server: %v\n", err)
		}
	}()
}
