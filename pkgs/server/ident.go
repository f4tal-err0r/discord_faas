package server

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/f4tal-err0r/discord_faas/pkgs/cache"
	"github.com/f4tal-err0r/discord_faas/pkgs/config"
)

var GuildCache = cache.New()
var UserGuildsCache = cache.New()

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

func GetUserGuildInfo(gid string, user *discordgo.User) *discordgo.Member {
	cfg, err := config.New()
	if err != nil {
		log.Fatal("ERR: Unable to fetch config: %w", err)
	}
	if v, ok := GuildCache.Get(gid + user.ID); !ok {
		return v.(*discordgo.Member)
	}
	botSession := GetSession(cfg)
	member, err := botSession.GuildMember(gid, user.ID)
	if err != nil {
		log.Fatalf("Error getting guild member: %v", err)
	}
	GuildCache.Set((gid + user.ID), member, (4 * time.Hour))
	return member
}
