package main

import (
	"fmt"

	"github.com/f4tal-err0r/discord_faas/pkgs/client"
	"github.com/spf13/cobra"
)

var ctxtoken string
var url string

func init() {
	rootCmd.AddCommand(context)
	context.AddCommand(newContext)
	context.AddCommand(listContexts)
	context.AddCommand(currentContext)
	newContext.Flags().StringVarP(&ctxtoken, "token", "t", "", "Token generated via the /login command in discord")
	newContext.MarkFlagRequired("token")
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
		client.NewContext(url, ctxtoken)
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
