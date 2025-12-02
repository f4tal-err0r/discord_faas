package functions

import (
	"context"

	"github.com/f4tal-err0r/discord_faas/pkgs/config"
	"github.com/gorilla/mux"
)

type Handler struct {
	cfg config.Config
}

type Platform interface {
	Start(ctx *context.Context) error
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		cfg: *cfg,
	}
}

func (h *Handler) AddRoute(r *mux.Router) {
	r.HandleFunc("/api/functions/{guildid}", h.GetFuncsHandler).Methods("GET")
	r.HandleFunc("/api/functions/{guildid}", h.DeployFuncHandler).Methods("POST")
	r.HandleFunc("/api/functions/{guildid}/{hash}", h.GetFuncsHandler).Methods("GET")
}

func (h *Handler) IsSecure() bool {
	return true
}
