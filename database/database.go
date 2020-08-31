package database

import (
	"go-api/config"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Test struct {
	Field string
}

var db *gorm.DB

func Open() {
	config := config.GetConfig()
	dsn := "user=" + config.Database.User + " password=" + config.Database.Password + " dbname=" + config.Database.Name + " port=" + config.Database.Port
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	log.Println(dsn)

	if err != nil {
		log.Fatal("Error when opening connection to the database")
	}

	db.AutoMigrate(&Test{})
}
