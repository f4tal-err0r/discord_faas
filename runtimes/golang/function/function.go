package function

import "github.com/bwmarrin/discordgo"

func DiscordFunction() *discordgo.Message {
	return &discordgo.Message{
		Content: "Hello World!",
	}
}
