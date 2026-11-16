package handlers

import (
	"net/http"
	"time"

	"school-management-backend/internal/config"
	"school-management-backend/internal/models"
	"school-management-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

type AttendanceEntry struct {
	StudentID uint                    `json:"student_id" binding:"required"`
	Status    models.AttendanceStatus `json:"status" binding:"required"`
}

type MarkAttendanceInput struct {
	ClassID uint              `json:"class_id" binding:"required"`
	Date    string            `json:"date" binding:"required"` // Format: YYYY-MM-DD
	Records []AttendanceEntry `json:"records" binding:"required"`
}

// MarkAttendance records or updates attendance for a class.
func MarkAttendance(c *gin.Context) {
	var input MarkAttendanceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var attendances []models.Attendance
	var absentStudentIDs []uint

	for _, record := range input.Records {
		attendances = append(attendances, models.Attendance{
			StudentID: record.StudentID,
			Date:      input.Date,
			Status:    record.Status,
		})

		if record.Status == models.StatusAbsent {
			absentStudentIDs = append(absentStudentIDs, record.StudentID)
		}
	}

	// Update the existing attendance record when the student and date already exist.
	err := config.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "student_id"}, {Name: "date"}},
		DoUpdates: clause.AssignmentColumns([]string{"status", "updated_at"}),
	}).Create(&attendances).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save attendance"})
		return
	}

	if len(absentStudentIDs) > 0 {
		var students []models.Student
		config.DB.Preload("User").Where("id IN ?", absentStudentIDs).Find(&students)

		var alertList []utils.AbsentStudentInfo
		for _, s := range students {
			alertList = append(alertList, utils.AbsentStudentInfo{
				Name:          s.User.Name,
				RollNumber:    s.RollNumber,
				GuardianPhone: s.GuardianPhone,
				Date:          input.Date,
			})
		}
		utils.SendAbsenceAlertsAsync(alertList)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Attendance recorded successfully!",
		"total_recorded": len(attendances),
		"absent_count":   len(absentStudentIDs),
	})
}

// GetClassAttendance returns attendance records for a class on a given date.
func GetClassAttendance(c *gin.Context) {
	classID := c.Param("class_id")
	date := c.DefaultQuery("date", time.Now().Format("2006-01-02"))

	var attendances []models.Attendance
	err := config.DB.
		Joins("JOIN students ON students.id = attendances.student_id").
		Preload("Student.User").
		Where("students.class_id = ? AND attendances.date = ?", classID, date).
		Find(&attendances).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attendance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"class_id":    classID,
		"date":        date,
		"attendances": attendances,
	})
}

// GetMyAttendanceReport returns the authenticated student's attendance summary.
func GetMyAttendanceReport(c *gin.Context) {
	userID, _ := c.Get("userID")

	var student models.Student
	if err := config.DB.Where("user_id = ?", userID).First(&student).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student profile not found for this user"})
		return
	}

	var records []models.Attendance
	config.DB.Where("student_id = ?", student.ID).Order("date desc").Find(&records)

	totalDays := len(records)
	presentCount := 0
	for _, r := range records {
		if r.Status == models.StatusPresent {
			presentCount++
		}
	}

	percentage := 0.0
	if totalDays > 0 {
		percentage = (float64(presentCount) / float64(totalDays)) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"student_id":    student.ID,
		"total_days":    totalDays,
		"present_days":  presentCount,
		"percentage":    percentage,
		"recent_logs":   records,
	})
}