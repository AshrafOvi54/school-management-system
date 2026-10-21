package main

import (
	"log"
	"net/http"
	"os"

	"school-management-backend/internal/config"
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

	err = config.DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("AutoMigrate Failed: ", err)
	}
	log.Println("✅ Database Migrated Successfully!")

	
	router := gin.Default()

	
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "School Management System Backend is running perfectly!",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	router.Run(":" + port)
}