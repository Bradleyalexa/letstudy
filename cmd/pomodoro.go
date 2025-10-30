package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"
	"github.com/gen2brain/beeep"
	"github.com/bradleyalexa/letstudy/data"
	"github.com/spf13/cobra"
)

var duration int
var sessionType string

var pomodoroCmd = &cobra.Command{
	Use:   "pomodoro",
	Short: "Start a Pomodoro timer",
	Long:  "Pomodoro timer with pause/resume, sound, and persistent tracking.",
	Run: func(cmd *cobra.Command, args []string) {
		startPomodoro(sessionType, duration)
	},
}

func init() {
	rootCmd.AddCommand(pomodoroCmd)
	pomodoroCmd.Flags().IntVarP(&duration, "minutes", "m", 25, "Set duration in minutes (default 25)")
	pomodoroCmd.Flags().StringVarP(&sessionType, "type", "t", "focus", "Session type: focus, break, or longbreak")
}

func startPomodoro(sessionType string, duration int) {
	if err := data.OpenDB(); err != nil {
		fmt.Println("❌ Cannot open database:", err)
		return
	}
	defer data.CloseDB()
	data.CreatePomodoroTable()

	fmt.Printf("🕒 Starting %s session for %d minutes...\n", strings.Title(sessionType), duration)
	fmt.Println("Press [p]ause / [r]esume / [q]uit anytime.")

	done := make(chan bool)
	pause := make(chan bool)
	resume := make(chan bool)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	startTime := time.Now()

	status := "completed"

	go func() {
		remaining := time.Duration(duration) * time.Minute
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		isPaused := false

		for {
			select {
			case <-pause:
				isPaused = true
			case <-resume:
				isPaused = false
			case <-ticker.C:
				if !isPaused {
					remaining -= time.Second
					printProgressBar(remaining, duration)
					if remaining <= 0 {
						fmt.Println("\n✅ Time’s up!")
						beeep.Alert("Pomodoro Finished!", "Time’s up!", "")
						done <- true
						return
					}
				}
			case <-done:
				return
			}
		}
	}()

	go func() {
		for {
			var input string
			fmt.Scanln(&input)
			switch input {
			case "p":
				fmt.Println("⏸ Paused")
				pause <- true
			case "r":
				fmt.Println("▶️ Resumed")
				resume <- true
			case "q":
				fmt.Println("🛑 Quit")
				status = "quit"
				done <- true
				return
			}
		}
	}()

	select {
	case <-done:
	case <-interrupt:
		fmt.Println("\n🛑 Interrupted.")
		status = "interrupted"
	}

	data.InsertPomodoroSession(sessionType, duration, status, startTime, time.Now())
	fmt.Println("💾 Session saved to database.")
}

func printProgressBar(remaining time.Duration, totalMinutes int) {
	totalSeconds := totalMinutes * 60
	remSec := int(remaining.Seconds())
	progress := int((float64(totalSeconds-remSec) / float64(totalSeconds)) * 30)
	bar := strings.Repeat("█", progress) + strings.Repeat("-", 30-progress)
	fmt.Printf("\r[%s] %02d:%02d remaining", bar, remSec/60, remSec%60)
}
