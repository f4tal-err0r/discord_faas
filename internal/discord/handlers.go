package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/f4tal-err0r/discord_faas/pkgs/config"
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

func (c *Client) RegisterCommands() error {
	c.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			if _, ok := defaultCommands[i.ApplicationCommandData().Name]; ok {
				s.InteractionRespond(i.Interaction, defaultCommands[i.ApplicationCommandData().Name].Function(i.Interaction, c))
			}
			s.InteractionRespond(i.Interaction, cmdRouter(i.Interaction, c))
		}
	})
	var commands []discordgo.ApplicationCommand

	//init default commands
	for k, v := range defaultCommands {
		commands = append(commands, discordgo.ApplicationCommand{
			Name:        k,
			Description: v.Description,
		})
	}

	guilds, err := c.Session.UserGuilds(100, "", "", false)
	if err != nil {
		return fmt.Errorf("error getting guilds: %v", err)
	}

	for _, guild := range guilds {
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

func cmdRouter(i *discordgo.Interaction, c *Client) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "TODO: Documentation for the help command",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}
}
