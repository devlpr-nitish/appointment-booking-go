package models

import "time"

type AvailabilitySlot struct {
	ID        uint      `gorm:"primaryKey"`
	ExpertID  uint      `gorm:"uniqueIndex"`
	StartTime time.Time
	EndTime   time.Time
	IsBooked  bool      `gorm:"default:false"`
	Expert    Expert    `gorm:"foreignKey:ExpertID"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}