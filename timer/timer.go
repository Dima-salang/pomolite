package timer

import (
	"fmt"
)

func CheckInput(minutes int, breakMinutes int) bool {
	// if minutes, breakMinutes, or workSeconds is less than or equal to 0, return false
	if minutes <= 0 || breakMinutes <= 0 {
		fmt.Println("Error: Invalid input. Minutes, break, and seconds must be greater than 0.")
		return false
	}

	return true
}
