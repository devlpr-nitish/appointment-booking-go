package repositories

import (
	"github.com/devlpr-nitish/appointment-booking-go/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OfferRepository interface {
	Create(offer *models.Offer) error
	FindByID(id uuid.UUID) (*models.Offer, error)
	FindByRequestID(requestID uuid.UUID) ([]models.Offer, error)
	AcceptOffer(offerID uuid.UUID, requestID uuid.UUID) error
}

type offerRepository struct {
	db *gorm.DB
}

func NewOfferRepository(db *gorm.DB) OfferRepository {
	return &offerRepository{db: db}
}

func (r *offerRepository) Create(offer *models.Offer) error {
	return r.db.Create(offer).Error
}

func (r *offerRepository) FindByID(id uuid.UUID) (*models.Offer, error) {
	var offer models.Offer
	err := r.db.Preload("Expert").Preload("Expert.User").First(&offer, "id = ?", id).Error
	return &offer, err
}

func (r *offerRepository) FindByRequestID(requestID uuid.UUID) ([]models.Offer, error) {
	var offers []models.Offer
	err := r.db.Preload("Expert").Preload("Expert.User").
		Where("request_id = ?", requestID).
		Find(&offers).Error
	return offers, err
}

func (r *offerRepository) AcceptOffer(offerID uuid.UUID, requestID uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Accept the specific offer
		if err := tx.Model(&models.Offer{}).
			Where("id = ?", offerID).
			Update("status", models.OfferStatusAccepted).Error; err != nil {
			return err
		}

		// 2. Decline all other offers for this request
		if err := tx.Model(&models.Offer{}).
			Where("request_id = ? AND id != ?", requestID, offerID).
			Update("status", models.OfferStatusDeclined).Error; err != nil {
			return err
		}

		// 3. Mark the request as ACCEPTED
		if err := tx.Model(&models.Request{}).
			Where("id = ?", requestID).
			Update("status", models.RequestStatusAccepted).Error; err != nil {
			return err
		}

		return nil
	})
}
