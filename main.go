package main

import (
	"fmt"
	"go-api/config"
	"go-api/database"
	"go-api/domain/entities/class"
	"go-api/domain/entities/user"
	userRoutes "go-api/interfaces/entities/user"
	"log"

	"github.com/gin-gonic/gin"
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

	r := gin.New()

	r.Use(gin.Logger())

	v1 := r.Group("v1")

	userRoutes.Router(v1.Group("/users"))

	r.Run()
}
