package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var rootCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "Discord Functions-as-a-Service",
	Long:  "Bot frontend to integrate discord with an eventing platform",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
