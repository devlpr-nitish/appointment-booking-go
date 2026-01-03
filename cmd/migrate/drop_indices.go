package main

import (
	"log"

	"github.com/devlpr-nitish/appointment-booking-go/internal/config"
	"github.com/devlpr-nitish/appointment-booking-go/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.LoadConfig()
	db := database.Connect(cfg)

	// Drop indices if they exist
	indices := []string{"idx_bookings_expert_id", "idx_bookings_user_id", "idx_bookings_slot_id"}

	for _, idx := range indices {
		log.Printf("Dropping index: %s", idx)
		if err := db.Exec("DROP INDEX IF EXISTS " + idx).Error; err != nil {
			log.Printf("Failed to drop index %s: %v", idx, err)
		} else {
			log.Printf("Successfully dropped index %s", idx)
		}
	}

	log.Println("Migration completed")
}
