package discord

import (
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"time"

	"github.com/f4tal-err0r/discord_faas/config"
	"github.com/f4tal-err0r/discord_faas/pkgs/config"
	"golang.org/x/oauth2"
)

type Credentials struct {
}

var Endpoint = oauth2.Endpoint{
	AuthURL:  "https://discord.com/oauth2/authorize",
	TokenURL: "https://discord.com/api/oauth2/token",
}

func init() {
	oauthCfg := &oauth2.Config{
		ClientID:     &config.Oauth.ClientID,
		ClientSecret: &config.Oauth.ClientSecret,
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"guilds", "guilds.members.read"},
	}
}

func generateCodeVerifier() string {
	rand.Seed(time.Now().UnixNano())
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
