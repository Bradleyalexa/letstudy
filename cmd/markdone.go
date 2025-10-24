/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"strconv"
	"github.com/bradleyalexa/letstudy/data"
	"github.com/spf13/cobra"
)

// markdoneCmd represents the markdone command
var markdoneCmd = &cobra.Command{
	Use:   "markdone [taskID]",
	Short: "Mark a task as done",
	Long:  "Mark a task as done by its ID, updating its status from 'not done' to 'done'.",
	Args:  cobra.ExactArgs(1), //ensure only exact 1 arguments received
	Run: func(cmd *cobra.Command, args []string) {
		taskID, err := strconv.Atoi(args[0]) //input type based string, conv into integer to send to func
		if err != nil {
			fmt.Println("❌ Invalid task ID. Please provide a numeric value.")
			return
		}

		err = data.MarkTaskDone(taskID)
		if err != nil {
			fmt.Printf("❌ %v\n", err)
			return
		}
	},
}

func init() {
	taskCmd.AddCommand(markdoneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// markdoneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// markdoneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
