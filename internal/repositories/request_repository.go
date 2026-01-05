package repositories

import (
	"github.com/devlpr-nitish/appointment-booking-go/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RequestRepository interface {
	Create(request *models.Request) error
	FindByID(id uuid.UUID) (*models.Request, error)
	FindAllOpenByCategory(categoryID uuid.UUID) ([]models.Request, error)
	UpdateStatus(id uuid.UUID, status models.RequestStatus) error
}

type requestRepository struct {
	db *gorm.DB
}

func NewRequestRepository(db *gorm.DB) RequestRepository {
	return &requestRepository{db: db}
}

func (r *requestRepository) Create(request *models.Request) error {
	return r.db.Create(request).Error
}

func (r *requestRepository) FindByID(id uuid.UUID) (*models.Request, error) {
	var request models.Request
	err := r.db.Preload("User").Preload("Category").First(&request, "id = ?", id).Error
	return &request, err
}

func (r *requestRepository) FindAllOpenByCategory(categoryID uuid.UUID) ([]models.Request, error) {
	var requests []models.Request
	err := r.db.Preload("User").
		Where("category_id = ? AND status = ?", categoryID, models.RequestStatusOpen).
		Find(&requests).Error
	return requests, err
}

func (r *requestRepository) UpdateStatus(id uuid.UUID, status models.RequestStatus) error {
	return r.db.Model(&models.Request{}).Where("id = ?", id).Update("status", status).Error
}
