package models

import (
	"time"

	"github.com/google/uuid"
)

// ExpertCategory is the join table for the Expert <-> Category many-to-many relationship
type ExpertCategory struct {
	ExpertID   uint      `gorm:"primaryKey" json:"expert_id"`
	CategoryID uuid.UUID `gorm:"primaryKey;type:uuid" json:"category_id"`
}

type Expert struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	UserID     uint       `gorm:"not null;uniqueIndex" json:"user_id"`
	Bio        string     `json:"bio"`
	HourlyRate float64    `json:"hourly_rate"`
	IsVerified bool       `json:"is_verified" gorm:"default:false"`
	// Legacy single-category field (kept for backward compat)
	CategoryID *uuid.UUID `gorm:"type:uuid" json:"category_id"`
	Category   Category   `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	// Multiple categories via join table
	Categories  []Category `gorm:"many2many:expert_categories;" json:"categories,omitempty"`
	User        User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	GeneratedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

type ExpertStats struct {
	TotalEarnings     float64 `json:"totalEarnings"`
	UpcomingSessions  int     `json:"upcomingSessions"`
	CompletedSessions int     `json:"completedSessions"`
	AverageRating     float64 `json:"averageRating"`
	PendingRequests   int     `json:"pendingRequests"`
}
