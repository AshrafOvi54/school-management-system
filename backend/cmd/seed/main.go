package main

import (
	"fmt"
	"time"

	"school-management-backend/internal/config"
	"school-management-backend/internal/models"
	"school-management-backend/internal/utils"

	"github.com/joho/godotenv"
	"gorm.io/gorm/clause"
)

func main() {
	godotenv.Load()
	config.ConnectDB()

	fmt.Println("🌱 Seeding realistic academic demo data...")

	adminHash, _ := utils.HashPassword("secretpassword")
	teacherHash, _ := utils.HashPassword("teacherpass123")
	studentHash, _ := utils.HashPassword("studentpass123")

	admin := models.User{
		Name:     "Headmaster Rahim Uddin",
		Email:    "admin@school.com",
		Password: adminHash,
		Role:     models.RoleAdmin,
	}
	config.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&admin)

	teachers := []models.User{
		{Name: "Professor Kamal Uddin", Email: "kamal@school.com", Password: teacherHash, Role: models.RoleTeacher},
		{Name: "Ms. Farhana Yasmin", Email: "farhana@school.com", Password: teacherHash, Role: models.RoleTeacher},
	}
	for _, t := range teachers {
		config.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&t)
	}

	classes := []models.Class{
		{Name: "Class 10", Section: "A", RoomNo: "Room 101"},
		{Name: "Class 10", Section: "B", RoomNo: "Room 102"},
		{Name: "Class 9", Section: "A", RoomNo: "Room 201"},
	}
	for i := range classes {
		config.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&classes[i])
	}

	subjects := []models.Subject{
		{ClassID: 1, Name: "General Mathematics", Code: "MATH101"},
		{ClassID: 1, Name: "English for Today", Code: "ENG101"},
		{ClassID: 1, Name: "Physics", Code: "PHY101"},
		{ClassID: 1, Name: "Bangla First Paper", Code: "BAN101"},
		{ClassID: 1, Name: "Information & Comm. Tech", Code: "ICT101"},
	}
	for _, s := range subjects {
		config.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&s)
	}

	exam := models.Exam{Name: "Annual Final Examination 2026", Term: "Final"}
	config.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&exam)

	studentsData := []struct {
		name, email, roll, phone, blood string
		classID                         uint
	}{
		{"Karim Ahmed", "karim@school.com", "101", "01711223344", "A+", 1},
		{"Shakib Al Hasan", "shakib@school.com", "102", "01811223344", "B+", 1},
		{"Nusrat Jahan", "nusrat@school.com", "103", "01911223344", "O+", 1},
		{"Tanvir Hossain", "tanvir@school.com", "201", "01611223344", "AB+", 2},
		{"Ayesha Siddiqua", "ayesha@school.com", "301", "01511223344", "O+", 3},
	}

	for _, s := range studentsData {
		u := models.User{Name: s.name, Email: s.email, Password: studentHash, Role: models.RoleStudent}
		if err := config.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&u).Error; err == nil {
			var existingUser models.User
			config.DB.Where("email = ?", s.email).First(&existingUser)

			st := models.Student{
				UserID:        existingUser.ID,
				ClassID:       s.classID,
				RollNumber:    s.roll,
				GuardianPhone: s.phone,
				BloodGroup:    s.blood,
			}
			config.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&st)
		}
	}

	marksData := []struct {
		subID uint
		marks float64
	}{
		{1, 88.0},
		{2, 74.5},
		{3, 68.0},
		{4, 82.0},
		{5, 92.5},
	}

	totalPoints := 0.0
	for _, m := range marksData {
		grade, pt := utils.CalculateGrade(m.marks)
		totalPoints += pt
		mark := models.Mark{
			ExamID:        1,
			StudentID:     1,
			SubjectID:     m.subID,
			MarksObtained: m.marks,
			Grade:         grade,
			GradePoint:    pt,
		}
		config.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "exam_id"}, {Name: "student_id"}, {Name: "subject_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"marks_obtained", "grade", "grade_point"}),
		}).Create(&mark)
	}

	finalGPA := totalPoints / float64(len(marksData)) // 4.50
	hash := utils.GenerateVerificationHash(1, 1, finalGPA)
	gpaStr := fmt.Sprintf("%.2f", finalGPA)

	verifyRec := models.ResultVerification{
		StudentID: 1,
		ExamID:    1,
		GPA:       gpaStr,
		Hash:      hash,
	}
	config.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"gpa", "hash"}),
	}).Create(&verifyRec)

	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	dayBefore := time.Now().AddDate(0, 0, -2).Format("2006-01-02")

	attRecords := []models.Attendance{
		{StudentID: 1, Date: today, Status: models.StatusPresent},
		{StudentID: 2, Date: today, Status: models.StatusPresent},
		{StudentID: 3, Date: today, Status: models.StatusAbsent},
		{StudentID: 1, Date: yesterday, Status: models.StatusPresent},
		{StudentID: 2, Date: yesterday, Status: models.StatusAbsent},
		{StudentID: 1, Date: dayBefore, Status: models.StatusPresent},
	}
	for _, a := range attRecords {
		config.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&a)
	}

	notices := []models.Notice{
		{Title: "Annual Sports Competition 2026", Content: "The annual sports competition will be held next Sunday at the main school stadium.", TargetRole: "all"},
		{Title: "Admit Card Distribution", Content: "All students are requested to collect their annual exam admit cards from their class teachers.", TargetRole: "student"},
		{Title: "Monthly Faculty & Staff Meeting", Content: "Urgent meeting on upcoming semester syllabus review at 2:00 PM in the Conference Hall.", TargetRole: "teacher"},
	}
	for _, n := range notices {
		config.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&n)
	}

	fmt.Println("✅ All demo data seeded successfully!")
	fmt.Printf("🎯 Karim's Official GPA: %s\n", gpaStr)
	fmt.Printf("🔐 Karim's Verification Hash: %s\n", hash)
}
