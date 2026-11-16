package main

import (
	"log"
	"net/http"
	"os"

	"school-management-backend/internal/config"
	"school-management-backend/internal/handlers"
	"school-management-backend/internal/middleware"
	"school-management-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: .env file not found, reading system env")
	}

	config.ConnectDB()

	err = config.DB.AutoMigrate(
		&models.User{},
		&models.Class{},
		&models.Student{},
		&models.Attendance{},
		&models.Subject{},
		&models.Exam{},
		&models.Mark{},
		&models.Notice{},
		&models.ResultVerification{},
	)
	if err != nil {
		log.Fatal("AutoMigrate Failed: ", err)
	}
	log.Println(" All Tables (including Notices, ResultVerification) Migrated Successfully!")

	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	// Public routes do not require authentication.
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Backend is running!"})
	})
	// Public result verification endpoint.
	router.GET("/api/public/verify-result/:hash", handlers.VerifyResultPublicly)

	authRoutes := router.Group("/api/auth")
	{
		authRoutes.POST("/register", handlers.Register)
		authRoutes.POST("/login", handlers.Login)
	}

	// Protected routes require a valid token.
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/classes", handlers.GetClasses)
		protected.GET("/students", handlers.GetStudents)
		protected.GET("/teachers", handlers.GetTeachers)
		protected.GET("/subjects", handlers.GetSubjects)
		protected.GET("/exams", handlers.GetExams)
		protected.GET("/notices", handlers.GetNotices)
		protected.GET("/attendance/my-report", handlers.GetMyAttendanceReport)
		protected.GET("/reports/student/:student_id/exam/:exam_id", handlers.GetStudentReportCard)

		// Routes available to teachers and administrators.
		staffRoutes := protected.Group("/")
		staffRoutes.Use(middleware.RoleGuard("admin", "teacher"))
		{
			staffRoutes.POST("/attendance/mark", handlers.MarkAttendance)
			staffRoutes.GET("/attendance/class/:class_id", handlers.GetClassAttendance)
			staffRoutes.POST("/marks/submit", handlers.SubmitMarks)
		}

		// Routes available only to administrators.
		adminRoutes := protected.Group("/admin")
		adminRoutes.Use(middleware.RoleGuard("admin"))
		{
			adminRoutes.POST("/classes", handlers.CreateClass)
			adminRoutes.POST("/students", handlers.CreateStudent)
			adminRoutes.POST("/teachers", handlers.CreateTeacher)
			adminRoutes.POST("/subjects", handlers.CreateSubject)
			adminRoutes.POST("/exams", handlers.CreateExam)
			adminRoutes.POST("/notices", handlers.CreateNotice)
			adminRoutes.GET("/dashboard-stats", handlers.GetDashboardStats)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	router.Run(":" + port)
}
