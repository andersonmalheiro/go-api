package user

import (
	"time"
)

// INUser models a user for insertion
type INUser struct {
	Name          string `json:"name" binding:"required"`
	BirthDate     string `json:"birth_date"`
	Email         string `json:"email" binding:"required"`
	Password      string `json:"password" binding:"required"`
	AvatarURL     string `json:"avatar_url"`
	Bio           string `json:"bio"`
	ContactNumber string `json:"contact_number"`
}

// OUTUser models a user for retrieval
type OUTUser struct {
	ID            *uint      `json:"id,omitempty"`
	Name          *string    `json:"name,omitempty"`
	BirthDate     *string    `json:"birth_date,omitempty"`
	Email         *string    `json:"email,omitempty"`
	Password      *string    `json:"password,omitempty"`
	AvatarURL     *string    `json:"avatar_url,omitempty"`
	Bio           *string    `json:"bio,omitempty"`
	ContactNumber *string    `json:"contact_number,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}

// OUTList models a list of users
type OUTList struct {
	Data []OUTUser
}
