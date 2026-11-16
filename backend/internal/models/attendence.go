package models

import (
	"time"
	"gorm.io/gorm"
)

type AttendanceStatus string

const (
	StatusPresent AttendanceStatus = "present"
	StatusAbsent  AttendanceStatus = "absent"
	StatusLate    AttendanceStatus = "late"
)

type Attendance struct {
	ID        uint             `gorm:"primaryKey" json:"id"`
	StudentID uint             `gorm:"not null;uniqueIndex:idx_student_date" json:"student_id"` // A student can have one attendance record per day.
	Student   Student          `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Date      string           `gorm:"type:date;not null;uniqueIndex:idx_student_date" json:"date"` // YYYY-MM-DD
	Status    AttendanceStatus `gorm:"type:varchar(10);not null" json:"status"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	DeletedAt gorm.DeletedAt   `gorm:"index" json:"-"`
}