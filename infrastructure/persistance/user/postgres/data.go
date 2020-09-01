package postgres

import (
	"go-api/domain/entities/user"
	"log"

	"gorm.io/gorm"
)

// PGUser is a base structure
// that implements methods for query execution
type PGUser struct {
	DB *gorm.DB
}

// Add insert an user into the database
func (pg *PGUser) Add(in *user.User) (err error) {
	if err := pg.DB.Create(in).Scan(&in); err != nil {
		log.Println(err.Error)
	}
	return nil
}
