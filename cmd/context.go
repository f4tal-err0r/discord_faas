package main

import (
	"fmt"

	"github.com/f4tal-err0r/discord_faas/pkgs/client"
	"github.com/spf13/cobra"
)

var guildid string
var url string

func init() {
	rootCmd.AddCommand(context)
	context.AddCommand(newContext)
	context.AddCommand(listContexts)
	context.AddCommand(currentContext)
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

var currentContext = &cobra.Command{
	Use:   "current",
	Short: "Show Current Context",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Current Server Context: %v", client.GetCurrentContext().GuildName)
	},
}
