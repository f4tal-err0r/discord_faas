package main

import (
	"fmt"
	"log"

	"github.com/f4tal-err0r/discord_faas/pkgs/discord"
	"github.com/spf13/cobra"
)

func init() {
	discordCmd.AddCommand(login)
}

var discordCmd = &cobra.Command{
	Use:   "discord",
	Short: "Commands for interfacing w/ Discord",
}

var login = &cobra.Command{
	Use:   "login",
	Short: "Auth to your Discord Guild",
	Run: func(cmd *cobra.Command, args []string) {
		if callback, err := discord.StartAuth(); err != nil {
			log.Fatal(err)
		} else {
			fmt.Print(callback)
		}
	},
}
