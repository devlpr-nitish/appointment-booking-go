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
func CreateAvailability(expertID uint, startTime, endTime string, isRecurring bool, dateStr string) (*models.AvailabilitySlot, error) {
	db := database.GetDB()

	// Parse date
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, errors.New("invalid date format, expected YYYY-MM-DD")
	}

	dayOfWeek := int(parsedDate.Weekday())

	// Validate that the expert exists
	var expert models.Expert
	if err := db.Where("id = ?", expertID).First(&expert).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("expert not found")
		}
		return nil, err
	}

	// Check for overlapping availability
	// Overlap check needs to consider:
	// 1. Same Expert
	// 2. Time overlap
	// 3. Condition:
	//    - If new is recurring: check against any recurring on same DayOfWeek, OR any specific date on that DayOfWeek
	//    - If new is specific: check against any recurring on same DayOfWeek, OR any specific date on same Date

	// Simplified overlap check: Check everything on this "DayOfWeek" that MIGHT conflict.
	// If new is recurring (Every Monday): Conflict if there's another Recurring Monday, OR a Specific Monday (e.g. Nov 11)??
	// Actually, if I have "Every Monday 9-5", adding "Nov 11 10-11" is redundant/overlap.
	// If I have "Nov 11 9-5", adding "Every Monday 9-5" is overlap.

	// Query finds any slot that overlaps in time AND:
	// (Existing.IsRecurring = true AND Existing.DayOfWeek = New.DayOfWeek) OR
	// (Existing.IsRecurring = false AND Existing.Date = New.Date) OR
	// (New.IsRecurring = true AND Existing.IsRecurring = false AND Existing.DayOfWeek = New.DayOfWeek) OR
	// (New.IsRecurring = false AND Existing.IsRecurring = true AND Existing.DayOfWeek = New.DayOfWeek)

	// Basically: If DayOfWeek matches AND (Both Recurring OR Both Specific Same Date OR Mixed)
	// Mixed is tricky. "Every Monday" vs "Monday Nov 11". They overlap in reality on Nov 11.
	// So we should just check: overlapping time on the SAME EFFECTIVE DAY.

	query := db.Where("expert_id = ? AND ((start_time < ? AND end_time > ?))", expertID, endTime, startTime)

	if isRecurring {
		// New is recurring: conflicts with ANY slot on this DayOfWeek (recurring or specific)
		query = query.Where("day_of_week = ?", dayOfWeek)
	} else {
		// New is specific: conflicts with Recurring slots on this DayOfWeek OR Specific slot on this Date
		query = query.Where("(is_recurring = ? AND day_of_week = ?) OR (is_recurring = ? AND date = ?)",
			true, dayOfWeek, false, parsedDate)
	}

	var existingSlot models.AvailabilitySlot
	if err := query.First(&existingSlot).Error; err == nil {
		return nil, errors.New("availability slot overlaps with existing slot")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	availability := models.AvailabilitySlot{
		ExpertID:    expertID,
		DayOfWeek:   dayOfWeek,
		StartTime:   startTime,
		EndTime:     endTime,
		IsRecurring: isRecurring,
	}

	if !isRecurring {
		availability.Date = &parsedDate
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
func UpdateAvailability(id, expertID uint, startTime, endTime string, isRecurring bool, dateStr string) (*models.AvailabilitySlot, error) {
	db := database.GetDB()

	// Parse date
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, errors.New("invalid date format, expected YYYY-MM-DD")
	}

	dayOfWeek := int(parsedDate.Weekday())

	var availability models.AvailabilitySlot
	if err := db.Where("id = ? AND expert_id = ?", id, expertID).First(&availability).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("availability slot not found")
		}
		return nil, err
	}

	// Check for overlapping availability (excluding current)
	query := db.Where("expert_id = ? AND id != ? AND ((start_time < ? AND end_time > ?))", expertID, id, endTime, startTime)

	if isRecurring {
		query = query.Where("day_of_week = ?", dayOfWeek)
	} else {
		query = query.Where("(is_recurring = ? AND day_of_week = ?) OR (is_recurring = ? AND date = ?)",
			true, dayOfWeek, false, parsedDate)
	}

	var existingSlot models.AvailabilitySlot
	if err := query.First(&existingSlot).Error; err == nil {
		return nil, errors.New("availability slot overlaps with existing slot")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Update fields
	availability.DayOfWeek = dayOfWeek
	availability.StartTime = startTime
	availability.EndTime = endTime
	availability.IsRecurring = isRecurring

	if !isRecurring {
		availability.Date = &parsedDate
	} else {
		availability.Date = nil
	}

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
	// Includes Recurring slots for this DayOfWeek AND Specific slots for this Date
	var availabilitySlots []models.AvailabilitySlot
	if err := db.Where("expert_id = ? AND ((is_recurring = ? AND day_of_week = ?) OR (is_recurring = ? AND date = ?))",
		expertID, true, dayOfWeek, false, parsedDate).
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
