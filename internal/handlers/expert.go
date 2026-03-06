package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/devlpr-nitish/appointment-booking-go/internal/models"
	"github.com/devlpr-nitish/appointment-booking-go/internal/services"
	"github.com/devlpr-nitish/appointment-booking-go/internal/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CreateExpertRequest struct {
	Bio        string  `json:"bio" validate:"required"`
	HourlyRate float64 `json:"hourly_rate" validate:"required,gt=0"`
}

// UpdateExpertRequest supports both a single category_id (legacy) and
// category_ids (multi-category). If category_ids is provided it takes precedence.
type UpdateExpertRequest struct {
	Bio         string      `json:"bio"`
	HourlyRate  float64     `json:"hourly_rate"`
	CategoryID  *uuid.UUID  `json:"category_id"`
	CategoryIDs []uuid.UUID `json:"category_ids"`
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

	expert, err := services.CreateExpertProfile(user.ID, req.Bio, req.HourlyRate)
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

	// Calculate profile completion percentage
	completionPercentage := calculateProfileCompletion(expert)

	// Build category_ids array for the response
	categoryIDs := make([]string, 0, len(expert.Categories))
	for _, cat := range expert.Categories {
		categoryIDs = append(categoryIDs, cat.ID.String())
	}

	response := map[string]interface{}{
		"expert":                expert,
		"completion_percentage": completionPercentage,
		"category_ids":          categoryIDs,
	}

	return utils.RespondSuccess(c, http.StatusOK, "expert profile retrieved successfully", response)
}

// calculateProfileCompletion calculates the profile completion percentage
func calculateProfileCompletion(expert *models.Expert) int {
	totalFields := 3 // bio, hourly_rate, at-least-one-category
	completedFields := 0

	if expert.Bio != "" {
		completedFields++
	}
	if expert.HourlyRate > 0 {
		completedFields++
	}
	if len(expert.Categories) > 0 || expert.CategoryID != nil {
		completedFields++
	}

	return (completedFields * 100) / totalFields
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

	// Decide which category list to use
	var categoryIDsPtr []uuid.UUID
	if len(req.CategoryIDs) > 0 {
		// Multi-category wins
		categoryIDsPtr = req.CategoryIDs
	} else if req.CategoryID != nil {
		// Fallback to single category
		categoryIDsPtr = []uuid.UUID{*req.CategoryID}
	}
	// If neither set, categoryIDsPtr remains nil → service won't change categories

	expert, err := services.UpdateExpertProfile(user.ID, req.Bio, req.HourlyRate, req.CategoryID, categoryIDsPtr)
	if err != nil {
		return utils.RespondError(c, http.StatusInternalServerError, err, "failed to update expert profile")
	}

	// Build category_ids for response
	categoryIDs := make([]string, 0, len(expert.Categories))
	for _, cat := range expert.Categories {
		categoryIDs = append(categoryIDs, cat.ID.String())
	}

	result := map[string]interface{}{
		"expert":       expert,
		"category_ids": categoryIDs,
	}

	return utils.RespondSuccess(c, http.StatusOK, "expert profile updated successfully", result)
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

func SearchExperts(c echo.Context) error {
	query := c.QueryParam("q")
	category := c.QueryParam("category")

	experts, err := services.SearchExperts(query, category)
	if err != nil {
		return utils.RespondError(c, http.StatusInternalServerError, err, "failed to search experts")
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

func GetExpertStats(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		return utils.RespondError(c, http.StatusUnauthorized, nil, "unauthorized")
	}

	expert, err := services.GetExpertProfile(user.ID)
	if err != nil {
		return utils.RespondError(c, http.StatusForbidden, err, "expert profile not found")
	}

	stats, err := services.GetExpertStats(expert.ID)
	if err != nil {
		return utils.RespondError(c, http.StatusInternalServerError, err, "failed to fetch expert stats")
	}

	return utils.RespondSuccess(c, http.StatusOK, "expert stats retrieved successfully", stats)
}
