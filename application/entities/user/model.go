package user

import (
	"time"
)

// INUser models a user for insertion
type INUser struct {
	Name          *string `json:"name" binding:"required" conversor:"name"`
	BirthDate     *string `json:"birth_date" conversor:"birth_date"`
	Email         *string `json:"email" binding:"required" conversor:"email"`
	Password      *string `json:"password" binding:"required" conversor:"password"`
	AvatarURL     *string `json:"avatar_url" conversor:"avatar_url"`
	Bio           *string `json:"bio" conversor:"bio"`
	ContactNumber *string `json:"contact_number" conversor:"contact_number"`
}

// OUTUser models a user for retrieval
type OUTUser struct {
	ID            *uint      `json:"id,omitempty" conversor:"id"`
	Name          *string    `json:"name,omitempty" conversor:"name"`
	BirthDate     *string    `json:"birth_date,omitempty" conversor:"birth_date"`
	Email         *string    `json:"email,omitempty" conversor:"email"`
	Password      *string    `json:"password,omitempty" conversor:"password"`
	AvatarURL     *string    `json:"avatar_url,omitempty" conversor:"avatar_url"`
	Bio           *string    `json:"bio,omitempty" conversor:"bio"`
	ContactNumber *string    `json:"contact_number,omitempty" conversor:"contact_number"`
	CreatedAt     *time.Time `json:"created_at,omitempty" conversor:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty" conversor:"updated_at"`
}

// OUTList models a list of users
type OUTList struct {
	Data []OUTUser
}
