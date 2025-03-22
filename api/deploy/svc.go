package deploy

import (
	"github.com/f4tal-err0r/discord_faas/pkgs/discord"
	"github.com/f4tal-err0r/discord_faas/pkgs/security"
)

type Handler struct {
	dbot   *discord.Client
	jwtsvc *security.JWTService
}

func NewHandler(dbot *discord.Client) *Handler {
	return &Handler{
		dbot: dbot,
	}
}
