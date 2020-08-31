package database

import (
	"database/sql"
	"go-api/config"
	"log"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Open() error {
	config := config.GetConfig()

	dsn := "user=" + config.Database.User + " host=" + config.Database.Host + " password=" + config.Database.Password + " port=" + config.Database.Port + " dbname=" + config.Database.Name + " sslmode=disable"

	sqlDB, err := sql.Open("postgres", dsn)

	defer sqlDB.Close()

	if err != nil {
		log.Fatal(err.Error())
		return err
	} else {
		log.Println("Connection stabilished")
	}

	db, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	if err != nil {
		log.Fatal("GORM failed to initialize database with current connection")
		return err
	}

	return nil
}

func GetDBSession() *gorm.DB {
	if db == nil {
		log.Fatal("Database not connected.")
		return nil
	}

	return db
}
