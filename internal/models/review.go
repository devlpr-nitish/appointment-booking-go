package models

import "time"


type Review struct {
	ID        uint `gorm:"primaryKey"`
	BookingID uint
	Rating    int
	Comment   string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}