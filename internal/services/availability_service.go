package services

import (
	"errors"
	"fmt"
	"time"

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

// TimeSlot represents a bookable time slot
type TimeSlot struct {
	Time      string `json:"time"`
	Available bool   `json:"available"`
	ID        uint   `json:"id"`
}

// GetAvailableSlots generates available time slots for an expert on a specific date
func GetAvailableSlots(expertID uint, date string) ([]TimeSlot, error) {
	db := database.GetDB()

	// Parse the date to get day of week
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, errors.New("invalid date format, expected YYYY-MM-DD")
	}

	dayOfWeek := int(parsedDate.Weekday())

	// Get all availability slots for this expert on this day of week
	var availabilitySlots []models.AvailabilitySlot
	if err := db.Where("expert_id = ? AND day_of_week = ?", expertID, dayOfWeek).
		Order("start_time ASC").
		Find(&availabilitySlots).Error; err != nil {
		return nil, err
	}

	// If no availability slots, return empty array
	if len(availabilitySlots) == 0 {
		return []TimeSlot{}, nil
	}

	// Get all bookings for this expert on this date
	var bookings []models.Booking
	if err := db.Preload("Slot").
		Where("expert_id = ? AND status != ?", expertID, models.BookingStatusCancelled).
		Find(&bookings).Error; err != nil {
		return nil, err
	}

	// Create a map of booked times for quick lookup
	bookedTimes := make(map[string]bool)
	for _, booking := range bookings {
		// For each booking, mark the time slot as booked
		bookedTimes[booking.Slot.StartTime] = true
	}

	// Generate time slots
	var timeSlots []TimeSlot
	slotDuration := 30 // 30 minutes per slot

	for _, availSlot := range availabilitySlots {
		// Parse start and end times
		startTime, err := time.Parse("15:04", availSlot.StartTime)
		if err != nil {
			continue
		}
		endTime, err := time.Parse("15:04", availSlot.EndTime)
		if err != nil {
			continue
		}

		// Generate slots in 30-minute intervals
		currentTime := startTime
		for currentTime.Before(endTime) {
			timeStr := currentTime.Format("15:04")

			// Check if this time is booked
			isAvailable := !bookedTimes[timeStr]

			timeSlots = append(timeSlots, TimeSlot{
				Time:      timeStr,
				Available: isAvailable,
				ID:        availSlot.ID,
			})

			// Move to next slot
			currentTime = currentTime.Add(time.Duration(slotDuration) * time.Minute)
		}
	}

	return timeSlots, nil
}
