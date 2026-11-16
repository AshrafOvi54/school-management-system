package utils

import (
	"fmt"
	"time"
)

type AbsentStudentInfo struct {
	Name          string
	RollNumber    string
	GuardianPhone string
	Date          string
}

// SendAbsenceAlertsAsync processes absence alerts in a background goroutine.
func SendAbsenceAlertsAsync(absentStudents []AbsentStudentInfo) {
	go func() {
		for _, s := range absentStudents {
			// Simulate the delay of sending a message through a notification gateway.
			time.Sleep(1 * time.Second)
			fmt.Printf("📢 [GOROUTINE ALERT] SMS sent to Guardian (%s): 'Dear Guardian, student %s (Roll: %s) was ABSENT on %s.'\n",
				s.GuardianPhone, s.Name, s.RollNumber, s.Date)
		}
	}()
}