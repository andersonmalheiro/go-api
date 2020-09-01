package user

import (
	"errors"
	"go-api/database"
	domain "go-api/domain/entities/user"
	repository "go-api/infrastructure/persistance/user"
)

// Add do the business logic of inserting an user into the database
func Add(in *INUser) (id uint, err error) {
	var repo domain.IUser = &repository.Repository{}

	tx, err := database.NewTransaction()

	if err != nil {
		return id, err
	}

	defer tx.Rollback()

	data := &domain.User{
		Name:     "Chico",
		Email:    "chico@test.com",
		Password: "12345",
	}

	// if err = utils.ConvertStruct(in, data); err != nil {
	// 	return id, errors.New("Erro when adding new user.\n" + err.Error())
	// }

	if err = repo.Add(data, tx); err != nil {
		return id, errors.New("Error when creating new user.\n" + err.Error())
	}

	tx.Commit()

	return data.ID, nil
}
