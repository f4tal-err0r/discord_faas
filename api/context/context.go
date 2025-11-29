package context

import (
	"net/http"
	"strings"

	v1 "github.com/f4tal-err0r/discord_faas/api/v1"
	"github.com/gorilla/mux"
	"google.golang.org/protobuf/proto"
)

func (h *Handler) Handler(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Token")

	guildid, ok := h.bot.ContextTokenCache.Load(token)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	} else {
		h.bot.ContextTokenCache.Delete(r.URL.Query().Get("token"))
	}

	guild, err := h.bot.Session.Guild(guildid.(string))
	//check for 404 and 429
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			http.Error(w, "Guild not found", http.StatusNotFound)
			return
		} else if strings.Contains(err.Error(), "429") {
			http.Error(w, "Rate limited, try again later", http.StatusTooManyRequests)
			return
		} else {
			http.Error(w, "Error fetching guild info", http.StatusInternalServerError)
			return
		}
	}

	ctxresp := v1.ContextResp{
		ClientId:  h.cfg.Discord.ClientID,
		GuildId:   guild.ID,
		GuildName: guild.Name,
	}
	resp, err := proto.Marshal(&ctxresp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

func (h *Handler) AddRoute(r *mux.Router) {
	r.HandleFunc("/api/context", h.Handler)
}

func (h *Handler) IsSecure() bool {
	return false
}
