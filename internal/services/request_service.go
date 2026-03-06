package services

import (
	"github.com/devlpr-nitish/appointment-booking-go/internal/models"
	"github.com/devlpr-nitish/appointment-booking-go/internal/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RequestService interface {
	CreateRequest(userID uint, categoryID uuid.UUID, amount float64, description string) (*models.Request, error)
	GetOpenRequestsByCategory(categoryID uuid.UUID) ([]models.Request, error)
	GetOpenRequestsByCategories(categoryIDs []uuid.UUID) ([]models.Request, error)
	GetRequestByID(id uuid.UUID) (*models.Request, error)
	CancelRequest(id uuid.UUID, userID uint) error
}

type requestService struct {
	requestRepo repositories.RequestRepository
}

func NewRequestService(repo repositories.RequestRepository) RequestService {
	return &requestService{requestRepo: repo}
}

func (s *requestService) CreateRequest(userID uint, categoryID uuid.UUID, amount float64, description string) (*models.Request, error) {
	req := &models.Request{
		UserID:        userID,
		CategoryID:    categoryID,
		InitialAmount: amount,
		Description:   description,
		Status:        models.RequestStatusOpen,
	}
	err := s.requestRepo.Create(req)
	return req, err
}

func (s *requestService) GetOpenRequestsByCategory(categoryID uuid.UUID) ([]models.Request, error) {
	return s.requestRepo.FindAllOpenByCategory(categoryID)
}

func (s *requestService) GetOpenRequestsByCategories(categoryIDs []uuid.UUID) ([]models.Request, error) {
	return s.requestRepo.FindAllOpenByCategories(categoryIDs)
}

func (s *requestService) GetRequestByID(id uuid.UUID) (*models.Request, error) {
	return s.requestRepo.FindByID(id)
}

func (s *requestService) CancelRequest(id uuid.UUID, userID uint) error {
	req, err := s.requestRepo.FindByID(id)
	if err != nil {
		return err
	}
	// Verify user owns the request
	if req.UserID != userID {
		return gorm.ErrRecordNotFound // Or a custom unauthorized error
	}
	return s.requestRepo.UpdateStatus(id, models.RequestStatusCanceled)
}
