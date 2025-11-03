package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bradleyalexa/letstudy/data"
	"github.com/spf13/cobra"
)

var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "Create and manage quick notes",
	Long:  "A simple quick-note feature to save your thoughts instantly from CLI.",
	Run: func(cmd *cobra.Command, args []string) {
		handleNoteCommand(args)
	},
}

func init() {
	rootCmd.AddCommand(noteCmd)
}

func handleNoteCommand(args []string) {
	if err := data.OpenDB(); err != nil {
		fmt.Println("❌ Cannot open database:", err)
		return
	}
	defer data.CloseDB()
	data.CreateNotesTable()

	// Tidak ada argumen → tampilkan bantuan
	if len(args) == 0 {
		showNoteHelp()
		return
	}

	switch strings.ToLower(args[0]) {
	case "help", "-h", "--help":
		showNoteHelp()

	case "list":
		if err := data.ListQuickNotes(); err != nil {
			fmt.Println("❌ Error listing notes:", err)
		}

	case "view":
		if len(args) < 2 {
			fmt.Println("⚠️ Please specify a note ID.")
			return
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("⚠️ Invalid ID format.")
			return
		}
		if err := data.ViewQuickNoteByID(id); err != nil {
			fmt.Println("❌ Error viewing note:", err)
		}

	case "delete":
		if len(args) < 2 {
			fmt.Println("⚠️ Please specify a note ID.")
			return
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("⚠️ Invalid ID format.")
			return
		}
		if err := data.DeleteQuickNoteByID(id); err != nil {
			fmt.Println("❌ Error deleting note:", err)
		}

	case "search":
		if len(args) < 2 {
			fmt.Println("⚠️ Please specify a keyword.")
			return
		}
		keyword := strings.Join(args[1:], " ")
		if err := data.SearchQuickNotes(keyword); err != nil {
			fmt.Println("❌ Error searching notes:", err)
		}

	default:
		
		content := strings.Join(args, " ")
		if strings.TrimSpace(content) == "" {
			fmt.Println("⚠️ Please provide note content.")
			return
		}

		if err := data.InsertQuickNote(content); err != nil {
			fmt.Printf("❌ Failed to add note: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Note added successfully.")
	}
}


func showNoteHelp() {
	fmt.Print(`
📘 LetStudy Notes Command Manual
--------------------------------------------
Usage:
  letstudy note "text"            Add a quick note
  letstudy note list              Show all notes
  letstudy note view <id>         View a note by ID
  letstudy note delete <id>       Delete a note by ID
  letstudy note search "keyword"  Search notes by keyword
  letstudy note help              Show this help menu

Examples:
  letstudy note "Learn about Go routines"
  letstudy note list
  letstudy note view 3
  letstudy note delete 2
  letstudy note search "Go"

✨ Tip: Use quotes around multi-word notes or search terms.
`)
}
