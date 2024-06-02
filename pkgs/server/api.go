package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/f4tal-err0r/discord_faas/pkgs/config"
)

// TODO: Implement JWT authentication

func NewServer() *http.Server {
	mux := Router()
	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}

func Router() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/api/context", ContextHandler)
	return router
}

func APIPatchAuth(w http.ResponseWriter, r *http.Request) bool {
	oauth, guild := r.Header.Get("X-Discord-Oauth"), r.Header.Get("X-Discord-GuildId")
	if oauth == "" || guild == "" {
		return false
	}

	//Check User exists in guild
	_, err := FetchGuildMember(oauth, guild)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		w.(http.Flusher).Flush()
		return false
	}

	//TODO: This needs to also check GetRolesByGuildID to see if the user has the correct role

	return true
}

func ContextHandler(w http.ResponseWriter, r *http.Request) {
	if !APIPatchAuth(w, r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	cfg, err := config.New()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	guildID := r.Header.Get("X-Discord-GuildId")
	guild, err := GetGuildInfo(GetSession(cfg), guildID)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			http.Error(w, "Guild not found", http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "429") {
			http.Error(w, "Rate Limit for SelfGuildInfo exceeded", http.StatusTooManyRequests)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := struct {
		ClientID  string `json:"client_id"`
		GuildID   string `json:"guild_id"`
		GuildName string `json:"guild_name"`
	}{
		ClientID:  cfg.Discord.ClientID,
		GuildID:   guildID,
		GuildName: guild.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
