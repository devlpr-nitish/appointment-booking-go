package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/devlpr-nitish/appointment-booking-go/internal/models"
	"github.com/devlpr-nitish/appointment-booking-go/internal/services"
	"github.com/devlpr-nitish/appointment-booking-go/internal/utils"
	"github.com/labstack/echo/v4"
)

type CreateExpertRequest struct {
	Bio        string  `json:"bio" validate:"required"`
	Expertise  string  `json:"expertise" validate:"required"`
	HourlyRate float64 `json:"hourly_rate" validate:"required,gt=0"`
}

type UpdateExpertRequest struct {
	Bio        string  `json:"bio"`
	Expertise  string  `json:"expertise"`
	HourlyRate float64 `json:"hourly_rate"`
}

func CreateExpertProfile(c echo.Context) error {
	var req CreateExpertRequest
	if err := c.Bind(&req); err != nil {
		return utils.RespondError(c, http.StatusBadRequest, err, "invalid request body")
	}

	user, ok := c.Get("user").(*models.User)
	if !ok {
		return utils.RespondError(c, http.StatusUnauthorized, nil, "unauthorized")
	}

	expert, err := services.CreateExpertProfile(user.ID, req.Bio, req.Expertise, req.HourlyRate)
	if err != nil {
		return utils.RespondError(c, http.StatusInternalServerError, err, "failed to create expert profile")
	}

	return utils.RespondSuccess(c, http.StatusCreated, "expert profile created successfully", expert)
}

func GetExpertProfile(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		return utils.RespondError(c, http.StatusUnauthorized, nil, "unauthorized")
	}

	expert, err := services.GetExpertProfile(user.ID)
	if err != nil {
		return utils.RespondError(c, http.StatusNotFound, err, "expert profile not found")
	}

	return utils.RespondSuccess(c, http.StatusOK, "expert profile retrieved successfully", expert)
}

func UpdateExpertProfile(c echo.Context) error {
	var req UpdateExpertRequest
	if err := c.Bind(&req); err != nil {
		return utils.RespondError(c, http.StatusBadRequest, err, "invalid request body")
	}

	user, ok := c.Get("user").(*models.User)
	if !ok {
		return utils.RespondError(c, http.StatusUnauthorized, nil, "unauthorized")
	}

	expert, err := services.UpdateExpertProfile(user.ID, req.Bio, req.Expertise, req.HourlyRate)
	if err != nil {
		return utils.RespondError(c, http.StatusInternalServerError, err, "failed to update expert profile")
	}

	return utils.RespondSuccess(c, http.StatusOK, "expert profile updated successfully", expert)
}

func GetExperts(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 10
	}

	experts, total, err := services.GetExperts(page, limit)
	if err != nil {
		return utils.RespondError(c, http.StatusInternalServerError, err, "failed to get experts")
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	response := map[string]interface{}{
		"experts": experts,
		"meta": map[string]interface{}{
			"current_page": page,
			"total_pages":  totalPages,
			"total_items":  total,
			"limit":        limit,
		},
	}

	return utils.RespondSuccess(c, http.StatusOK, "experts retrieved successfully", response)
}

func GetExpertByCatergoryName(c echo.Context) error {
	categoryName := c.QueryParam("category")
	if categoryName == "" {
		return utils.RespondError(c, http.StatusBadRequest, nil, "category name is required")
	}

	experts, err := services.GetExpertByCatergoryName(categoryName)
	if err != nil {
		return utils.RespondError(c, http.StatusInternalServerError, err, "failed to get experts")
	}

	return utils.RespondSuccess(c, http.StatusOK, "experts retrieved successfully", experts)
}

func GetExpertById(c echo.Context) error {
	idStr := c.Param("id")
	if idStr == "" {
		return utils.RespondError(c, http.StatusBadRequest, nil, "expert id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.RespondError(c, http.StatusBadRequest, err, "invalid expert id")
	}

	expert, err := services.GetExpertById(uint(id))
	if err != nil {
		return utils.RespondError(c, http.StatusInternalServerError, err, "failed to get expert")
	}

	return utils.RespondSuccess(c, http.StatusOK, "expert retrieved successfully", expert)
}
