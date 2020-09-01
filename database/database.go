package database

import (
	"database/sql"
	"errors"
	"go-api/config"
	"log"
	"reflect"

	// PostgreSQL dialetc for opening connection
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db    *gorm.DB
	sqlDB *sql.DB
)

// Open tries to open a database connection
func Open() error {
	config := config.GetConfig()

	dsn := "user=" + config.Database.User + " host=" + config.Database.Host + " password=" + config.Database.Password + " port=" + config.Database.Port + " dbname=" + config.Database.Name + " sslmode=disable"

	sqlDB, err := sql.Open("postgres", dsn)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	log.Println("Connection stabilished")

	db, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if err != nil {
		log.Println("GORM failed to initialize database with current connection")
		return err
	}

	return nil
}

// NewTransaction start a new transaction
func NewTransaction() (*gorm.DB, error) {
	if db == nil {
		log.Println("Database not connected")
		return nil, errors.New("Database not connected")
	}

	tx := db.Session(&gorm.Session{SkipDefaultTransaction: true})

	return tx, nil
}

// Close tries to close the current connectin
func Close() {
	if sqlDB == nil {
		log.Println("Database not connected")
		return
	}
	sqlDB.Close()
}

// GetDBSession returns the current database opened session
func GetDBSession() *gorm.DB {
	if db == nil {
		log.Println("Database not connected.")
		return nil
	}
	return db
}

// ApplyMigrations run gorm AutoMigrate to the given model
func ApplyMigrations(model interface{}) {
	modelType := reflect.Indirect(reflect.ValueOf(model)).Type()

	err := db.AutoMigrate(model)

	if err != nil {
		log.Println("Failed applying migrations.")
		return
	}

	log.Printf("Migrations applied successfully for model %v\n", modelType)
}
