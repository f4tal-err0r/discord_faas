package deploy

import (
	"os"

	"github.com/f4tal-err0r/discord_faas/pkgs/discord"
	"github.com/f4tal-err0r/discord_faas/pkgs/security"
)

type Handler struct {
	dbot   *discord.Client
	jwtsvc *security.JWTService
}

func init() {
	if err := os.Mkdir("/app/data/artifacts", 0755); !os.IsExist(err) && err != nil {
		panic(err)
	}
}

func NewHandler(dbot *discord.Client) *Handler {
	return &Handler{
		dbot: dbot,
	}
}
