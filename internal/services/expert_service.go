package services

import (
	"errors"
	"time"

	"github.com/devlpr-nitish/appointment-booking-go/internal/database"
	"github.com/devlpr-nitish/appointment-booking-go/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateExpertProfile(userID uint, bio string, hourlyRate float64) (*models.Expert, error) {
	db := database.GetDB()

	// Check if expert profile already exists
	var existingExpert models.Expert
	if err := db.Where("user_id = ?", userID).First(&existingExpert).Error; err == nil {
		return nil, errors.New("expert profile already exists for this user")
	}

	tx := db.Begin()

	expert := models.Expert{
		UserID:     userID,
		Bio:        bio,
		HourlyRate: hourlyRate,
		IsVerified: false,
	}

	if err := tx.Create(&expert).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Update user role to expert
	if err := tx.Model(&models.User{}).Where("id = ?", userID).Update("role", models.RoleExpert).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &expert, nil
}

func GetExpertProfile(userID uint) (*models.Expert, error) {
	db := database.GetDB()
	var expert models.Expert
	if err := db.
		Preload("User").
		Preload("Category").
		Preload("Categories").
		Where("user_id = ?", userID).
		First(&expert).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("expert profile not found")
		}
		return nil, err
	}
	return &expert, nil
}

// UpdateExpertProfile updates the expert profile including multiple categories.
func UpdateExpertProfile(userID uint, bio string, hourlyRate float64, categoryID *uuid.UUID, categoryIDs []uuid.UUID) (*models.Expert, error) {
	db := database.GetDB()
	var expert models.Expert

	if err := db.Preload("Categories").Where("user_id = ?", userID).First(&expert).Error; err != nil {
		return nil, errors.New("expert profile not found")
	}

	if bio != "" {
		expert.Bio = bio
	}
	if hourlyRate > 0 {
		expert.HourlyRate = hourlyRate
	}

	// Handle multi-category update
	if categoryIDs != nil {
		// Build slice of Category objects for association
		var newCategories []models.Category
		for _, cid := range categoryIDs {
			newCategories = append(newCategories, models.Category{ID: cid})
		}
		// Replace all associations (GORM will sync the join table)
		if err := db.Model(&expert).Association("Categories").Replace(newCategories); err != nil {
			return nil, err
		}
		// Also update the legacy single-category field (first one for backward compat)
		if len(categoryIDs) > 0 {
			first := categoryIDs[0]
			expert.CategoryID = &first
		} else {
			expert.CategoryID = nil
		}
	} else if categoryID != nil {
		// Fallback: single-category update
		expert.CategoryID = categoryID
		// Also add to the multi-category association if not already present
		cat := models.Category{ID: *categoryID}
		if err := db.Model(&expert).Association("Categories").Append(&cat); err != nil {
			return nil, err
		}
	}

	if err := db.Save(&expert).Error; err != nil {
		return nil, err
	}

	// Reload with all associations
	if err := db.Preload("User").Preload("Category").Preload("Categories").First(&expert, expert.ID).Error; err != nil {
		return nil, err
	}

	return &expert, nil
}

func GetExperts(page, limit int) ([]models.Expert, int64, error) {
	db := database.GetDB()
	var experts []models.Expert
	var total int64

	if err := db.Model(&models.Expert{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := db.Preload("User").Preload("Categories").Offset(offset).Limit(limit).Find(&experts).Error; err != nil {
		return nil, 0, err
	}
	return experts, total, nil
}

func SearchExperts(query, category string) ([]models.Expert, error) {
	db := database.GetDB()
	var experts []models.Expert

	tx := db.Model(&models.Expert{}).
		Joins("LEFT JOIN users ON users.id = experts.user_id").
		Preload("User").
		Preload("Categories")

	if category != "" && category != "all" {
		tx = tx.Where("experts.expertise = ?", category)
	}

	if query != "" {
		searchQuery := "%" + query + "%"
		tx = tx.Where(
			"experts.expertise ILIKE ? OR users.name ILIKE ?",
			searchQuery, searchQuery,
		)
	}

	if err := tx.Find(&experts).Error; err != nil {
		return nil, err
	}

	return experts, nil
}

func GetExpertById(id uint) (*models.Expert, error) {
	db := database.GetDB()
	var expert models.Expert
	if err := db.Preload("User").Preload("Categories").Where("id = ?", id).First(&expert).Error; err != nil {
		return nil, err
	}
	return &expert, nil
}

func GetExpertStats(expertID uint) (*models.ExpertStats, error) {
	db := database.GetDB()
	var stats models.ExpertStats

	// 1. Total Earnings (Sum of TotalPrice for completed bookings)
	var totalEarnings float64
	if err := db.Model(&models.Booking{}).
		Where("expert_id = ? AND status = ?", expertID, models.BookingStatusCompleted).
		Select("COALESCE(SUM(total_price), 0)").Scan(&totalEarnings).Error; err != nil {
		return nil, err
	}
	stats.TotalEarnings = totalEarnings

	// 2. Upcoming Sessions (Count of confirmed bookings with start_time in future)
	var upcomingSessions int64
	today := time.Now().Format("2006-01-02")
	if err := db.Model(&models.Booking{}).
		Where("expert_id = ? AND status = ? AND booking_date >= ?", expertID, models.BookingStatusConfirmed, today).
		Count(&upcomingSessions).Error; err != nil {
		return nil, err
	}
	stats.UpcomingSessions = int(upcomingSessions)

	// 3. Completed Sessions
	var completedSessions int64
	if err := db.Model(&models.Booking{}).
		Where("expert_id = ? AND status = ?", expertID, models.BookingStatusCompleted).
		Count(&completedSessions).Error; err != nil {
		return nil, err
	}
	stats.CompletedSessions = int(completedSessions)

	// 4. Pending Requests
	var pendingRequests int64
	if err := db.Model(&models.Booking{}).
		Where("expert_id = ? AND status = ?", expertID, models.BookingStatusPending).
		Count(&pendingRequests).Error; err != nil {
		return nil, err
	}
	stats.PendingRequests = int(pendingRequests)

	// 5. Average Rating (Placeholder)
	stats.AverageRating = 0.0

	return &stats, nil
}
