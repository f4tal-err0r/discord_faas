package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/f4tal-err0r/discord_faas/pkgs/cache"
)

type OAuth2Response struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type GuildMember struct {
	User     OAuth2Response `json:"user"`
	Nickname string         `json:"nick"`
	Roles    []string       `json:"roles"`
	GuildID  string         `json:"id"`
}

var GuildCache = cache.New()
var GuildMemberCache = cache.New()

func FetchGuildMember(token, gid string) (*GuildMember, error) {
	endpoint := fmt.Sprintf("https://discord.com/api/users/@me/guilds/%s/member", gid)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch guild member: %s", resp.Status)
	}

	var member GuildMember
	if err := json.NewDecoder(resp.Body).Decode(&member); err != nil {
		return nil, err
	}

	return &member, nil
}

func GetGuildInfo(session *discordgo.Session, gid string) (*discordgo.Guild, error) {
	if v, _ := GuildCache.Get(gid); v != nil {
		return v.(*discordgo.Guild), nil
	}
	guild, err := session.Guild(gid)
	if err != nil {
		return nil, err
	}
	GuildCache.Set(gid, guild, (4 * time.Hour))
	return guild, nil
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
