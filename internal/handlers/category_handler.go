package handlers

import (
	"net/http"

	"github.com/devlpr-nitish/appointment-booking-go/internal/models"
	"github.com/devlpr-nitish/appointment-booking-go/internal/repositories"
	"github.com/devlpr-nitish/appointment-booking-go/internal/utils"
	"github.com/labstack/echo/v4"
)

type CategoryHandler struct {
	repo repositories.CategoryRepository
}

func NewCategoryHandler(repo repositories.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{repo: repo}
}

// GET /categories — list all categories
func (h *CategoryHandler) GetAllCategories(c echo.Context) error {
	categories, err := h.repo.FindAll()
	if err != nil {
		return utils.RespondError(c, http.StatusInternalServerError, err, "failed to fetch categories")
	}
	return utils.RespondSuccess(c, http.StatusOK, "categories retrieved successfully", categories)
}

// GET /categories/search?q=... — fuzzy name search (up to 10 results)
func (h *CategoryHandler) SearchCategories(c echo.Context) error {
	q := c.QueryParam("q")
	if q == "" {
		// Return all when nothing typed yet
		categories, err := h.repo.FindAll()
		if err != nil {
			return utils.RespondError(c, http.StatusInternalServerError, err, "failed to fetch categories")
		}
		return utils.RespondSuccess(c, http.StatusOK, "categories retrieved successfully", categories)
	}

	categories, err := h.repo.SearchByName(q)
	if err != nil {
		return utils.RespondError(c, http.StatusInternalServerError, err, "failed to search categories")
	}
	return utils.RespondSuccess(c, http.StatusOK, "categories retrieved successfully", categories)
}

// POST /categories — create a new category (or return existing one with same name)
func (h *CategoryHandler) CreateCategory(c echo.Context) error {
	var body struct {
		Name string `json:"name"`
	}
	if err := c.Bind(&body); err != nil || body.Name == "" {
		return utils.RespondError(c, http.StatusBadRequest, nil, "name is required")
	}

	category, err := h.repo.FindOrCreate(body.Name)
	if err != nil {
		return utils.RespondError(c, http.StatusInternalServerError, err, "failed to create category")
	}
	return utils.RespondSuccess(c, http.StatusOK, "category ready", category)
}

// POST /categories/raw — create strictly new (original behaviour, used by admin)
func (h *CategoryHandler) CreateCategoryRaw(c echo.Context) error {
	var category models.Category
	if err := c.Bind(&category); err != nil {
		return utils.RespondError(c, http.StatusBadRequest, err, "invalid request body")
	}
	if err := h.repo.Create(&category); err != nil {
		return utils.RespondError(c, http.StatusInternalServerError, err, "failed to create category")
	}
	return utils.RespondSuccess(c, http.StatusCreated, "category created successfully", category)
}
