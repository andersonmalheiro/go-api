package user

import (
	"fmt"
	"go-api/database"
	domain "go-api/domain/entities/user"
	repository "go-api/infrastructure/persistance/user"
	"go-api/oops"
	"go-api/utils"
)

// Add do the business logic of inserting an user into the database
func Add(in *INUser) (id uint, err error) {
	var repo domain.IUser = &repository.Repository{}

	fmt.Printf("\nBusiness in:  %+v\n", in)

	tx, err := database.NewTransaction()

	if err != nil {
		return id, oops.Wrap(err, "Erro when adding new user.")
	}

	defer tx.Rollback()

	data := &domain.User{}

	if err = utils.ConvertStruct(in, data); err != nil {
		return id, oops.Wrap(err, "Erro when adding new user.")
	}

	fmt.Printf("\nBusiness data:  %+v\n", data)

	if err = repo.Add(data, tx); err != nil {
		return id, oops.Wrap(err, "Erro when adding new user.")
	}

	tx.Commit()

	return data.ID, nil
}
