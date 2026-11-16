package models

import (
	"time"
	"gorm.io/gorm"
)

// Class represents a school class and section.
type Class struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	Section   string         `gorm:"size:10;not null" json:"section"`
	RoomNo    string         `gorm:"size:20" json:"room_no"`
	Students  []Student      `gorm:"foreignKey:ClassID" json:"students,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Student links a user account to a class.
type Student struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"not null;unique" json:"user_id"`
	User          User           `gorm:"foreignKey:UserID" json:"user"`
	ClassID       uint           `gorm:"not null" json:"class_id"`
	Class         Class          `gorm:"foreignKey:ClassID" json:"class"`
	RollNumber    string         `gorm:"size:20;not null" json:"roll_number"`
	GuardianPhone string         `gorm:"size:20;not null" json:"guardian_phone"`
	BloodGroup    string         `gorm:"size:10" json:"blood_group"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}