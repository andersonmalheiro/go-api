package user

import (
	"fmt"
	"go-api/database"
	domain "go-api/domain/entities/user"
	repository "go-api/infrastructure/persistance/user"
	"go-api/oops"
	"go-api/utils"
	"log"
)

// Add do the business logic of inserting an user into the database
func Add(in *INUser) (id uint, err error) {
	var repo domain.IUser = &repository.Repository{}

	fmt.Printf("\nBusiness in:  %+v\n", in)

	tx, err := database.NewTransaction()

	if err != nil {
		return id, oops.Wrap(err, "Error when initializing transaction.")
	}

	defer tx.Rollback()

	data := &domain.User{}

	if err = utils.ConvertStruct(in, data); err != nil {
		log.Println(err)
		return id, oops.Wrap(err, "Error when converting struct.")
	}

	fmt.Printf("\nBusiness data:  %+v\n", data)

	if err = repo.Add(data, tx); err != nil {
		return id, oops.Wrap(err, "Error when adding new user.")
	}

	tx.Commit()

	return *data.ID, nil
}
