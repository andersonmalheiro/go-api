package user

import (
	"time"

	"gorm.io/gorm"
)

// User struct defines the fields of user table1
type User struct {
	Name          *string         `gorm:"not null" conversor:"name"`
	Email         *string         `gorm:"unique; not null" conversor:"email"`
	Password      *string         `gorm:"not null" conversor:"password"`
	BirthDate     *string         `conversor:"birth_date"`
	AvatarURL     *string         `conversor:"avatar_url"`
	ContactNumber *string         `conversor:"contact_number"`
	Bio           *string         `conversor:"bio"`
	ID            *uint           `gorm:"primaryKey" conversor:"id"`
	CreatedAt     *time.Time      `conversor:"created_at"`
	UpdatedAt     *time.Time      `conversor:"updated_at"`
	DeletedAt     *gorm.DeletedAt `gorm:"index" conversor:"deleted_at"`
	// Classes       []class.Class  `gorm:"foreignKey:TeacherID"`
}
