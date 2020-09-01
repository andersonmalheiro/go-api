package class

import (
	"time"

	"gorm.io/gorm"
)

// Class struct defines the fields of class table
type Class struct {
	Name      *string         `gorm:"not null" conversor:"name"`
	Price     *int64          `gorm:"not null" conversor:"price"`
	TeacherID *uint           `conversor:"teacher_id"`
	ID        *uint           `gorm:"primaryKey" conversor:"id"`
	CreatedAt *time.Time      `conversor:"created_at"`
	UpdatedAt *time.Time      `conversor:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" conversor:"deleted_at"`
	// Schedules []Schedule `gorm:"foreignKey:ClassID"`
}

// Schedule defines the field of a class schedule
type Schedule struct {
	Date      *string         `gorm:"not null" conversor:"date"`
	Start     *string         `gorm:"not null" conversor:"start"`
	End       *string         `gorm:"not null" conversor:"end"`
	ClassID   *uint           `conversor:"class_id"`
	ID        *uint           `gorm:"primaryKey" conversor:"id"`
	CreatedAt *time.Time      `conversor:"created_at"`
	UpdatedAt *time.Time      `conversor:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" conversor:"deleted_at"`
}
