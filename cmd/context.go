package main

import (
	"github.com/f4tal-err0r/discord_faas/pkgs/client"
	"github.com/spf13/cobra"
)

var guildid string
var url string

func init() {
	rootCmd.AddCommand(context)
	context.AddCommand(newContext)
	context.AddCommand(listContexts)
	newContext.Flags().StringVarP(&guildid, "guildid", "", "", "GuildID of server")
	newContext.MarkFlagRequired("guildid")
	newContext.Flags().StringVarP(&url, "url", "", "", "GuildID of server")
	newContext.MarkFlagRequired("url")

}

var context = &cobra.Command{
	Use:     "context",
	Aliases: []string{"ctx"},
	Short:   "Working context of Discord server",
}

var newContext = &cobra.Command{
	Use:   "connect",
	Short: "Working context of Discord server",
	Run: func(cmd *cobra.Command, args []string) {
		client.NewContext(url, guildid)
	},
}

var listContexts = &cobra.Command{
	Use:   "ls",
	Short: "List available contexts",
	Run: func(cmd *cobra.Command, args []string) {
		client.ListContexts()
	},
}
