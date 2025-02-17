package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/f4tal-err0r/discord_faas/pkgs/platform"
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
		if err := platform.FunctionTemplate(args[0], false, runtime); err != nil {
			fmt.Println(err)
			return
		}
	},
}

var funcRuntimeCmd = &cobra.Command{
	Use:   "runtimes",
	Short: "List available runtimes",
	Run: func(cmd *cobra.Command, args []string) {
		for runtime, _ := range platform.UserLangDir {
			fmt.Println(runtime)
		}
	},
}
