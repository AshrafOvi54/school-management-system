package utils

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"
)

// CalculateGrade returns the grade and grade point for a score.
func CalculateGrade(marks float64) (string, float64) {
	switch {
	case marks >= 80:
		return "A+", 5.00
	case marks >= 70:
		return "A", 4.00
	case marks >= 60:
		return "A-", 3.50
	case marks >= 50:
		return "B", 3.00
	case marks >= 40:
		return "C", 2.00
	case marks >= 33:
		return "D", 1.00
	default:
		return "F", 0.00
	}
}

// GenerateVerificationHash creates a result signature using the configured secret key.
func GenerateVerificationHash(studentID, examID uint, gpa float64) string {
	secretKey := os.Getenv("RESULT_SECRET_KEY")
	if secretKey == "" {
		log.Fatal("FATAL SECURITY ERROR: RESULT_SECRET_KEY environment variable is missing in .env")
	}
	data := fmt.Sprintf("STUDENT:%d_EXAM:%d_GPA:%.2f_SECRET:%s", studentID, examID, gpa, secretKey)
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}
