package models

import "time"

type PaymentStatus string

const (
	PaymentInitiated PaymentStatus = "initiated"
	PaymentCompleted PaymentStatus = "completed"
	PaymentFailed    PaymentStatus = "failed"
	PaymentRefunded  PaymentStatus = "refunded"
)

type Payment struct {
	ID        uint `gorm:"primaryKey"`
	BookingID uint
	Amount    float64
	Status    PaymentStatus `gorm:"type:varchar(20)"`
	Provider  string
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
