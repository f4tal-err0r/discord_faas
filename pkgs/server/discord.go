package server

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/f4tal-err0r/discord_faas/pkgs/config"
)

var (
	systemCommands = []discordgo.ApplicationCommand{
		{
			Name:        "allowuser",
			Description: "Allows a guild , this is only available to the bot owner",
		},
		{
			Name:        "allowguild",
			Description: "Add a command, this is only available to the bot owner",
		},
	}
	defaultCommands = []discordgo.ApplicationCommand{
		{
			Name:        "help",
			Description: "Get help",
		},
	}
)

func RegisterCommands(db *sql.DB, session *discordgo.Session) error {
	commands := defaultCommands
	guilds, err := session.UserGuilds(100, "", "", false)
	if err != nil {
		return fmt.Errorf("error getting guilds: %v", err)
	}

	for _, guild := range guilds {
		cmdb, err := GetCmdsDb(db, StrToInt(guild.ID))
		if err != nil {
			return fmt.Errorf("error getting commands: %v", err)
		}

		for _, cmd := range cmdb {
			commands = append(commands, discordgo.ApplicationCommand{
				Name:        cmd.Command,
				Description: cmd.Desc,
				//TODO: Add options
			})
		}
		for _, cmd := range commands {
			_, err := session.ApplicationCommandCreate(session.State.User.ID, guild.ID, &cmd)
			_, err := session.ApplicationCommandCreate(session.State.User.ID, guild.ID, &cmd)
			if err != nil {
				return fmt.Errorf("error creating command on guild %v: %v", guild.ID, err)
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
	db, err := NewDB(cfg)
	if err != nil {
		log.Fatalf("ERR: Unable to create db: %v", err)
	}
	if err := dc.Open(); err != nil {
		log.Fatal(fmt.Errorf("websocket error: %v", err))
	}
	if err := RegisterCommands(db, dc); err != nil {
		log.Fatalf("ERR: %s", err)
	}
	return dc
}

func loginCommand(i *discordgo.Interaction) *discordgo.InteractionResponse {

	message := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{},
	}

	cfg, err := config.New()
	if err != nil {
		message.Data.Content = "Error getting config"
		return message
	}

	token, err := GenerateRandomHash()
	if err != nil {
		message.Data.Content = "Error generating token"
		return message
	}
	contextToken.Store(token, i.GuildID)

	message.Data.Content = "To Login to the command line client, use the following command: `dfaas context connect --url https://" + cfg.URLDomain + " --token " + token + "`\nTo Download the CLI Client: https://github.com/f4tal-err0r/discord_faas/"
	return message
}
