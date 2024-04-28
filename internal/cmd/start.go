package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Short: "Starts the Discord FAAS Platform",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}
