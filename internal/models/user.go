package models

import "time"

type UserRole string

const (
	RoleUser   UserRole = "user"
	RoleAdmin  UserRole = "admin"
	RoleExpert UserRole = "expert"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	Password  string    `json:"-"` // Don't return password
	Role      UserRole  `gorm:"type:varchar(20)" json:"role"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
