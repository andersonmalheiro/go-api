package main

import (
	"fmt"
	"go-api/config"
	"go-api/database"
	"go-api/domain/entities/class"
	"go-api/domain/entities/user"
	"log"
)

var models = []interface{}{
	user.User{},
	class.Class{},
	class.Schedule{},
}

func main() {
	err := config.LoadConfig()

	if err != nil {
		log.Fatal("Error when loading config")
	}

	err = database.Open()

	if err != nil {
		log.Fatal("Error when stabilishing connection to the database")
	}

	log.Println("Applying migrations...")
	fmt.Println()

	for _, model := range models {
		database.ApplyMigrations(model)
	}

	fmt.Println()
	log.Println("Migrations finished")

	defer database.Close()
}
