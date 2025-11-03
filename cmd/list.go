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
	Short: "List of to-do tasks",
	Long: `Showing lists of to-do tasks`,
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
