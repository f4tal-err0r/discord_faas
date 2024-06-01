package server

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/f4tal-err0r/discord_faas/pkgs/config"
)

var (
	defaultCommands = []discordgo.ApplicationCommand{
		{
			Name:        "delete",
			Description: "Delete a command, this is only available to the bot owner",
		},
		{
			Name:        "delegate",
			Description: "Add a command, this is only available to the bot owner",
		},
	}
)

func RegisterCommands(db *sql.DB, session *discordgo.Session) error {
	commands := defaultCommands
	guilds, err := session.UserGuilds(100, "", "", false)
	if err != nil {
		return fmt.Errorf("Error getting guilds: %v", err)
	}

	for _, guild := range guilds {
		cmdb, err := GetCmdsDb(db, StrToInt(guild.ID))
		if err != nil {
			return fmt.Errorf("Error getting commands: %v", err)
		}

		for _, cmd := range cmdb {
			commands = append(commands, discordgo.ApplicationCommand{
				Name:        cmd.Command,
				Description: cmd.Desc,
			})
		}
		for _, cmd := range commands {
			_, err := session.ApplicationCommandCreate(session.State.User.ID, guilds[0].ID, &cmd)
			if err != nil {
				return fmt.Errorf("Error creating command on guild %v: %v", guild.ID, err)
			}
		}

	}

	return nil
}

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
