package main

import (
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Functions-as-a-Service Server",
	Long:  "Discord FaaS kubernetes controller",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}
