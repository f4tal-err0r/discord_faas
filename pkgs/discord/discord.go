package discord

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/f4tal-err0r/discord_faas/pkgs/config"
	"github.com/f4tal-err0r/discord_faas/pkgs/db"
)

type FaasDB interface {
	GetGuild(guildID int) (*db.GuildMetadata, error)
	GetCommandsByGuild(guildID int) ([]db.Command, error)
	InsertGuildMetadata(guild db.GuildMetadata) error
}

type Client struct {
	Session           *discordgo.Session
	ContextTokenCache *sync.Map
	dbc               FaasDB
}

func NewClient(dbc FaasDB, cfg *config.Config) (*Client, error) {
	dsession, err := discordgo.New("Bot " + cfg.Discord.Token)
	if err != nil {
		log.Fatalf("ERR: %s", err)
	}

	return &Client{
		Session:           dsession,
		ContextTokenCache: &sync.Map{},
		dbc:               dbc,
	}, nil
}

func (c *Client) InitGuildData() error {
	var newGuilds []string

	botGuilds, err := c.Session.UserGuilds(100, "", "", false)
	if err != nil {
		log.Fatalf("error getting guilds: %v", err)
	}

	for _, guild := range botGuilds {
		gid, err := strconv.Atoi(guild.ID)
		if err != nil {
			return fmt.Errorf("error converting guild ID: %v", err)
		}
		_, err = c.dbc.GetGuild(gid)
		if err != nil {
			if err == sql.ErrNoRows {
				newGuilds = append(newGuilds, guild.ID)
			} else {
				return fmt.Errorf("error querying guild: %v", err)
			}
		}
	}

	for _, gid := range newGuilds {
		guild, err := c.Session.Guild(gid)
		if err != nil {
			return fmt.Errorf("error getting guild info: %v", err)
		}

		guildID, err := strconv.Atoi(guild.ID)
		if err != nil {
			return fmt.Errorf("error converting guild ID: %v", err)
		}

		guildMetadata := db.GuildMetadata{
			GuildID: guildID,
			Source:  "discord",
			Name:    guild.Name,
			Owner:   guild.OwnerID,
		}

		if err := c.dbc.InsertGuildMetadata(guildMetadata); err != nil {
			return fmt.Errorf("error inserting guild metadata: %v", err)
		}
	}

	return nil
}
