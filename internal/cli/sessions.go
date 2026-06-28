package cli

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Dima-salang/pomolite/internal/storage"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// sessionsCmd represents the sessions command
var sessionsCmd = &cobra.Command{
	Use:   "sessions",
	Short: "list the sessions",
	Long: `List the sessions.

	This command lists all the sessions saved in the database.
	Each session includes the label, start time, and end time.
	The sessions are ordered by start time in descending order.`,
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		store, err := storage.NewSQLiteStorage("")
		if err != nil {
			fmt.Println(color.RedString("Error: %v", err))
			return
		}
		defer store.Close()

		sessions, err := store.List(limit)
		if err != nil {
			fmt.Println(color.RedString("Error: %v", err))
			return
		}
		if len(sessions) == 0 {
			fmt.Println(color.YellowString("No sessions found."))
			return
		}

		ansi := regexp.MustCompile("\x1b\\[[0-9;]*m")
		visibleLen := func(s string) int {
			return len([]rune(ansi.ReplaceAllString(s, "")))
		}
		padRightANSI := func(s string, width int) string {
			v := visibleLen(s)
			if v >= width {
				return s
			}
			return s + strings.Repeat(" ", width-v)
		}

		idW := visibleLen("ID")
		labelW := visibleLen("Label")
		startW := visibleLen("Start Time")
		endW := visibleLen("End Time")
		durW := visibleLen("Duration")

		for _, s := range sessions {
			idStr := fmt.Sprintf("%d", s.ID)
			if len(idStr) > idW {
				idW = len(idStr)
			}
			if visibleLen(s.Label) > labelW {
				labelW = visibleLen(s.Label)
			}
			startStr := s.StartTime.Format("2006-01-02 15:04:05")
			if visibleLen(startStr) > startW {
				startW = visibleLen(startStr)
			}
			endStr := s.EndTime.Format("2006-01-02 15:04:05")
			if visibleLen(endStr) > endW {
				endW = visibleLen(endStr)
			}
			durStr := s.SessionDuration.Round(time.Second).String()
			if visibleLen(durStr) > durW {
				durW = visibleLen(durStr)
			}
		}

		hID := color.CyanString("ID")
		hLabel := color.CyanString("Label")
		hStart := color.CyanString("Start Time")
		hEnd := color.CyanString("End Time")
		hDur := color.CyanString("Duration")

		sepLen := idW + labelW + startW + endW + durW + 4*2
		fmt.Printf("%s  %s  %s  %s  %s\n",
			padRightANSI(hID, idW),
			padRightANSI(hLabel, labelW),
			padRightANSI(hStart, startW),
			padRightANSI(hEnd, endW),
			padRightANSI(hDur, durW),
		)
		fmt.Println(strings.Repeat("-", sepLen))

		for i, s := range sessions {
			idStr := fmt.Sprintf("%d", s.ID)
			labelColored := color.GreenString(s.Label)
			if i%2 == 1 {
				labelColored = color.YellowString(s.Label)
			}
			startStr := s.StartTime.Format("2006-01-02 15:04:05")
			endStr := s.EndTime.Format("2006-01-02 15:04:05")
			durStr := s.SessionDuration.Round(time.Second).String()
			durColored := color.MagentaString("%s", durStr)

			fmt.Printf("%s  %s  %s  %s  %s\n",
				padRightANSI(idStr, idW),
				padRightANSI(labelColored, labelW),
				padRightANSI(startStr, startW),
				padRightANSI(endStr, endW),
				padRightANSI(durColored, durW),
			)
		}

		if limit > 0 {
			fmt.Println()
			fmt.Println(color.HiBlackString("Showing last %d session(s).", limit))
		}
	},
}

func init() {
	rootCmd.AddCommand(sessionsCmd)
	sessionsCmd.Flags().IntP("limit", "l", 0, "number of sessions to list")
}
