/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"errors"
	"os"
	"time"
	"github.com/bradleyalexa/letstudy/data"
	"github.com/spf13/cobra"
	"github.com/manifoldco/promptui"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create new task",
	Long: `New command`,
	Run: func(cmd *cobra.Command, args []string) {
		createNewNote()
	},
}

func init() {
	taskCmd.AddCommand(newCmd)
}

type promptContent struct {
	errorMsg string
	label    string
}

func promptGetInput(pc promptContent) string {
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New(pc.errorMsg)
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     pc.label,
		Templates: templates,
		Validate:  validate,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Input: %s\n", result)

	return result
}

func createNewNote() {

	if _, err := os.Stat("./sqlite-database.db"); os.IsNotExist(err) {
		fmt.Println("⚠️ Database not found — creating a new one...")
	}

	if err := data.OpenDB(); err != nil {
	fmt.Printf("❌ Failed to open database: %v\n", err)
	return
	}

	data.CreateTable()


	contentPromptContent := promptContent{
		"Please provide what task you want to input.",
		"What task you want to do?",
	}
	content := promptGetInput(contentPromptContent)

	 promptDate := promptui.Prompt{
        Label: "Enter due date (YYYY-MM-DD) — press Enter to skip:",
    }

    dateInput, _ := promptDate.Run()

	 var taskDate *time.Time
    if dateInput != "" {
        parsedDate, err := time.Parse("2006-01-02", dateInput)
        if err != nil {
            fmt.Println("⚠️ Invalid date format. Please use YYYY-MM-DD.")
            return
        }
        taskDate = &parsedDate
    } else {
        taskDate = nil 
    }
	
	err := data.InsertNote(content, taskDate)
    if err != nil {
        fmt.Printf("❌ Failed to insert task: %v\n", err)
        return
    }

    fmt.Println("✅ Task added successfully! (status: not done)")
}
