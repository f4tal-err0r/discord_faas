package main

import (
	"fmt"
	"os"

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
	funcRootCmd.AddCommand(funcDeploy)
}

var funcRootCmd = &cobra.Command{
	Use:   "func",
	Short: "Create, Deploy, and Manage Functions",
}

var funcCreateCmd = &cobra.Command{
	Use:   "create [function_name]",
	Short: "Create a function",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		functionName := args[0]
		if err := runtimes.GenerateFunc(functionName, runtime); err != nil {
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

var funcDeploy = &cobra.Command{
	Use:   "deploy [function_path]",
	Short: "Upload and deploy a function",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("No file path provided")
			return
		}
		if _, err := os.Stat(args[0]); os.IsNotExist(err) {
			fmt.Printf("Directory does not exist: %s\n", args[0])
			return
		}

		if err := client.DeployFunc(args[0]); err != nil {
			fmt.Printf("Unable to deploy function: %v", err)
			return
		}
	},
}
