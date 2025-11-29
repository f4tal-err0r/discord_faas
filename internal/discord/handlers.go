package discord

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	pb "github.com/f4tal-err0r/discord_faas/proto"
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
	notif = map[string]chan *pb.DiscordResp{}
)

type handlerFunc func(*discordgo.Interaction, *Client) *discordgo.InteractionResponse

type defaultCommandData struct {
	Description string
	Function    handlerFunc
}

func (c *Client) StartBotHandler() error {
	if err := c.Session.Open(); err != nil {
		return fmt.Errorf("error opening discord session: %v", err)
	}

	log.Print("Bot Started...")

	if err := c.RegisterCommands(); err != nil {
		return fmt.Errorf("error registering commands: %v", err)
	}

	c.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		//if bot is not member of the guild, ignore
		if _, err := s.State.Guild(i.GuildID); err != nil {
			log.Printf("Ignoring interaction from guild %s: bot not a member", i.GuildID)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Bot is not a registered member of this guild.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}
		if h, ok := defaultCommands[i.ApplicationCommandData().Name]; ok {
			s.InteractionRespond(i.Interaction, h.Function(i.Interaction, c))
		}
		s.InteractionRespond(i.Interaction, cmdRouter(i.Interaction, c))
	})
	return nil
}

func (c *Client) RegisterCommands() error {
	//list guilds
	guilds, err := c.Session.UserGuilds(100, "", "", false)
	if err != nil {
		return fmt.Errorf("error getting user guilds: %v", err)
	}

	for _, guild := range guilds {
		var commands []discordgo.ApplicationCommand

		//init default commands
		for k, v := range defaultCommands {
			commands = append(commands, discordgo.ApplicationCommand{
				Name:        k,
				Description: v.Description,
			})
		}
		cmdb, err := c.dbc.GetCommandsByGuild(strToInt(guild.ID))
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
			} else {
				log.Printf("Registered command %s on guild %s", cmd.Name, guild.Name)
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

	token, err := generateRandomHash()
	if err != nil {
		message.Data.Content = "Error generating token"
		return message
	}
	c.ContextTokenCache.Store(token, i.GuildID)

	message.Data.Content = "To Login to the command line client, use the following command: `dfaas context create --url https://" + c.cfg.URLDomain + " --token " + token +
		"`\nTo Download the CLI Client: https://github.com/f4tal-err0r/discord_faas/"

	guild, err := c.Session.Guild(i.GuildID)
	if err != nil {
		log.Printf("Error getting guild info: %v", err)
	}
	fmt.Printf("Generated login token for guild %s: user: %s\n", guild.Name, i.Member.User.Username)

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

func cmdRouter(_ *discordgo.Interaction, _ *Client) *discordgo.InteractionResponse {
	message := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Command received",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}

	return message
}
