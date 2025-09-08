package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var runtime string

func init() {
	rootCmd.AddCommand(funcRootCmd)
	funcRootCmd.AddCommand(funcCreateCmd)
	funcCreateCmd.Flags().StringVar(&runtime, "runtime", "", "Runtime to use for genearting a function")
	funcCreateCmd.MarkFlagRequired("runtime")
	funcRootCmd.AddCommand(funcRuntimeCmd)
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

		// TODO: Clone func template from repo
	},
}

var funcRuntimeCmd = &cobra.Command{
	Use:   "runtimes",
	Short: "List available runtimes",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: List available runtimes from api
	},
}
