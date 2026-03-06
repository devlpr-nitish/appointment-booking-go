package repositories

import (
	"github.com/devlpr-nitish/appointment-booking-go/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *models.Category) error
	FindAll() ([]models.Category, error)
	FindByID(id uuid.UUID) (*models.Category, error)
	SearchByName(query string) ([]models.Category, error)
	FindOrCreate(name string) (*models.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) FindAll() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Order("name ASC").Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) FindByID(id uuid.UUID) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, "id = ?", id).Error
	return &category, err
}

// SearchByName returns categories whose name contains the query (case-insensitive), limit 10.
func (r *categoryRepository) SearchByName(query string) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Where("name ILIKE ?", "%"+query+"%").Order("name ASC").Limit(10).Find(&categories).Error
	return categories, err
}

// FindOrCreate looks up a category by (case-insensitive) exact name and creates it if missing.
func (r *categoryRepository) FindOrCreate(name string) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("LOWER(name) = LOWER(?)", name).First(&category).Error
	if err == nil {
		return &category, nil // already exists
	}
	// Create new
	category = models.Category{Name: name}
	if err := r.db.Create(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}
