package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/f4tal-err0r/discord_faas/pkgs/client"
	"github.com/f4tal-err0r/discord_faas/pkgs/platform"
)

var runtime string

func init() {
	rootCmd.AddCommand(funcRootCmd)
	funcRootCmd.AddCommand(funcCreateCmd)
	funcCreateCmd.Flags().StringVar(&runtime, "runtime", "", "Runtime to use for genearting a function")
	funcCreateCmd.MarkFlagRequired("runtime")
	funcRootCmd.AddCommand(funcRuntimeCmd)
	funcRootCmd.AddCommand(funcDeployCmd)
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
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Print(err)
			return
		}

		fp := cwd + "/" + args[0]

		if err := platform.FunctionTemplate(fp, runtime); err != nil {
			fmt.Println(err)
			return
		}
	},
}

var funcRuntimeCmd = &cobra.Command{
	Use:   "runtimes",
	Short: "List available runtimes",
	Run: func(cmd *cobra.Command, args []string) {
		for runtime := range platform.UserLangDir {
			fmt.Println(runtime)
		}
	},
}

var funcDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a function",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			fmt.Println("No directory provided")
			return
		}

		dfaasPath := args[0]

		if _, err := os.Stat(filepath.Join(dfaasPath, "/dfaas.yaml")); os.IsNotExist(err) {
			fmt.Printf("ERR: dfaas.yaml not found in %s\n", dfaasPath)
			return
		} else {
			c, err := client.GetCurrentContext()
			if err != nil {
				fmt.Println("ERR: GetCurrentContext", err)
				return
			}
			if err := client.DeployFunc(c, dfaasPath); err != nil {
				fmt.Println("ERR: DeployFunc", err)
				return
			}
		}
	},
}
