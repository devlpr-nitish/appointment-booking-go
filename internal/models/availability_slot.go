package models

import "time"

type AvailabilitySlot struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ExpertID  uint      `gorm:"not null" json:"expert_id"`
	DayOfWeek int       `gorm:"not null" json:"day_of_week"` // 0-6 (Sunday-Saturday)
	StartTime string    `gorm:"not null" json:"start_time"`  // Format: "HH:MM"
	EndTime   string    `gorm:"not null" json:"end_time"`    // Format: "HH:MM"
	Expert    Expert    `gorm:"foreignKey:ExpertID" json:"-"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
