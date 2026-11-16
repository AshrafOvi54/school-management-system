package handlers

import (
	"fmt"
	"net/http"

	"school-management-backend/internal/config"
	"school-management-backend/internal/models"
	"school-management-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

type CreateSubjectInput struct {
	ClassID uint   `json:"class_id" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Code    string `json:"code" binding:"required"`
}

type CreateExamInput struct {
	Name string `json:"name" binding:"required"`
	Term string `json:"term"`
}

type MarkEntryInput struct {
	ExamID        uint    `json:"exam_id" binding:"required"`
	StudentID     uint    `json:"student_id" binding:"required"`
	SubjectID     uint    `json:"subject_id" binding:"required"`
	MarksObtained float64 `json:"marks_obtained" binding:"required,min=0,max=100"`
}

// CreateSubject creates a subject.
func CreateSubject(c *gin.Context) {
	var input CreateSubjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subject := models.Subject{
		ClassID: input.ClassID,
		Name:    input.Name,
		Code:    input.Code,
	}

	if err := config.DB.Create(&subject).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subject"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Subject created successfully!", "subject": subject})
}

// CreateExam creates an exam.
func CreateExam(c *gin.Context) {
	var input CreateExamInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exam := models.Exam{
		Name: input.Name,
		Term: input.Term,
	}

	if err := config.DB.Create(&exam).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create exam"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Exam created successfully!", "exam": exam})
}

// SubmitMarks records marks and calculates the grade.
func SubmitMarks(c *gin.Context) {
	var input MarkEntryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grade, point := utils.CalculateGrade(input.MarksObtained)

	mark := models.Mark{
		ExamID:        input.ExamID,
		StudentID:     input.StudentID,
		SubjectID:     input.SubjectID,
		MarksObtained: input.MarksObtained,
		Grade:         grade,
		GradePoint:    point,
	}

	// Update the existing mark when the exam, student, and subject already exist.
	err := config.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "exam_id"}, {Name: "student_id"}, {Name: "subject_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"marks_obtained", "grade", "grade_point", "updated_at"}),
	}).Create(&mark).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save mark"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Marks recorded successfully!",
		"mark":    mark,
	})
}

// GetStudentReportCard returns the student's results and GPA.
func GetStudentReportCard(c *gin.Context) {
	studentID := c.Param("student_id")
	examID := c.Param("exam_id")

	var student models.Student
	if err := config.DB.Preload("User").Preload("Class").First(&student, studentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	var marks []models.Mark
	config.DB.Preload("Subject").Where("student_id = ? AND exam_id = ?", studentID, examID).Find(&marks)

	if len(marks) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No marks found for this exam"})
		return
	}

	totalPoints := 0.0
	hasFailed := false

	for _, m := range marks {
		if m.Grade == "F" {
			hasFailed = true
		}
		totalPoints += m.GradePoint
	}

	gpa := 0.0
	if !hasFailed && len(marks) > 0 {
		gpa = totalPoints / float64(len(marks))
	}


	verificationHash := utils.GenerateVerificationHash(student.ID, marks[0].ExamID, gpa)
	gpaStr := fmt.Sprintf("%.2f", gpa)

	var verifyRec models.ResultVerification
	if err := config.DB.Where("student_id = ? AND exam_id = ?", student.ID, marks[0].ExamID).First(&verifyRec).Error; err != nil {
		config.DB.Create(&models.ResultVerification{
			StudentID: student.ID,
			ExamID:    marks[0].ExamID,
			GPA:       gpaStr,
			Hash:      verificationHash,
		})
	} else {
		config.DB.Model(&verifyRec).Updates(models.ResultVerification{
			GPA:  gpaStr,
			Hash: verificationHash,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"student_info": gin.H{
			"id":          student.ID,
			"name":        student.User.Name,
			"roll_number": student.RollNumber,
			"class":       student.Class.Name + " - " + student.Class.Section,
		},
		"results":           marks,
		"gpa":               fmt.Sprintf("%.2f", gpa),
		"status":            map[bool]string{true: "FAILED", false: "PASSED"}[hasFailed],
		"verification_hash": verificationHash,
		"verify_url":        fmt.Sprintf("http://localhost:8080/api/public/verify-result/%s", verificationHash),
	})
}



// VerifyResultPublicly verifies a result without requiring authentication.
func VerifyResultPublicly(c *gin.Context) {
	hash := c.Param("hash")

	var record models.ResultVerification
	err := config.DB.
		Preload("Student.User").
		Preload("Student.Class").
		Preload("Exam").
		Where("hash = ?", hash).
		First(&record).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"verified": false,
			"status":   "INVALID / FORGED RESULT",
			"message":  "No authentic academic record found with this cryptographic hash!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"verified":      true,
		"status":        "AUTHENTIC & VERIFIED",
		"issuer":        "Official School Academic Board",
		"student_name":  record.Student.User.Name,
		"roll_number":   record.Student.RollNumber,
		"class":         record.Student.Class.Name + " - " + record.Student.Class.Section,
		"exam_name":     record.Exam.Name,
		"official_gpa":  record.GPA,
		"certified_at":  record.CreatedAt.Format("2006-01-02 15:04:05"),
	})
}


// GetSubjects returns all subjects, optionally filtered by class.
func GetSubjects(c *gin.Context) {
	var subjects []models.Subject
	classID := c.Query("class_id")
	query := config.DB.Preload("Class")
	if classID != "" {
		query = query.Where("class_id = ?", classID)
	}
	if err := query.Find(&subjects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subjects"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"subjects": subjects})
}

// GetExams returns all exams.
func GetExams(c *gin.Context) {
	var exams []models.Exam
	if err := config.DB.Order("id desc").Find(&exams).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch exams"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"exams": exams})
}
