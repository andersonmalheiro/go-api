package postgres

import (
	"fmt"
	"go-api/domain/entities/user"
	"go-api/oops"

	"gorm.io/gorm"
)

// PGUser is a base structure
// that implements methods for query execution
type PGUser struct {
	DB *gorm.DB
}

// Add insert an user into the database
func (pg *PGUser) Add(in *user.User) (err error) {
	fmt.Printf("\nPGUser in:  %+v\n", in)
	if err := pg.DB.Create(in).Scan(&in); err != nil {
		return oops.Err(err.Error)
	}
	return nil
}
