package discord

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/f4tal-err0r/discord_faas/pkgs/config"
)

func GetSession(cfg *config.Config) *discordgo.Session {
	dc, err := discordgo.New("Bot " + cfg.Discord.Token)
	if err != nil {
		log.Fatalf("ERR: %s", err)
	}
	if err := dc.Open(); err != nil {
		log.Fatal(err)
	}
	return dc
}

func StartDiscordBot() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal("ERR: Unable to fetch config: %w", err)
	}

	dc := GetSession(cfg)
	log.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	dc.Close()
}
