package cmd

import (
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "Discord FaaS server",
	Short: "Discord Functions-as-a-Service Server",
	Long:  "Discord FaaS kubernetes controller",
}

func init() {
	rootCmd.AddCommand(serverCmd)
	// startCmd.PersistentFlags()
}
