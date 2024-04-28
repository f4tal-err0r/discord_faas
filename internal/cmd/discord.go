package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(discordCmd)
	// startCmd.PersistentFlags()
}

var discordCmd = &cobra.Command{
	Short: "Commands for interfacing w/ Discord",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}
