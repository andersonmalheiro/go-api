package class

import "gorm.io/gorm"

// Class struct defines the fields of class table
type Class struct {
	gorm.Model
	Name      string `gorm:"not null"`
	Price     int64  `gorm:"not null"`
	Schedules []Schedule
	TeacherID uint
}

// Schedule defines the field of a class schedule
type Schedule struct {
	gorm.Model
	Date    string `gorm:"not null"`
	Start   string `gorm:"not null"`
	End     string `gorm:"not null"`
	ClassID uint
}
