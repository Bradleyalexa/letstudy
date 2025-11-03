package cmd

import (
	"fmt"
	"time"

	"github.com/bradleyalexa/letstudy/data"
	"github.com/spf13/cobra"
)

// remindCmd represents the remind command
var remindCmd = &cobra.Command{
	Use:   "remind",
	Short: "Show reminders for tasks that are due soon",
	Long: `Checks your task list and shows reminders for tasks
that are approaching their due date (e.g., within 24 hours).`,
	Run: func(cmd *cobra.Command, args []string) {
		checkReminders()
	},
}

func init() {
	taskCmd.AddCommand(remindCmd)
}

func checkReminders() {
	if err := data.OpenDB(); err != nil {
		fmt.Println("❌ Cannot open database:", err)
		return
	}
	defer data.CloseDB()

	tasks, err := data.GetUpcomingTasks(24 * time.Hour) // check 24 hours ahead
	if err != nil {
		fmt.Printf("❌ Failed to fetch tasks: %v\n", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("✅ No upcoming tasks in the next 24 hours.")
		return
	}

	fmt.Println("⏰ Upcoming Tasks (within 24 hours):")
	fmt.Println("------------------------------------")
	for _, t := range tasks {
		fmt.Printf("[%d] %s — Due: %s (%s left)\n",
			t.ID, t.Content, t.Due.Format("2 Jan 2006, 15:04"),
			time.Until(t.Due).Truncate(time.Minute))
		beep()
	}
}

// Simple notification beep
func beep() {
	fmt.Print("\a") // plays terminal beep sound (cross-platform)
}
