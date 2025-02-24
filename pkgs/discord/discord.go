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

var (
	defaultCommands = map[string]defaultCommandData{
		"help": {
			Description: "Information around Discord FaaS",
			Function:    helpCommand,
		},
		"login": {
			Description: "Generate a login token for the command line client",
			Function:    loginCommand,
		},
	}
)

type helpFunc func(*discordgo.Interaction, *Client) *discordgo.InteractionResponse

type defaultCommandData struct {
	Description string
	Function    helpFunc
}

type Client struct {
	Session           *discordgo.Session
	ContextTokenCache *sync.Map
}

func NewClient(cfg *config.Config) (*Client, error) {
	dsession, err := discordgo.New("Bot " + cfg.Discord.Token)
	if err != nil {
		log.Fatalf("ERR: %s", err)
	}
	return &Client{
		Session:           dsession,
		ContextTokenCache: &sync.Map{},
	}, nil
}

func (c *Client) RegisterCommands(db *db.DBHandler) error {
	c.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			if _, ok := defaultCommands[i.ApplicationCommandData().Name]; ok {
				s.InteractionRespond(i.Interaction, defaultCommands[i.ApplicationCommandData().Name].Function(i.Interaction, c))
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Pong!",
				},
			})
		}
	})
	var commands []discordgo.ApplicationCommand

	//init default commands
	for k, v := range defaultCommands {
		commands = append(commands, discordgo.ApplicationCommand{
			Name:        k,
			Description: v.Description,
			//TODO: Add options
		})
	}

	guilds, err := c.Session.UserGuilds(100, "", "", false)
	if err != nil {
		return fmt.Errorf("error getting guilds: %v", err)
	}

	for _, guild := range guilds {
		cmdb, err := db.GetCommandsByGuild(strToInt(guild.ID))
		if err != nil {
			return fmt.Errorf("error getting commands: %v", err)
		}

		for _, cmd := range cmdb {
			commands = append(commands, discordgo.ApplicationCommand{
				Name:        cmd.Command,
				Description: cmd.Description,
				//TODO: Add options
			})
		}
		for _, cmd := range commands {
			_, err := c.Session.ApplicationCommandCreate(c.Session.State.User.ID, guild.ID, &cmd)
			if err != nil {
				return fmt.Errorf("error creating command on guild %v: %v", guild.ID, err)
			}
		}

	}

	return nil
}

func loginCommand(i *discordgo.Interaction, c *Client) *discordgo.InteractionResponse {
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

	token, err := generateRandomHash()
	if err != nil {
		message.Data.Content = "Error generating token"
		return message
	}
	c.ContextTokenCache.Store(token, i.GuildID)

	message.Data.Content = "To Login to the command line client, use the following command: `dfaas context connect --url https://" + cfg.URLDomain + " --token " + token +
		"`\nTo Download the CLI Client: https://github.com/f4tal-err0r/discord_faas/"
	return message
}

func helpCommand(_ *discordgo.Interaction, _ *Client) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "TODO: Documentation for the help command",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}
}

func (c *Client) InitGuildData(dbc *db.DBHandler) error {
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
		_, err = dbc.GetGuild(gid)
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

		if err := dbc.InsertGuildMetadata(guildMetadata); err != nil {
			return fmt.Errorf("error inserting guild metadata: %v", err)
		}
	}

	return nil
}
