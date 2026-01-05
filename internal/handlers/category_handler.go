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

func (h *CategoryHandler) GetAllCategories(c echo.Context) error {
	categories, err := h.repo.FindAll()
	if err != nil {
		return utils.RespondError(c, http.StatusInternalServerError, err, "failed to fetch categories")
	}
	return utils.RespondSuccess(c, http.StatusOK, "categories retrieved successfully", categories)
}

func (h *CategoryHandler) CreateCategory(c echo.Context) error {
	var category models.Category
	if err := c.Bind(&category); err != nil {
		return utils.RespondError(c, http.StatusBadRequest, err, "invalid request body")
	}

	if err := h.repo.Create(&category); err != nil {
		return utils.RespondError(c, http.StatusInternalServerError, err, "failed to create category")
	}

	return utils.RespondSuccess(c, http.StatusCreated, "category created successfully", category)
}
