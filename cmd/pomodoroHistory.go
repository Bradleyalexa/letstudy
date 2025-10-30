package cmd

import (
	"fmt"

	"github.com/bradleyalexa/letstudy/data"
	"github.com/spf13/cobra"
)

var pomodoroHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Show Pomodoro session history",
	Run: func(cmd *cobra.Command, args []string) {
		showPomodoroHistory()
	},
}

func init() {
	pomodoroCmd.AddCommand(pomodoroHistoryCmd)
}

func showPomodoroHistory() {
	if err := data.OpenDB(); err != nil {
		fmt.Println("❌ Failed to open database:", err)
		return
	}
	defer data.CloseDB()

	data.CreatePomodoroTable()
	sessions, err := data.GetAllPomodoroSessions()
	if err != nil {
		fmt.Println("❌ Failed to fetch history:", err)
		return
	}

	if len(sessions) == 0 {
		fmt.Println("📭 No Pomodoro sessions yet.")
		return
	}

	fmt.Println("\n📜 Pomodoro History:")
	fmt.Println("────────────────────────────────────────────────────────────")
	for _, s := range sessions {
		start := s.StartTime.Format("2006-01-02 15:04")
		end := s.EndTime.Format("15:04")
		fmt.Printf("🕒 %-8s | %2d min | %s → %s | %-10s\n",
			s.SessionType, s.Duration, start, end, s.Status)
	}
	fmt.Println("────────────────────────────────────────────────────────────")
	fmt.Printf("Total sessions: %d\n", len(sessions))
}
