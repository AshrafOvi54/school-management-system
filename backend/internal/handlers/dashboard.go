package handlers

import (
	"net/http"
	"time"

	"school-management-backend/internal/config"
	"school-management-backend/internal/models"

	"github.com/gin-gonic/gin"
)

type CreateNoticeInput struct {
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required"`
	TargetRole string `json:"target_role"`
}

// CreateNotice publishes a notice.
func CreateNotice(c *gin.Context) {
	var input CreateNoticeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	target := input.TargetRole
	if target == "" {
		target = "all"
	}

	notice := models.Notice{
		Title:      input.Title,
		Content:    input.Content,
		TargetRole: target,
	}

	if err := config.DB.Create(&notice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to post notice"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Notice published successfully!", "notice": notice})
}

// GetNotices returns notices visible to the authenticated user's role.
func GetNotices(c *gin.Context) {
	userRole, _ := c.Get("userRole")

	var notices []models.Notice
	config.DB.Where("target_role = 'all' OR target_role = ?", userRole).Order("created_at desc").Find(&notices)

	c.JSON(http.StatusOK, gin.H{"notices": notices})
}

// GetDashboardStats returns school and daily attendance statistics.
func GetDashboardStats(c *gin.Context) {
	var totalStudents int64
	var totalTeachers int64
	var totalClasses int64
	var totalSubjects int64
	var totalExams int64
	var todayPresent int64
	var todayAbsent int64

	today := time.Now().Format("2006-01-02")

	config.DB.Model(&models.Student{}).Count(&totalStudents)
	config.DB.Model(&models.User{}).Where("role = ?", models.RoleTeacher).Count(&totalTeachers)
	config.DB.Model(&models.Class{}).Count(&totalClasses)
	config.DB.Model(&models.Subject{}).Count(&totalSubjects)
	config.DB.Model(&models.Exam{}).Count(&totalExams)

	config.DB.Model(&models.Attendance{}).Where("date = ? AND status = ?", today, models.StatusPresent).Count(&todayPresent)
	config.DB.Model(&models.Attendance{}).Where("date = ? AND status = ?", today, models.StatusAbsent).Count(&todayAbsent)

	attendanceRate := 0.0
	totalAttendanceToday := todayPresent + todayAbsent
	if totalAttendanceToday > 0 {
		attendanceRate = (float64(todayPresent) / float64(totalAttendanceToday)) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"total_students":              totalStudents,
		"total_teachers":              totalTeachers,
		"total_classes":               totalClasses,
		"total_subjects":              totalSubjects,
		"total_exams":                 totalExams,
		"today_present":               todayPresent,
		"today_absent":                todayAbsent,
		"today_attendance_rate_%":     attendanceRate,
		"today_attendance_percentage": attendanceRate,
	})
}
