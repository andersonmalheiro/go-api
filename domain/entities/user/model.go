package user

import (
	"go-api/domain/entities/class"

	"gorm.io/gorm"
)

// User struct defines the fields of user table1
type User struct {
	gorm.Model
	Name          string `gorm:"not null"`
	Email         string `gorm:"unique; not null"`
	Password      string `gorm:"not null"`
	BirthDate     string
	AvatarURL     string
	ContactNumber string
	Bio           string
	Classes       []class.Class `gorm:"foreignKey:TeacherID"`
}
