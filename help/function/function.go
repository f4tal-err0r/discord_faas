package function

import (
	"encoding/json"

	"github.com/bwmarrin/discordgo"
)

func Handler(content *DiscordContent) (*DiscordResp, error) {
	embed := discordgo.MessageEmbed{
		Title:       content.Command,
		Description: content.Command,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Command",
				Value:  content.Command,
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Discord FaaS",
		},
	}

	//convert embed struct to json
	embedJSON, err := json.Marshal(embed)
	if err != nil {
		return nil, err
	}

	return &DiscordResp{
		Message: content.Command,
		Embed:   string(embedJSON),
		Files:   nil,
	}, nil
}
