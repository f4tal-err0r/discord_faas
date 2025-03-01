package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/f4tal-err0r/discord_faas/pkgs/api"
	"github.com/f4tal-err0r/discord_faas/pkgs/discord"
	"github.com/f4tal-err0r/discord_faas/pkgs/security"
	jwt "github.com/golang-jwt/jwt/v5"
)

func TestNewRouter(t *testing.T) {
	bot := &discord.Client{}
	jwtService := &security.JWTService{}
	router := api.NewRouter(bot, jwtService)

	if router.Router == nil {
		t.Error("router not created successfully")
	}
}

func TestDeployHandler(t *testing.T) {
	bot := &discord.Client{}
	jwtService := &security.JWTService{}
	router := api.NewRouter(bot, jwtService)

	req, err := http.NewRequest("GET", "/api/deploy", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	router.Router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestContextHandler(t *testing.T) {
	bot := &discord.Client{}
	jwtService := &security.JWTService{}
	router := api.NewRouter(bot, jwtService)

	req, err := http.NewRequest("GET", "/api/context", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	router.Router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestAuthHandler(t *testing.T) {
	bot := &discord.Client{}
	jwtService := &security.JWTService{}
	router := api.NewRouter(bot, jwtService)

	req, err := http.NewRequest("GET", "/api/context/auth", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	router.Router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestDecodeHandlerInvalidToken(t *testing.T) {
	bot := &discord.Client{}
	jwtService := &security.JWTService{}
	router := api.NewRouter(bot, jwtService)

	req, err := http.NewRequest("GET", "/api/context/decode", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	router.Router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestDecodeHandlerValidToken(t *testing.T) {
	bot := &discord.Client{}
	jwtService, _ := security.NewJWT() // Initialize JWTService with keys
	router := api.NewRouter(bot, jwtService)

	claims := security.Claims{
		UserID:  "user123",
		GuildID: "guild123",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(6 * time.Hour)),
		},
	}
	token, err := jwtService.CreateToken(claims)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/api/context/decode", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.Router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var respClaims security.Claims
	err = json.Unmarshal(w.Body.Bytes(), &respClaims)
	if err != nil {
		t.Fatal(err)
	}

	if respClaims.UserID != claims.UserID || respClaims.GuildID != claims.GuildID {
		t.Errorf("expected claims %v, got %v", claims, respClaims)
	}
}
