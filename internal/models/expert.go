package models

import "time"

type Expert struct {
	ID         uint `gorm:"primaryKey"`
	UserID     uint `gorm:"uniqueIndex"`
	Bio        string
	Expertise  string
	HourlyRate float64
	IsVerified bool      `gorm:"default:false"`
	User       User      `gorm:"foreignKey:UserID"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}
