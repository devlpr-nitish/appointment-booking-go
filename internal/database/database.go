package database

import (
	"log"

	"github.com/devlpr-nitish/appointment-booking-go/internal/config"
	"github.com/devlpr-nitish/appointment-booking-go/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(cfg *config.Config) *gorm.DB {
	var dialector gorm.Dialector

	switch cfg.DBDriver {
	case "postgres":
		if cfg.DBUrl == "" {
			log.Fatal("Database url not found for postgres driver")
		}
		dialector = postgres.Open(cfg.DBUrl)
	case "mysql":
		// Import "gorm.io/driver/mysql"
		// dialector = mysql.Open(cfg.DBUrl)
	case "sqlite":
		// Import "gorm.io/driver/sqlite"
		// dialector = sqlite.Open(cfg.DBUrl)
	default:
		log.Fatalf("Unsupported database driver: %s", cfg.DBDriver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to open DB connection: %v", err)
	}

	// AutoMigrate all models
	err = db.AutoMigrate(
		&models.User{},
		&models.Expert{},
		&models.AvailabilitySlot{},
		&models.Booking{},
		&models.Payment{},
		&models.Review{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate models: %v", err)
	}

	DB = db

	log.Println("Db connected successfully")

	return db
}

func GetDB() *gorm.DB {
	return DB
}
