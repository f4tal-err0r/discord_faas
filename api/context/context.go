package context

import (
	"log"
	"net/http"
	"strings"

	"github.com/f4tal-err0r/discord_faas/pkgs/config"
	pb "github.com/f4tal-err0r/discord_faas/proto"
	"github.com/gorilla/mux"
	"google.golang.org/protobuf/proto"
)

func (h *Handler) Handler(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.New()
	if err != nil {
		log.Println("Error getting config:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token := r.Header.Get("Token")

	guildid, ok := h.bot.ContextTokenCache.Load(token)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	} else {
		h.bot.ContextTokenCache.Delete(r.URL.Query().Get("token"))
	}

	guild, err := h.bot.Session.Guild(guildid.(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if strings.Contains(err.Error(), "404") {
		http.Error(w, "Guild not found", http.StatusNotFound)
		return
	} else if strings.Contains(err.Error(), "429") {
		http.Error(w, "Rate Limit for SelfGuildInfo exceeded", http.StatusTooManyRequests)
		return
	}

	ctxPb, err := proto.Marshal(&pb.ContextResp{
		ClientID:  cfg.Discord.ClientID,
		GuildID:   guild.ID,
		GuildName: guild.Name,
	})
	if err != nil {
		log.Println("Error marshalling ctx: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(ctxPb)
}

func (h *Handler) AddRoute(r *mux.Router) {
	r.HandleFunc("/api/context", h.Handler)
}

func (h *Handler) IsSecure() bool {
	return false
}
