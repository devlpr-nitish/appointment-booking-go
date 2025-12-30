package services

import (
	"errors"
	"fmt"

	"github.com/devlpr-nitish/appointment-booking-go/internal/database"
	"github.com/devlpr-nitish/appointment-booking-go/internal/models"
	"gorm.io/gorm"
)

// CreateAvailability creates a new availability slot for an expert
func CreateAvailability(expertID uint, dayOfWeek int, startTime, endTime string) (*models.AvailabilitySlot, error) {
	db := database.GetDB()

	// Validate day of week
	if dayOfWeek < 0 || dayOfWeek > 6 {
		return nil, errors.New("day of week must be between 0 (Sunday) and 6 (Saturday)")
	}

	// Validate that the expert exists
	var expert models.Expert
	if err := db.Where("id = ?", expertID).First(&expert).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("expert not found")
		}
		return nil, err
	}

	// Check for overlapping availability on the same day
	var existingSlot models.AvailabilitySlot
	err := db.Where("expert_id = ? AND day_of_week = ? AND ((start_time <= ? AND end_time > ?) OR (start_time < ? AND end_time >= ?) OR (start_time >= ? AND end_time <= ?))",
		expertID, dayOfWeek, startTime, startTime, endTime, endTime, startTime, endTime).First(&existingSlot).Error

	if err == nil {
		return nil, errors.New("availability slot overlaps with existing slot")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	availability := models.AvailabilitySlot{
		ExpertID:  expertID,
		DayOfWeek: dayOfWeek,
		StartTime: startTime,
		EndTime:   endTime,
	}

	if err := db.Create(&availability).Error; err != nil {
		return nil, err
	}

	return &availability, nil
}

// GetAvailabilityByExpertID retrieves all availability slots for an expert
func GetAvailabilityByExpertID(expertID uint) ([]models.AvailabilitySlot, error) {
	db := database.GetDB()

	var availability []models.AvailabilitySlot
	if err := db.Where("expert_id = ?", expertID).Order("day_of_week ASC, start_time ASC").Find(&availability).Error; err != nil {
		return nil, err
	}

	return availability, nil
}

// UpdateAvailability updates an existing availability slot
func UpdateAvailability(id, expertID uint, dayOfWeek int, startTime, endTime string) (*models.AvailabilitySlot, error) {
	db := database.GetDB()

	// Validate day of week
	if dayOfWeek < 0 || dayOfWeek > 6 {
		return nil, errors.New("day of week must be between 0 (Sunday) and 6 (Saturday)")
	}

	var availability models.AvailabilitySlot
	if err := db.Where("id = ? AND expert_id = ?", id, expertID).First(&availability).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("availability slot not found")
		}
		return nil, err
	}

	// Check for overlapping availability on the same day (excluding current slot)
	var existingSlot models.AvailabilitySlot
	err := db.Where("expert_id = ? AND day_of_week = ? AND id != ? AND ((start_time <= ? AND end_time > ?) OR (start_time < ? AND end_time >= ?) OR (start_time >= ? AND end_time <= ?))",
		expertID, dayOfWeek, id, startTime, startTime, endTime, endTime, startTime, endTime).First(&existingSlot).Error

	if err == nil {
		return nil, errors.New("availability slot overlaps with existing slot")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Update fields
	availability.DayOfWeek = dayOfWeek
	availability.StartTime = startTime
	availability.EndTime = endTime

	if err := db.Save(&availability).Error; err != nil {
		return nil, err
	}

	return &availability, nil
}

// DeleteAvailability deletes an availability slot
func DeleteAvailability(id, expertID uint) error {
	db := database.GetDB()

	result := db.Where("id = ? AND expert_id = ?", id, expertID).Delete(&models.AvailabilitySlot{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("availability slot not found or you don't have permission to delete it")
	}

	return nil
}
