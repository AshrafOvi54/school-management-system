package handlers

import (
	"net/http"

	"school-management-backend/internal/config"
	"school-management-backend/internal/models"
	"school-management-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type CreateClassInput struct {
	Name    string `json:"name" binding:"required"`
	Section string `json:"section" binding:"required"`
	RoomNo  string `json:"room_no"`
}

type CreateStudentInput struct {
	Name          string `json:"name" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	Password      string `json:"password" binding:"required,min=6"`
	ClassID       uint   `json:"class_id" binding:"required"`
	RollNumber    string `json:"roll_number" binding:"required"`
	GuardianPhone string `json:"guardian_phone" binding:"required"`
	BloodGroup    string `json:"blood_group"`
}

// CreateClass creates a class.
func CreateClass(c *gin.Context) {
	var input CreateClassInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	class := models.Class{
		Name:    input.Name,
		Section: input.Section,
		RoomNo:  input.RoomNo,
	}

	if err := config.DB.Create(&class).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create class"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Class created successfully!",
		"class":   class,
	})
}

// GetClasses returns all classes.
func GetClasses(c *gin.Context) {
	var classes []models.Class
	if err := config.DB.Find(&classes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch classes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"classes": classes})
}

// CreateStudent creates a student and their user account in one transaction.
func CreateStudent(c *gin.Context) {
	var input CreateStudentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var class models.Class
	if err := config.DB.First(&class, input.ClassID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Class not found with given class_id"})
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	tx := config.DB.Begin()

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
		Role:     models.RoleStudent,
	}
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists or failed to create user"})
		return
	}

	student := models.Student{
		UserID:        user.ID,
		ClassID:       input.ClassID,
		RollNumber:    input.RollNumber,
		GuardianPhone: input.GuardianPhone,
		BloodGroup:    input.BloodGroup,
	}
	if err := tx.Create(&student).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create student profile"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Student registered successfully!",
		"student": gin.H{
			"student_id":     student.ID,
			"name":           user.Name,
			"email":          user.Email,
			"roll_number":    student.RollNumber,
			"class":          class.Name + " - " + class.Section,
			"guardian_phone": student.GuardianPhone,
		},
	})
}

// GetStudents returns students with their user and class details.
func GetStudents(c *gin.Context) {
	var students []models.Student

	if err := config.DB.Preload("User").Preload("Class").Find(&students).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch students"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"students": students})
}
type CreateTeacherInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// CreateTeacher creates a teacher account.
func CreateTeacher(c *gin.Context) {
	var input CreateTeacherInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	teacher := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
		Role:     models.RoleTeacher,
	}

	if err := config.DB.Create(&teacher).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Teacher account created successfully!",
		"teacher": gin.H{
			"id":    teacher.ID,
			"name":  teacher.Name,
			"email": teacher.Email,
			"role":  teacher.Role,
		},
	})
}

// GetTeachers returns all teacher accounts.
func GetTeachers(c *gin.Context) {
	var teachers []models.User
	if err := config.DB.Where("role = ?", models.RoleTeacher).Find(&teachers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teachers"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"teachers": teachers})
}
