package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/JunNishimura/Chatify/config"
	"github.com/JunNishimura/spotify/v2"
	spotifyauth "github.com/JunNishimura/spotify/v2/auth"
	"github.com/google/uuid"
	"github.com/pkg/browser"
)

const (
	authPort    = "8888"
	redirectURI = "http://localhost:8888/callback"
)

type Client struct {
	SpotifyChannel chan *spotify.Client
	cfg            *config.Config
	auth           *spotifyauth.Authenticator
	server         *http.Server
	state          string
}

func New(cfg *config.Config) *Client {
	auth := spotifyauth.New(
		spotifyauth.WithClientID(cfg.GetClientValue(config.SpotifyIDKey)),
		spotifyauth.WithClientSecret(cfg.GetClientValue(config.SpotifySecretKey)),
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate),
	)

	return &Client{
		SpotifyChannel: make(chan *spotify.Client, 1),
		cfg:            cfg,
		auth:           auth,
		server:         &http.Server{Addr: fmt.Sprintf(":%s", authPort)},
		state:          uuid.New().String(),
	}
}

func (a *Client) Authorize() {
	http.HandleFunc("/callback", a.completeAuth)
	go func() {
		if err := a.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	authURL := a.auth.AuthURL(a.state)
	if err := browser.OpenURL(authURL); err != nil {
		log.Fatalf("fail to open %s: %v", authURL, err)
	}
}

func (a *Client) completeAuth(w http.ResponseWriter, r *http.Request) {
	token, err := a.auth.Token(r.Context(), a.state, r)
	if err != nil {
		http.Error(w, "fail to get token", http.StatusForbidden)
		log.Fatalf("fail to get token: %v\n", err)
	}
	if getState := r.FormValue("state"); getState != a.state {
		http.NotFound(w, r)
		log.Fatalf("state mismatch: got = %s, expected = %s\n", getState, a.state)
	}

	if err := a.cfg.SetToken(token); err != nil {
		log.Fatalf("fail to set token: %v", err)
	}

	a.SpotifyChannel <- spotify.New(a.auth.Client(r.Context(), token))

	go func() {
		if err := a.server.Shutdown(context.Background()); err != nil {
			log.Fatalf("fail to shutdown server: %v\n", err)
		}
	}()
}
