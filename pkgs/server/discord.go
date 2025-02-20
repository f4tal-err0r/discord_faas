package server

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/f4tal-err0r/discord_faas/pkgs/config"
)

var (
	contextToken    sync.Map
	defaultCommands = map[string]DefaultCommandData{
		"help": {
			Description: "Information around Discord FaaS",
		},
		"login": {
			Description: "Generate a login token for the command line client",
			Function:    loginCommand,
		},
	}
)

type DefaultCommandData struct {
	Description string
	Function    func(i *discordgo.Interaction) *discordgo.InteractionResponse
}

func RegisterCommands(db *sql.DB, session *discordgo.Session) error {
	var commands []discordgo.ApplicationCommand

	//init default commands
	for k, v := range defaultCommands {
		commands = append(commands, discordgo.ApplicationCommand{
			Name:        k,
			Description: v.Description,
			//TODO: Add options
		})
	}

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
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
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

	message.Data.Content = "To Login to the command line client, use the following command: `dfaas context connect --url https://" + cfg.URLDomain + " --token " + token +
		"`\nTo Download the CLI Client: https://github.com/f4tal-err0r/discord_faas/"
	return message
}

func helpCommand(i *discordgo.Interaction) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "TODO: Documentation for the help command",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}
}
