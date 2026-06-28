package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/Dima-salang/pomolite/internal/storage"
	"github.com/Dima-salang/pomolite/internal/timer"
	"github.com/Dima-salang/pomolite/internal/ui"
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
		if !timer.CheckInput(minutes, breakMinutes) {
			return
		}
		store, err := storage.NewSQLiteStorage("")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		defer store.Close()

		totalWorkDuration := time.Duration(minutes) * time.Minute
		totalBreakDuration := time.Duration(breakMinutes) * time.Minute

		m := ui.NewMainModel(totalWorkDuration, totalBreakDuration, label, store)
		p := tea.NewProgram(m, tea.WithAltScreen())

		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringVarP(&label, "label", "l", "Work", "label for the work session")
	startCmd.Flags().IntVarP(&minutes, "minutes", "m", 30, "minutes to work")
	startCmd.Flags().IntVarP(&breakMinutes, "break", "b", 5, "minutes to take a break")
}
