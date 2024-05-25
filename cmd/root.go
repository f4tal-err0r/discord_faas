package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "Discord FaaS",
	Short: "Discord Functions-as-a-Service",
	Long:  "Bot frontend to integrate discord with an eventing platform",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
