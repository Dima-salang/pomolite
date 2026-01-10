/*
Copyright © 2023 LUIS GABRIELLE PUTAN <luisgabrielle1026@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/Dima-salang/pomolite/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pomo",
	Short: "Lightweight CLI Pomodoro Timer",
	Long: fmt.Sprintf(`%s PomoLite is a lightweight CLI Pomodoro application desgned for students to efficiently accomplish tasks and maximize their learning potential
########################################################################################

Example usage:
	
pomo start -m 30 -b 5 :: Starts a 30-minute Pomodoro timer with a 5-minute break

To pause the Pomodoro Timer, you can press the 'p' button on your keyboard.

To unpause or resume the Pomodoro Timer, press 'r'.

To quit the Pomodoro Timer and save it for stats, press 'q'.



Developed by PUTAN LUIS GABRIELLE <luisgabrielle1026@gmail.com>

#######################################################################################`, timer.ASCIIArt),
	Run: func(cmd *cobra.Command, args []string) {
		storage, err := timer.NewSQLiteStorage("./pomodoro.db")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		defer storage.Close()

		// Default values for home page if started from here
		m := timer.NewMainModel(30*time.Minute, 5*time.Minute, "Work", storage)
		p := tea.NewProgram(m, tea.WithAltScreen())

		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.PomoLite.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
