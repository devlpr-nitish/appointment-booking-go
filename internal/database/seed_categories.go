package database

import (
	"log"

	"github.com/devlpr-nitish/appointment-booking-go/internal/models"
	"gorm.io/gorm"
)

// SeedCategories seeds the database with initial categories
func SeedCategories(db *gorm.DB) error {
	categories := []models.Category{
		{Name: "Web Development"},
		{Name: "Mobile Development"},
		{Name: "UI/UX Design"},
		{Name: "Digital Marketing"},
		{Name: "Content Writing"},
		{Name: "Business Consulting"},
		{Name: "Financial Advisory"},
		{Name: "Legal Consulting"},
		{Name: "Career Coaching"},
		{Name: "Language Tutoring"},
	}

	for _, category := range categories {
		// Check if category already exists
		var existingCategory models.Category
		if err := db.Where("name = ?", category.Name).First(&existingCategory).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Category doesn't exist, create it
				if err := db.Create(&category).Error; err != nil {
					log.Printf("Error creating category %s: %v", category.Name, err)
					return err
				}
				log.Printf("Created category: %s", category.Name)
			} else {
				log.Printf("Error checking category %s: %v", category.Name, err)
				return err
			}
		} else {
			log.Printf("Category already exists: %s", category.Name)
		}
	}

	log.Println("Category seeding completed successfully")
	return nil
}
