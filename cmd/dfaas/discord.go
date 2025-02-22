package main

import (
	"fmt"
	"log"

	"github.com/f4tal-err0r/discord_faas/pkgs/client"
	"github.com/spf13/cobra"
)

func init() {
	discordCmd.AddCommand(login)
	// discordCmd.AddCommand(currentUser)
}

var discordCmd = &cobra.Command{
	Use:   "discord",
	Short: "Commands for interfacing w/ Discord",
}

var login = &cobra.Command{
	Use:   "login",
	Short: "Auth to your Discord Guild",
	Run: func(cmd *cobra.Command, args []string) {
		dauth := client.NewUserAuth()
		if token, err := dauth.StartAuth(); err != nil {
			log.Fatal(err)
		} else {
			fmt.Printf("Auth token: %s\n", token.AccessToken)
		}
	},
}

// var currentUser = &cobra.Command{
// 	Use:   "current-user",
// 	Short: "Get the current user",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		fmt.Printf("Current User: %s", client.GetCurrentUser().String())
// 	},
// }
