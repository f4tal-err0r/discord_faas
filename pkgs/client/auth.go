package client

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/browser"
	"golang.org/x/oauth2"
)

type DiscordUserAuth struct {
	Token    *oauth2.Token
	Config   *oauth2.Config
	Filepath string
	Browser  bool
}

func NewUserAuth(opts ...func(*DiscordUserAuth)) *DiscordUserAuth {
	var userauth DiscordUserAuth
	context, err := GetCurrentContext()
	if err != nil {
		log.Fatal(err)
	}

	for _, opt := range opts {
		opt(&userauth)
	}
	oauthCfg := &oauth2.Config{
		ClientID:    context.ClientID,
		RedirectURL: "http://localhost:8085/callback",
		Scopes:      []string{"guilds", "guilds.members.read", "identify"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discord.com/oauth2/authorize",
			TokenURL: "https://discord.com/api/oauth2/token",
		},
	}

	userauth = DiscordUserAuth{
		Token:    nil,
		Config:   oauthCfg,
		Filepath: FetchCacheDir("auth"),
	}

	return &userauth
}

func WithToken(token *oauth2.Token) func(*DiscordUserAuth) {
	return func(d *DiscordUserAuth) {
		d.Token = token
	}
}

func generateRand() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func generateCodeChallenge(codeVerifier string) (string, error) {
	h := sha256.New()
	_, err := h.Write([]byte(codeVerifier))
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil)), nil
}

func (d *DiscordUserAuth) StartAuth() (*oauth2.Token, error) {
	tokenChan := make(chan *oauth2.Token)
	state := generateRand()
	codeVerifier := generateRand()

	codeChallenge, err := generateCodeChallenge(codeVerifier)
	if err != nil {
		log.Fatalf("Unable to generate code challenge: %s", err)
	}

	url := d.Config.AuthCodeURL(state, oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"))

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state {
			log.Fatal("State is not valid")
			return
		}

		code := r.URL.Query().Get("code")
		token, err := d.Config.Exchange(context.Background(), code, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
		if err != nil {
			log.Fatal("Failed to exchange token: " + err.Error())
			return
		}

		// Send token to channel
		tokenChan <- token
	})

	go func() {
		log.Fatal(http.ListenAndServe(":8085", nil))
	}()

	// Open the browser to start the auth flow
	if err := browser.OpenURL(url); err != nil {
		return nil, fmt.Errorf("open browser: %w", err)
	}

	// Wait for the auth flow to complete and return the token
	token := <-tokenChan

	return token, nil
}
