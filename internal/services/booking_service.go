package services

import (
	"errors"

	"github.com/devlpr-nitish/appointment-booking-go/internal/database"
	"github.com/devlpr-nitish/appointment-booking-go/internal/models"
	"gorm.io/gorm"
)

func CreateBooking(userID, expertID, slotID uint) (*models.Booking, error) {
	db := database.GetDB()

	// 1. Validate Expert
	var expert models.Expert
	if err := db.First(&expert, expertID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("expert not found")
		}
		return nil, err
	}

	// 2. Validate Slot
	var slot models.AvailabilitySlot
	if err := db.First(&slot, slotID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("slot not found")
		}
		return nil, err
	}

	if slot.IsBooked {
		return nil, errors.New("slot is already booked")
	}

	if slot.ExpertID != expertID {
		return nil, errors.New("slot does not belong to the specified expert")
	}

	// 3. Create Booking & Update Slot (Transaction)
	tx := db.Begin()

	booking := models.Booking{
		UserID:   userID,
		ExpertID: expertID,
		SlotID:   slotID,
		Status:   models.BookingStatusConfirmed, // Auto-confirming for now
	}

	if err := tx.Create(&booking).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Model(&slot).Update("is_booked", true).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &booking, nil
}
