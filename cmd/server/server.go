package main

import (
	"github.com/f4tal-err0r/discord_faas/pkgs/server"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.AddCommand(startCmd)
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
		server.Start()
	},
}
