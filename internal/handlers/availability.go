package handlers

import (
	"net/http"
	"strconv"

	"github.com/devlpr-nitish/appointment-booking-go/internal/models"
	"github.com/devlpr-nitish/appointment-booking-go/internal/services"
	"github.com/devlpr-nitish/appointment-booking-go/internal/utils"
	"github.com/labstack/echo/v4"
)

type CreateAvailabilityRequest struct {
	DayOfWeek int    `json:"day_of_week" validate:"required,min=0,max=6"`
	StartTime string `json:"start_time" validate:"required"`
	EndTime   string `json:"end_time" validate:"required"`
}

type UpdateAvailabilityRequest struct {
	DayOfWeek int    `json:"day_of_week" validate:"required,min=0,max=6"`
	StartTime string `json:"start_time" validate:"required"`
	EndTime   string `json:"end_time" validate:"required"`
}

// CreateAvailability creates a new availability slot for the authenticated expert
func CreateAvailability(c echo.Context) error {
	var req CreateAvailabilityRequest
	if err := c.Bind(&req); err != nil {
		return utils.RespondError(c, http.StatusBadRequest, err, "invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return utils.RespondError(c, http.StatusBadRequest, err, "validation failed")
	}

	user, ok := c.Get("user").(*models.User)
	if !ok {
		return utils.RespondError(c, http.StatusUnauthorized, nil, "unauthorized")
	}

	// Get expert profile to ensure user is an expert
	expert, err := services.GetExpertProfile(user.ID)
	if err != nil {
		return utils.RespondError(c, http.StatusNotFound, err, "expert profile not found")
	}

	availability, err := services.CreateAvailability(expert.ID, req.DayOfWeek, req.StartTime, req.EndTime)
	if err != nil {
		return utils.RespondError(c, http.StatusBadRequest, err, "failed to create availability")
	}

	return utils.RespondSuccess(c, http.StatusCreated, "availability created successfully", availability)
}

// GetAvailability retrieves all availability slots for the authenticated expert
func GetAvailability(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		return utils.RespondError(c, http.StatusUnauthorized, nil, "unauthorized")
	}

	// Get expert profile
	expert, err := services.GetExpertProfile(user.ID)
	if err != nil {
		return utils.RespondError(c, http.StatusNotFound, err, "expert profile not found")
	}

	availability, err := services.GetAvailabilityByExpertID(expert.ID)
	if err != nil {
		return utils.RespondError(c, http.StatusInternalServerError, err, "failed to get availability")
	}

	return utils.RespondSuccess(c, http.StatusOK, "availability retrieved successfully", availability)
}

// UpdateAvailability updates an existing availability slot
func UpdateAvailability(c echo.Context) error {
	var req UpdateAvailabilityRequest
	if err := c.Bind(&req); err != nil {
		return utils.RespondError(c, http.StatusBadRequest, err, "invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return utils.RespondError(c, http.StatusBadRequest, err, "validation failed")
	}

	idStr := c.Param("id")
	if idStr == "" {
		return utils.RespondError(c, http.StatusBadRequest, nil, "availability id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.RespondError(c, http.StatusBadRequest, err, "invalid availability id")
	}

	user, ok := c.Get("user").(*models.User)
	if !ok {
		return utils.RespondError(c, http.StatusUnauthorized, nil, "unauthorized")
	}

	// Get expert profile
	expert, err := services.GetExpertProfile(user.ID)
	if err != nil {
		return utils.RespondError(c, http.StatusNotFound, err, "expert profile not found")
	}

	availability, err := services.UpdateAvailability(uint(id), expert.ID, req.DayOfWeek, req.StartTime, req.EndTime)
	if err != nil {
		return utils.RespondError(c, http.StatusBadRequest, err, "failed to update availability")
	}

	return utils.RespondSuccess(c, http.StatusOK, "availability updated successfully", availability)
}

// DeleteAvailability deletes an availability slot
func DeleteAvailability(c echo.Context) error {
	idStr := c.Param("id")
	if idStr == "" {
		return utils.RespondError(c, http.StatusBadRequest, nil, "availability id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.RespondError(c, http.StatusBadRequest, err, "invalid availability id")
	}

	user, ok := c.Get("user").(*models.User)
	if !ok {
		return utils.RespondError(c, http.StatusUnauthorized, nil, "unauthorized")
	}

	// Get expert profile
	expert, err := services.GetExpertProfile(user.ID)
	if err != nil {
		return utils.RespondError(c, http.StatusNotFound, err, "expert profile not found")
	}

	if err := services.DeleteAvailability(uint(id), expert.ID); err != nil {
		return utils.RespondError(c, http.StatusBadRequest, err, "failed to delete availability")
	}

	return utils.RespondSuccess(c, http.StatusOK, "availability deleted successfully", nil)
}

// GetAvailableSlots returns available time slots for an expert on a specific date
func GetAvailableSlots(c echo.Context) error {
	expertIDStr := c.QueryParam("expertId")
	if expertIDStr == "" {
		return utils.RespondError(c, http.StatusBadRequest, nil, "expertId query parameter is required")
	}

	expertID, err := strconv.ParseUint(expertIDStr, 10, 32)
	if err != nil {
		return utils.RespondError(c, http.StatusBadRequest, err, "invalid expertId")
	}

	date := c.QueryParam("date")
	if date == "" {
		return utils.RespondError(c, http.StatusBadRequest, nil, "date query parameter is required")
	}

	slots, err := services.GetAvailableSlots(uint(expertID), date)
	if err != nil {
		return utils.RespondError(c, http.StatusInternalServerError, err, "failed to get available slots")
	}

	return utils.RespondSuccess(c, http.StatusOK, "available slots retrieved successfully", map[string]interface{}{
		"slots": slots,
	})
}
