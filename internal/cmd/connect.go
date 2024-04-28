package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(connectCmd)
	// startCmd.PersistentFlags()
}

var connectCmd = &cobra.Command{
	Short: "[HACK] Lightweight admin interface to connect to remote platform .",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("placehold")
	},
}
