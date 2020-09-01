package main

import (
	"fmt"
	app "go-api/application/entities/user"
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
		log.Println("Error when loading config")
		return
	}

	err = database.Open()

	if err != nil {
		log.Println("Error when stabilishing connection to the database")
		return
	}

	defer database.Close()

	log.Println("Applying migrations...")
	fmt.Println()

	for _, model := range models {
		database.ApplyMigrations(model)
	}

	fmt.Println()
	log.Println("Migrations finished")

	// Testing user insertion
	// Remove later
	var user app.INUser

	id, err := app.Add(&user)

	if err != nil {
		return
	}

	fmt.Printf("user id %v\n", id)
	// Testing user insertion
	// Remove later
}
