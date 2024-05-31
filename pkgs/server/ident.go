package server

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/f4tal-err0r/discord_faas/pkgs/cache"
	"github.com/f4tal-err0r/discord_faas/pkgs/config"
)

var UserGuildsCache = cache.New()
var GuildCache = cache.New()

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
	if v, ok := UserGuildsCache.Get(gid + user.ID); !ok {
		return v.(*discordgo.Member)
	}
	botSession := GetSession(cfg)
	member, err := botSession.GuildMember(gid, user.ID)
	if err != nil {
		log.Fatalf("Error getting guild member: %v", err)
	}
	UserGuildsCache.Set((gid + user.ID), member, (4 * time.Hour))
	return member
}

func GetGuildInfo(gid string) *discordgo.Guild {
	cfg, err := config.New()
	if err != nil {
		log.Fatal("ERR: Unable to fetch config: %w", err)
	}
	if v, _ := GuildCache.Get(gid); v != nil {
		return v.(*discordgo.Guild)
	}
	botSession := GetSession(cfg)
	guild, err := botSession.Guild(gid)
	if err != nil {
		log.Fatalf("Error getting guild: %v", err)
	}
	GuildCache.Set(gid, guild, (4 * time.Hour))
	return guild
}

func GetDefaultChannel(session *discordgo.Session, guildID string) (*discordgo.Channel, error) {
	channels, err := session.GuildChannels(guildID)
	if err != nil {
		return nil, err
	}

	var defaultChannel *discordgo.Channel
	for _, channel := range channels {
		if channel.Type == discordgo.ChannelTypeGuildText {
			if defaultChannel == nil || channel.Position < defaultChannel.Position {
				defaultChannel = channel
			}
		}
	}

	if defaultChannel == nil {
		return nil, fmt.Errorf("no default channel found")
	}
	return defaultChannel, nil
}
