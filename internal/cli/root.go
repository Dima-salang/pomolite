/*
Copyright © 2023 LUIS GABRIELLE PUTAN <luisgabrielle1026@gmail.com>
*/
package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/Dima-salang/pomolite/internal/storage"
	"github.com/Dima-salang/pomolite/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pomo",
	Short: "Lightweight CLI Pomodoro Timer",
	Long: fmt.Sprintf(`%s PomoLite is a lightweight CLI Pomodoro application designed for students to efficiently accomplish tasks and maximize their learning potential
########################################################################################

Example usage:
	
pomo start -m 30 -b 5 :: Starts a 30-minute Pomodoro timer with a 5-minute break

To pause the Pomodoro Timer, you can press the 'p' button on your keyboard.

To unpause or resume the Pomodoro Timer, press 'r'.

To quit the Pomodoro Timer and save it for stats, press 'q'.



Developed by PUTAN LUIS GABRIELLE <luisgabrielle1026@gmail.com>

#######################################################################################`, ui.ASCIIArt),
	Run: func(cmd *cobra.Command, args []string) {
		store, err := storage.NewSQLiteStorage("")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		defer store.Close()

		m := ui.NewMainModel(30*time.Minute, 5*time.Minute, "Work", store)
		p := tea.NewProgram(m, tea.WithAltScreen())

		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
