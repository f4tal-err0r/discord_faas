package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/f4tal-err0r/discord_faas/pkgs/discord"
	"github.com/f4tal-err0r/discord_faas/pkgs/security"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type Handler struct {
	Jwtsvc  *security.JWTService
	discord *discord.Client
}

func NewAuthHandler(Jwtsvc *security.JWTService, discord *discord.Client) *Handler {
	return &Handler{Jwtsvc: Jwtsvc, discord: discord}
}

func (h *Handler) AuthHandler(w http.ResponseWriter, r *http.Request) {
	token := getTokenFromHeader(r)
	if token == "" {
		http.Error(w, "Unauthorized: Token is missing", http.StatusUnauthorized)
	}
	guildID := r.Header.Get("GuildID")

	guildMember, err := discord.IdentGuildMember(token, guildID)
	if err != nil {
		http.Error(w, "Error identifying guild member: "+err.Error(), http.StatusInternalServerError)
	}

	if (guildMember.Permissions & (1 << 3)) != 0 {
		http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
	}

	jwtToken, err := h.Jwtsvc.CreateToken(security.Claims{
		UserID:  guildMember.User.ID,
		GuildID: guildID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(6 * time.Hour)),
		},
	})
	if err != nil {
		http.Error(w, "Error creating JWT: "+err.Error(), http.StatusInternalServerError)
	}

	w.Write([]byte(jwtToken))
}

func (h *Handler) AddRoute(r *mux.Router) {
	r.HandleFunc("/api/context/auth", h.AuthHandler)
}

func (h *Handler) IsSecure() bool {
	return false
}

func getTokenFromHeader(r *http.Request) string {
	token := r.Header.Get("Authorization")
	if !strings.HasPrefix(token, "Bearer ") {
		return ""
	}

	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		return ""
	}

	return token
}
