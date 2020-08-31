package database

import (
	"database/sql"
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
		log.Fatal(err.Error())
		return err
	}

	log.Println("Connection stabilished")

	db, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	if err != nil {
		log.Fatal("GORM failed to initialize database with current connection")
		return err
	}

	return nil
}

// Close tries to close the current connectin
func Close() {
	if sqlDB == nil {
		log.Fatal("Database not connected")
	}

	sqlDB.Close()
}

// GetDBSession returns the current database opened session
func GetDBSession() *gorm.DB {
	if db == nil {
		log.Fatal("Database not connected.")
		return nil
	}

	return db
}

// ApplyMigrations run gorm AutoMigrate to the given model
func ApplyMigrations(model interface{}) {
	modelType := reflect.Indirect(reflect.ValueOf(model)).Type()

	err := db.AutoMigrate(model)

	if err != nil {
		log.Fatal("Failed applying migrations.")
	} else {
		log.Printf("Migrations applied successfully for model %v\n", modelType)
	}
}
