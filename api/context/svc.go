package context

import "github.com/f4tal-err0r/discord_faas/pkgs/discord"

type Handler struct {
	bot *discord.Client
}

func NewHandler(bot *discord.Client) *Handler {
	return &Handler{
		bot: bot,
	}
}
