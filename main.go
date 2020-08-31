package main

import (
	"fmt"
	"go-api/config"
	"go-api/database"
	"log"
)

func main() {
	err := config.LoadConfig()

	if err != nil {
		log.Fatal("Error when loading config")
	}

	err = database.Open()

	if err != nil {
		log.Fatal("Error when stabilishing connection to the database")
	}

	db := database.GetDBSession()

	if db == nil {
		return
	}

	fmt.Println(db.Name())
}
