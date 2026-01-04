package models

import "time"

type Expert struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null;uniqueIndex" json:"user_id"`
	Bio         string    `json:"bio"`
	Expertise   string    `json:"expertise"`
	HourlyRate  float64   `json:"hourly_rate"`
	IsVerified  bool      `json:"is_verified" gorm:"default:false"`
	User        User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	GeneratedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type ExpertStats struct {
	TotalEarnings     float64 `json:"totalEarnings"`
	UpcomingSessions  int     `json:"upcomingSessions"`
	CompletedSessions int     `json:"completedSessions"`
	AverageRating     float64 `json:"averageRating"`
	PendingRequests   int     `json:"pendingRequests"`
}
