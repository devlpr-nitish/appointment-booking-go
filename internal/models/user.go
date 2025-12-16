package models

import "time"

type UserRole string

const (
	RoleUser UserRole = "user"
	RoleAdmin UserRole = "admin"
	RoleExpert UserRole = "expert"
)

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string
	Email     string    `gorm:"uniqueIndex"`
	Password  string
	Role      UserRole  `gorm:"type:varchar(20)"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}