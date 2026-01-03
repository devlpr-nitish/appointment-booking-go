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

// TimeSlot represents a free time interval
type TimeSlot struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// GetAvailableSlots calculates available time ranges for an expert on a specific date
func GetAvailableSlots(expertID uint, date string) ([]TimeSlot, error) {
	db := database.GetDB()

	// Parse the date to get day of week
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, errors.New("invalid date format, expected YYYY-MM-DD")
	}

	dayOfWeek := int(parsedDate.Weekday())

	// Get all availability slots for this expert (Master Availability)
	var availabilitySlots []models.AvailabilitySlot
	if err := db.Where("expert_id = ? AND day_of_week = ?", expertID, dayOfWeek).
		Order("start_time ASC").
		Find(&availabilitySlots).Error; err != nil {
		return nil, err
	}

	if len(availabilitySlots) == 0 {
		return []TimeSlot{}, nil
	}

	// Get all confirmed bookings for this date
	var bookings []models.Booking
	if err := db.Where("expert_id = ? AND booking_date = ? AND status != ?", expertID, date, models.BookingStatusCancelled).
		Order("start_time ASC").
		Find(&bookings).Error; err != nil {
		return nil, err
	}

	var freeSlots []TimeSlot

	// Process each availability slot
	for _, avail := range availabilitySlots {
		// Start with the full availability slot as free
		currentRanges := []TimeSlot{{Start: avail.StartTime, End: avail.EndTime}}

		// Subtract each booking from the current ranges
		for _, booking := range bookings {
			var nextRanges []TimeSlot
			for _, r := range currentRanges {
				// Check for overlap
				// Booking: [bS, bE], Range: [rS, rE]
				// Overlap if bS < rE AND bE > rS
				if booking.StartTime < r.End && booking.EndTime > r.Start {
					// Overlap exists. Calculate remaining parts.

					// Part before booking?
					if booking.StartTime > r.Start {
						nextRanges = append(nextRanges, TimeSlot{Start: r.Start, End: booking.StartTime})
					}

					// Part after booking?
					if booking.EndTime < r.End {
						nextRanges = append(nextRanges, TimeSlot{Start: booking.EndTime, End: r.End})
					}
				} else {
					// No overlap, keep range
					nextRanges = append(nextRanges, r)
				}
			}
			currentRanges = nextRanges
		}
		freeSlots = append(freeSlots, currentRanges...)
	}

	return freeSlots, nil
}
