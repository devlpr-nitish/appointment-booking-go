package services

import (
	"errors"
	"time"

	"github.com/devlpr-nitish/appointment-booking-go/internal/database"
	"github.com/devlpr-nitish/appointment-booking-go/internal/models"
	"gorm.io/gorm"
)

func CreateBooking(userID, expertID uint, bookingDate, startTime, endTime string) (*models.Booking, error) {
	db := database.GetDB()

	// 1. Validate Time Format & Logic
	start, err := time.Parse("15:04", startTime)
	if err != nil {
		return nil, errors.New("invalid start time format (HH:MM)")
	}
	end, err := time.Parse("15:04", endTime)
	if err != nil {
		return nil, errors.New("invalid end time format (HH:MM)")
	}

	if !start.Before(end) {
		return nil, errors.New("start time must be before end time")
	}

	date, err := time.Parse("2006-01-02", bookingDate)
	if err != nil {
		return nil, errors.New("invalid date format (YYYY-MM-DD)")
	}

	// 2. Validate Expert
	var expert models.Expert
	if err := db.First(&expert, expertID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("expert not found")
		}
		return nil, err
	}

	// 3. Check Master Availability (is it a valid working time?)
	dayOfWeek := int(date.Weekday())
	var slots []models.AvailabilitySlot
	if err := db.Where("expert_id = ? AND day_of_week = ?", expertID, dayOfWeek).Find(&slots).Error; err != nil {
		return nil, err
	}

	isWithinAvailability := false
	for _, slot := range slots {
		slotStart, _ := time.Parse("15:04", slot.StartTime)
		slotEnd, _ := time.Parse("15:04", slot.EndTime)

		// Check strict containment: slotStart <= start AND end <= slotEnd
		if (slotStart.Before(start) || slotStart.Equal(start)) && (slotEnd.After(end) || slotEnd.Equal(end)) {
			isWithinAvailability = true
			break
		}
	}

	if !isWithinAvailability {
		return nil, errors.New("requested time is not within expert's availability")
	}

	// 4. Check Intersecting Bookings
	// Overlap condition: Not (ExistingEnd <= NewStart OR ExistingStart >= NewEnd)
	// => ExistingEnd > NewStart AND ExistingStart < NewEnd
	var count int64
	err = db.Model(&models.Booking{}).
		Where("expert_id = ? AND booking_date = ? AND status != 'cancelled'", expertID, bookingDate).
		Where("start_time < ? AND end_time > ?", endTime, startTime).
		Count(&count).Error

	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, errors.New("time slot is already booked")
	}

	// 5. Calculate Total Price
	duration := end.Sub(start).Hours()
	totalPrice := duration * expert.HourlyRate

	// 6. Create Booking
	booking := models.Booking{
		UserID:      userID,
		ExpertID:    expertID,
		BookingDate: date.Format("2006-01-02"),
		StartTime:   start.Format("15:04"),
		EndTime:     end.Format("15:04"),
		TotalPrice:  totalPrice,
		Status:      models.BookingStatusPending,
	}

	if err := db.Create(&booking).Error; err != nil {
		return nil, err
	}

	return &booking, nil
}

func GetBookingsByUserID(userID uint) ([]models.Booking, error) {
	db := database.GetDB()
	var bookings []models.Booking
	if err := db.Preload("Expert.User").Where("user_id = ?", userID).Order("booking_date desc, start_time desc").Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}
