package models

import (
	"time"
	"gorm.io/gorm"
)

// Notice represents an announcement targeted to one or more roles.
type Notice struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Title      string         `gorm:"size:200;not null" json:"title"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	TargetRole string         `gorm:"size:20;not null;default:'all'" json:"target_role"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// ResultVerification stores a signed result record for public verification.
type ResultVerification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	StudentID uint      `gorm:"not null" json:"student_id"`
	Student   Student   `gorm:"foreignKey:StudentID" json:"student"`
	ExamID    uint      `gorm:"not null" json:"exam_id"`
	Exam      Exam      `gorm:"foreignKey:ExamID" json:"exam"`
	GPA       string    `gorm:"size:10;not null" json:"gpa"`
	Hash      string    `gorm:"size:64;uniqueIndex;not null" json:"hash"`
	CreatedAt time.Time `json:"created_at"`
}