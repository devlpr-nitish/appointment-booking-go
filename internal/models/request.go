package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RequestStatus string

const (
	RequestStatusOpen     RequestStatus = "OPEN"
	RequestStatusAccepted RequestStatus = "ACCEPTED"
	RequestStatusClosed   RequestStatus = "CLOSED"
)

type Request struct {
	ID            uuid.UUID     `gorm:"type:uuid;primary_key;" json:"id"`
	UserID        uint          `gorm:"not null" json:"user_id"`
	CategoryID    uuid.UUID     `gorm:"type:uuid;not null" json:"category_id"`
	InitialAmount float64       `json:"initial_amount"`
	Description   string        `json:"description"`
	Status        RequestStatus `gorm:"type:varchar(20);default:'OPEN'" json:"status"`
	Models        []Offer       `gorm:"foreignKey:RequestID" json:"offers,omitempty"`
	User          User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Category      Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

func (r *Request) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return
}
