package main

import (
	"go-api/config"
	"go-api/database"
	"log"
)

func main() {
	err := config.LoadConfig()

	if err != nil {
		log.Fatal("Error when loading config")
	}

	database.Open()
}
