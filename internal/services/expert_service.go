package services

import (
	"errors"

	"github.com/devlpr-nitish/appointment-booking-go/internal/database"
	"github.com/devlpr-nitish/appointment-booking-go/internal/models"
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
	if err := db.Preload("User").Where("user_id = ?", userID).First(&expert).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("expert profile not found")
		}
		return nil, err
	}
	return &expert, nil
}

func UpdateExpertProfile(userID uint, bio, expertise string, hourlyRate float64) (*models.Expert, error) {
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

func GetExpertByCatergoryName(categoryName string) ([]models.Expert, error) {
	db := database.GetDB()
	var experts []models.Expert
	if err := db.Preload("User").Where("expertise = ?", categoryName).Find(&experts).Error; err != nil {
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