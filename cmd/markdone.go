/*
Copyright © 2025 NAME HERE
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
	Long:  "Mark a task as done by its ID, and reflect on what you learned from completing it.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskID, err := strconv.Atoi(args[0])
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
}
