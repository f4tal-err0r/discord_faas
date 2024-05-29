package main

import (
	"log"
	"time"

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
		for {
			time.Sleep(10 * time.Second)
			log.Println("Test Out Live Cmd")
		}
	},
}
