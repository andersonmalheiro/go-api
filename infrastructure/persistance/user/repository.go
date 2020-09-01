package user

import (
	"go-api/domain/entities/user"
	"go-api/infrastructure/persistance/user/postgres"

	"gorm.io/gorm"
)

// Repository is a base structure that
// implements IUser methods
type Repository struct{}

// Add is a function that manage the flow of user insertion into database
func (r *Repository) Add(in *user.User, db *gorm.DB) error {
	data := postgres.PGUser{DB: db}
	return data.Add(in)
}

// Update updates an user
func (r *Repository) Update(u *user.User) error {
	return nil
}

// Delete removes an user
func (r *Repository) Delete(id int64) error {
	return nil
}

// Get returns an user by his ID
func (r *Repository) Get(u *user.User) error {
	return nil
}

// GetAll list all users
func (r *Repository) GetAll(interface{}) error {
	return nil
}
