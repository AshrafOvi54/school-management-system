package models

import (
	"time"
	"gorm.io/gorm"
)

// Subject represents a subject offered to a class.
type Subject struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ClassID   uint           `gorm:"not null" json:"class_id"`
	Class     Class          `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Code      string         `gorm:"size:20;not null" json:"code"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Exam represents an assessment.
type Exam struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Term      string         `gorm:"size:50" json:"term"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Mark stores a student's score and grade for one subject in an exam.
type Mark struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ExamID        uint           `gorm:"not null;uniqueIndex:idx_exam_student_sub" json:"exam_id"`
	Exam          Exam           `gorm:"foreignKey:ExamID" json:"exam,omitempty"`
	StudentID     uint           `gorm:"not null;uniqueIndex:idx_exam_student_sub" json:"student_id"`
	Student       Student        `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	SubjectID     uint           `gorm:"not null;uniqueIndex:idx_exam_student_sub" json:"subject_id"`
	Subject       Subject        `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
	MarksObtained float64        `gorm:"type:decimal(5,2);not null" json:"marks_obtained"`
	Grade         string         `gorm:"size:5;not null" json:"grade"`
	GradePoint    float64        `gorm:"type:decimal(3,2);not null" json:"grade_point"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}