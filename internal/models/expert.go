package models

import "time"

type Expert struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"uniqueIndex" json:"user_id"`
	Bio        string    `json:"bio"`
	Expertise  string    `json:"expertise"`
	HourlyRate float64   `json:"hourly_rate"`
	IsVerified bool      `gorm:"default:false" json:"is_verified"`
	User       User      `gorm:"foreignKey:UserID" json:"user"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
