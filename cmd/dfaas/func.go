package main

import (
	"fmt"

	"github.com/f4tal-err0r/discord_faas/internal/client"
	"github.com/f4tal-err0r/discord_faas/pkgs/runtimes"
	"github.com/spf13/cobra"
)

var runtime string

func init() {
	rootCmd.AddCommand(funcRootCmd)
	funcRootCmd.AddCommand(funcCreateCmd)
	funcCreateCmd.Flags().StringVar(&runtime, "runtime", "", "Runtime to use for genearting a function")
	funcCreateCmd.MarkFlagRequired("runtime")
	funcRootCmd.AddCommand(funcRuntimeCmd)
	funcRootCmd.AddCommand(funcList)
}

var funcRootCmd = &cobra.Command{
	Use:   "func",
	Short: "Create, Deploy, and Manage Functions",
}

var funcCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a function",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("No function name provided")
			return
		}

		if err := runtimes.GenerateFunc(args[0], runtime); err != nil {
			fmt.Printf("Unable to generate function: %v", err)
			return
		}
	},
}

var funcRuntimeCmd = &cobra.Command{
	Use:   "runtimes",
	Short: "List available runtimes",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: List available runtimes from api
	},
}

var funcList = &cobra.Command{
	Use:   "ls",
	Short: "List available functions",
	Run: func(cmd *cobra.Command, args []string) {
		client.ListFunctions()
	},
}
