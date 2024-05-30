package client

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/f4tal-err0r/discord_faas/pkgs/shared"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"
)

var (
	oauthCfg     *oauth2.Config
	state        string
	codeVerifier string
	tokenChan    chan *oauth2.Token
	cachefp      string
)

func init() {
	cachefp = FetchCache("token")
	oauthCfg = &oauth2.Config{
		// ClientID:    cfg.Oauth.ClientID,
		RedirectURL: "http://localhost:8080/callback",
		Scopes:      []string{"guilds", "guilds.members.read", "identify"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discord.com/oauth2/authorize",
			TokenURL: "https://discord.com/api/oauth2/token",
		},
	}
	state = "VIgSXcWvBgLtHt4T9MVPg0jr" // you should generate this randomly
	codeVerifier = generateCodeVerifier()
	tokenChan = make(chan *oauth2.Token)
}

func generateCodeVerifier() string {
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

func StartAuth() (*oauth2.Token, error) {
	http.HandleFunc("/callback", handleCallback)

	codeChallenge, err := generateCodeChallenge(codeVerifier)
	if err != nil {
		log.Fatalf("Unable to generate code challenge: %s", err)
	}

	url := oauthCfg.AuthCodeURL(state, oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"))

	fmt.Printf("Open Browser to auth: %s", url)
	browser.OpenURL(url)

	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	token := <-tokenChan

	if err := saveToken(token); err != nil {
		log.Printf("\nWARN: Unable to cache Oauth2 token: %v", err)
	}
	return token, nil
}

func GetToken() (string, error) {
	var token oauth2.Token

	f, err := os.Open(cachefp)
	if err != nil {
		token, err := StartAuth()
		if err != nil {
			log.Fatalf("\nERROR: Failed to auth to Discord: %v", err)
		}
		return token.AccessToken, nil
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&token)
	if err != nil {
		return "", err
	}

	tokenSource := oauthCfg.TokenSource(context.Background(), &token)
	refreshToken, err := tokenSource.Token()
	if err != nil {
		log.Println(err)
		newToken, err := StartAuth()
		if err != nil {
			log.Fatalf("\nERROR: Failed to auth to Discord: %v", err)
		}
		if err := saveToken(newToken); err != nil {
			log.Printf("\nWARN: Unable to cache Oauth2 token: %v", err)
		}
		return newToken.AccessToken, nil
	}

	if err := saveToken(refreshToken); err != nil {
		log.Printf("\nWARN: Unable to cache Oauth2 token: %v", err)
	}

	user := shared.GetCurrentUser(refreshToken.AccessToken)

	log.Printf("User: %+v", user)
	log.Printf("ID: %+v", user.ID)
	return refreshToken.AccessToken, nil
}

func saveToken(token *oauth2.Token) error {
	f, err := os.Create(cachefp)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(token)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("state") != state {
		log.Fatal("State is not valid")
		return
	}

	code := r.URL.Query().Get("code")
	token, err := oauthCfg.Exchange(context.Background(), code, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	if err != nil {
		log.Fatal("Failed to exchange token: " + err.Error())
		return
	}

	// Send the token to the channel
	tokenChan <- token
}
