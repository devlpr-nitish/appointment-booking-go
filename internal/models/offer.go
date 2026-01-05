package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OfferStatus string

const (
	OfferStatusPending  OfferStatus = "PENDING"
	OfferStatusAccepted OfferStatus = "ACCEPTED"
	OfferStatusDeclined OfferStatus = "DECLINED"
)

type Offer struct {
	ID        uuid.UUID   `gorm:"type:uuid;primary_key;" json:"id"`
	RequestID uuid.UUID   `gorm:"type:uuid;not null" json:"request_id"`
	ExpertID  uint        `gorm:"not null" json:"expert_id"`
	Amount    float64     `json:"amount"`
	Status    OfferStatus `gorm:"type:varchar(20);default:'PENDING'" json:"status"`
	Request   Request     `gorm:"foreignKey:RequestID" json:"request,omitempty"`
	Expert    Expert      `gorm:"foreignKey:ExpertID" json:"expert,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

func (o *Offer) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return
}
