/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/bradleyalexa/letstudy/data"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		doneFlag, _ := cmd.Flags().GetBool("done")

		if doneFlag {
			fmt.Println("📋 Showing completed tasks:")
			data.DisplayTasks("done")
		} else {
			fmt.Println("📋 Showing active tasks:")
			data.DisplayTasks("not done")
		}
	},
}

func init() {
	taskCmd.AddCommand(listCmd)
	listCmd.Flags().Bool("done", false, "Show only completed tasks")
}
