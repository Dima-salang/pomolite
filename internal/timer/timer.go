package timer

import (
	"fmt"
)

// CheckInput validates the user's customized Pomodoro durations.
func CheckInput(minutes int, breakMinutes int) bool {
	if minutes <= 0 || breakMinutes <= 0 {
		fmt.Println("Error: Invalid input. Minutes and break must be greater than 0.")
		return false
	}
	return true
}
