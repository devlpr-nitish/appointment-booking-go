package models

import "time"

type BookingStatus string

const (
	BookingStatusPending BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusCompleted BookingStatus = "completed"
)

type Booking struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"uniqueIndex"`
	ExpertID  uint      `gorm:"uniqueIndex"`
	SlotID    uint      `gorm:"uniqueIndex"`
	Status    BookingStatus
	User      User      `gorm:"foreignKey:UserID"`
	Expert    Expert    `gorm:"foreignKey:ExpertID"`
	Slot      AvailabilitySlot `gorm:"foreignKey:SlotID"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
