package main

import (
	"fmt"

	"github.com/f4tal-err0r/discord_faas/pkgs/discord"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(refreshCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Functions-as-a-Service Server",
	Long:  "Discord FaaS kubernetes controller",
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Discord bot",
	Run: func(cmd *cobra.Command, args []string) {
		discord.StartDiscordBot()
	},
}

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Start Discord bot",
	Run: func(cmd *cobra.Command, args []string) {
		if d, err := discord.GetToken(); err != nil {
			fmt.Println(d)
		}
	},
}
