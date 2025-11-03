package cmd

import (
	"fmt"
	"strconv"

	"github.com/bradleyalexa/letstudy/data"
	"github.com/spf13/cobra"
)

// 🔹 reflectCmd utama
var reflectCmd = &cobra.Command{
	Use:   "reflect",
	Short: "View your task reflections",
	Long:  "See all your saved reflections for completed tasks.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := data.OpenDB(); err != nil {
			fmt.Println("❌ Failed to open database:", err)
			return
		}
		defer data.CloseDB()

		data.CreateReflectionTable()
		data.ListReflectionSummaries()
	},
}

// 🔹 reflect view [id]
var reflectViewCmd = &cobra.Command{
	Use:   "view [reflectionID]",
	Short: "View details of a specific reflection",
	Long:  "Display detailed information about a specific reflection entry by its ID.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("❌ Invalid reflection ID. Please provide a numeric value.")
			return
		}

		if err := data.OpenDB(); err != nil {
			fmt.Println("❌ Failed to open database:", err)
			return
		}
		defer data.CloseDB()

		data.ViewReflectionByID(id)
	},
}

func init() {
	rootCmd.AddCommand(reflectCmd)
	reflectCmd.AddCommand(reflectViewCmd)
}
