package services

import (
	"errors"
	"time"

	"github.com/devlpr-nitish/appointment-booking-go/internal/database"
	"github.com/devlpr-nitish/appointment-booking-go/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateExpertProfile(userID uint, bio, expertise string, hourlyRate float64) (*models.Expert, error) {
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
		Expertise:  expertise,
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
	if err := db.Preload("User").Preload("Category").Where("user_id = ?", userID).First(&expert).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("expert profile not found")
		}
		return nil, err
	}
	return &expert, nil
}

func UpdateExpertProfile(userID uint, bio, expertise string, hourlyRate float64, categoryID *uuid.UUID) (*models.Expert, error) {
	db := database.GetDB()
	var expert models.Expert

	if err := db.Where("user_id = ?", userID).First(&expert).Error; err != nil {
		return nil, errors.New("expert profile not found")
	}

	if bio != "" {
		expert.Bio = bio
	}
	if expertise != "" {
		expert.Expertise = expertise
	}
	if hourlyRate > 0 {
		expert.HourlyRate = hourlyRate
	}
	if categoryID != nil {
		expert.CategoryID = categoryID
	}

	if err := db.Save(&expert).Error; err != nil {
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
	if err := db.Preload("User").Offset(offset).Limit(limit).Find(&experts).Error; err != nil {
		return nil, 0, err
	}
	return experts, total, nil
}

func SearchExperts(query, category string) ([]models.Expert, error) {
	db := database.GetDB()
	var experts []models.Expert

	tx := db.Model(&models.Expert{}).
		Joins("LEFT JOIN users ON users.id = experts.user_id").
		Preload("User")

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
	if err := db.Preload("User").Where("id = ?", id).First(&expert).Error; err != nil {
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
