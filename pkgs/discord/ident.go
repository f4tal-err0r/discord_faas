package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/f4tal-err0r/discord_faas/pkgs/config"
)

func GetCurrentUser(token string) *discordgo.User {
	session, err := discordgo.New("Bearer " + token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	user, err := session.User("@me")
	if err != nil {
		log.Fatalf("Error getting current user: %v", err)
	}
	return user
}

func GetUserGuildInfo(gid string, user *discordgo.User) []string {
	cfg, err := config.New()
	if err != nil {
		log.Fatal("ERR: Unable to fetch config: %w", err)
	}
	botSession := GetSession(cfg)
	member, err := botSession.GuildMember(gid, user.ID)
	if err != nil {
		log.Fatalf("Error getting guild member: %v", err)
	}
	return member.Roles
}
