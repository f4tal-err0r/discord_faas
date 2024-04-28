package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/f4tal-err0r/discord_faas/pkgs/config"
)

func StartBot(c *config.Config) *discordgo.Session {
	dc, err := discordgo.New("Bot " + c.Discord.Token)
	if err != nil {
		log.Fatalf("ERR: %s", err)
	}
	return dc
}
