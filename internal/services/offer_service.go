package services

import (
	"errors"

	"github.com/devlpr-nitish/appointment-booking-go/internal/models"
	"github.com/devlpr-nitish/appointment-booking-go/internal/repositories"
	"github.com/google/uuid"
)

type OfferService interface {
	CreateOffer(requestID uuid.UUID, expertID uint, amount float64) (*models.Offer, error)
	GetOffersByRequestID(requestID uuid.UUID) ([]models.Offer, error)
	AcceptOffer(offerID uuid.UUID) error
}

type offerService struct {
	offerRepo   repositories.OfferRepository
	requestRepo repositories.RequestRepository
}

func NewOfferService(offerRepo repositories.OfferRepository, requestRepo repositories.RequestRepository) OfferService {
	return &offerService{
		offerRepo:   offerRepo,
		requestRepo: requestRepo,
	}
}

func (s *offerService) CreateOffer(requestID uuid.UUID, expertID uint, amount float64) (*models.Offer, error) {
	// Validate request state
	req, err := s.requestRepo.FindByID(requestID)
	if err != nil {
		return nil, err
	}
	if req.Status != models.RequestStatusOpen {
		return nil, errors.New("cannot offer on non-open request")
	}

	offer := &models.Offer{
		RequestID: requestID,
		ExpertID:  expertID,
		Amount:    amount,
		Status:    models.OfferStatusPending,
	}
	err = s.offerRepo.Create(offer)
	return offer, err
}

func (s *offerService) GetOffersByRequestID(requestID uuid.UUID) ([]models.Offer, error) {
	return s.offerRepo.FindByRequestID(requestID)
}

func (s *offerService) AcceptOffer(offerID uuid.UUID) error {
	offer, err := s.offerRepo.FindByID(offerID)
	if err != nil {
		return err
	}
	if offer.Status != models.OfferStatusPending {
		return errors.New("offer is not pending")
	}

	return s.offerRepo.AcceptOffer(offerID, offer.RequestID)
}
