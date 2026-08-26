package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "migrator",
	Short: "Deep Tech Migration Batch Processor",
	Long:  `An Enterprise AI tool that converts legacy codebases (PHP, Python) to modern Go using AST graphs and Agentic TDD.`,
}

func init() {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Web UI Dashboard",
	Run: func(cmd *cobra.Command, args []string) {
		port := "8080"
		StartAPIServer(port)
	},
}

// Execute runs the root CLI command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
