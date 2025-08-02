package context

import (
	"github.com/f4tal-err0r/discord_faas/internal/discord"
	"github.com/f4tal-err0r/discord_faas/pkgs/config"
)

type Handler struct {
	bot *discord.Client
	cfg *config.Config
}

func NewHandler(bot *discord.Client, cfg *config.Config) *Handler {
	return &Handler{
		bot: bot,
		cfg: cfg,
	}
}
