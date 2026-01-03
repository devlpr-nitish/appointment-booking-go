package models

import "time"

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusCompleted BookingStatus = "completed"
)

type Booking struct {
	ID                 uint          `gorm:"primaryKey"`
	UserID             uint          `gorm:"index"`
	ExpertID           uint          `gorm:"index"`
	BookingDate        string        `gorm:"not null;default:'1970-01-01'" json:"booking_date"` // YYYY-MM-DD
	StartTime          string        `gorm:"not null;default:'00:00'" json:"start_time"`        // HH:MM
	EndTime            string        `gorm:"not null;default:'00:00'" json:"end_time"`          // HH:MM
	TotalPrice         float64       `gorm:"not null;default:0" json:"total_price"`
	CancellationReason string        `json:"cancellation_reason"`
	Status             BookingStatus `json:"status"`
	User               User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Expert             Expert        `gorm:"foreignKey:ExpertID" json:"expert,omitempty"`
	CreatedAt          time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}
