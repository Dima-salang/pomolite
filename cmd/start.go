package cmd

import (
	"fmt"
	"time"

	"github.com/Dima-salang/pomolite/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var minutes int
var breakMinutes int
var label string

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start a Pomodoro Timer.",
	Long: `Start a Pomodoro Timer.

	FLAGS:
	-l : label for the work session
	-m : minutes of work
	-b : minutes of break`,
	Run: func(cmd *cobra.Command, args []string) {
		// check for the validity of the input
		if !timer.CheckInput(minutes, breakMinutes) {
			return
		}
		storage, err := timer.NewSQLiteStorage("./pomodoro.db")

		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		defer storage.Close()

		totalWorkDuration := time.Duration(minutes) * time.Minute
		totalBreakDuration := time.Duration(breakMinutes) * time.Minute

		m := timer.NewMainModel(totalWorkDuration, totalBreakDuration, label, storage)
		p := tea.NewProgram(m, tea.WithAltScreen())

		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringVarP(&label, "label", "l", "Work", "label for the work session")
	startCmd.Flags().IntVarP(&minutes, "minutes", "m", 30, "minutes to work")
	startCmd.Flags().IntVarP(&breakMinutes, "break", "b", 5, "minutes to take a break")
}
