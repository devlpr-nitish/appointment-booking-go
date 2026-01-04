package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/devlpr-nitish/appointment-booking-go/internal/database"
	"github.com/devlpr-nitish/appointment-booking-go/internal/models"
	"github.com/devlpr-nitish/appointment-booking-go/internal/services"
	"github.com/devlpr-nitish/appointment-booking-go/internal/utils"
	"github.com/labstack/echo/v4"
)

type CreateBookingRequest struct {
	ExpertID    uint   `json:"expert_id" validate:"required"`
	BookingDate string `json:"booking_date" validate:"required"`
	StartTime   string `json:"start_time" validate:"required"`
	EndTime     string `json:"end_time" validate:"required"`
}

func CreateBooking(c echo.Context) error {
	var req CreateBookingRequest
	if err := c.Bind(&req); err != nil {
		return utils.RespondError(c, http.StatusBadRequest, err, "invalid request body")
	}

	user, ok := c.Get("user").(*models.User)
	if !ok {
		return utils.RespondError(c, http.StatusUnauthorized, nil, "unauthorized")
	}

	booking, err := services.CreateBooking(user.ID, req.ExpertID, req.BookingDate, req.StartTime, req.EndTime)
	if err != nil {
		status := http.StatusInternalServerError
		errMsg := err.Error()
		if strings.Contains(errMsg, "not found") {
			status = http.StatusNotFound
		} else if strings.Contains(errMsg, "booked") || strings.Contains(errMsg, "overlap") {
			status = http.StatusConflict
		} else if strings.Contains(errMsg, "invalid") || strings.Contains(errMsg, "availability") || strings.Contains(errMsg, "within") {
			status = http.StatusBadRequest
		}
		return utils.RespondError(c, status, err, "failed to create booking")
	}

	return utils.RespondSuccess(c, http.StatusCreated, "booking created successfully", booking)
}

func GetUserBookings(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		return utils.RespondError(c, http.StatusUnauthorized, nil, "unauthorized")
	}

	bookings, err := services.GetBookingsByUserID(user.ID)
	if err != nil {
		return utils.RespondError(c, http.StatusInternalServerError, err, "failed to fetch bookings")
	}

	return utils.RespondSuccess(c, http.StatusOK, "bookings fetched successfully", bookings)
}

type CancelBookingRequest struct {
	Reason string `json:"reason" validate:"required"`
}

func CancelBooking(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return utils.RespondError(c, http.StatusBadRequest, nil, "invalid booking id")
	}
	bookingID := uint(id)

	var req CancelBookingRequest
	if err := c.Bind(&req); err != nil {
		return utils.RespondError(c, http.StatusBadRequest, err, "invalid request body")
	}

	user, ok := c.Get("user").(*models.User)
	if !ok {
		return utils.RespondError(c, http.StatusUnauthorized, nil, "unauthorized")
	}

	if err := services.CancelBooking(bookingID, user.ID, req.Reason); err != nil {
		return utils.RespondError(c, http.StatusInternalServerError, err, "failed to cancel booking")
	}

	return utils.RespondSuccess(c, http.StatusOK, "booking cancelled successfully", nil)
}

func GetExpertBookings(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		return utils.RespondError(c, http.StatusUnauthorized, nil, "unauthorized")
	}

	// Ideally we should check if user is an expert here, or rely on middleware/service check
	// Assuming an expert user has an associated Expert record or ID.
	// Based on the code, we might need to find the expert ID associated with this user.
	// Let's assume for now the user context has what we need or we look it up.
	// Checking the service CreateBooking, it takes expertID passed from frontend.
	// But here the logged in user IS the expert.
	// We need to find the Expert record for this UserID.

	db := database.GetDB()
	var expert models.Expert
	if err := db.Where("user_id = ?", user.ID).First(&expert).Error; err != nil {
		return utils.RespondError(c, http.StatusForbidden, err, "expert profile not found")
	}

	status := c.QueryParam("status")
	bookings, err := services.GetBookingsByExpertID(expert.ID, status)
	if err != nil {
		return utils.RespondError(c, http.StatusInternalServerError, err, "failed to fetch bookings")
	}

	return utils.RespondSuccess(c, http.StatusOK, "bookings fetched successfully", bookings)
}

type UpdateBookingStatusRequest struct {
	Status string `json:"status" validate:"required"`
}

func UpdateBookingStatus(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return utils.RespondError(c, http.StatusBadRequest, nil, "invalid booking id")
	}
	bookingID := uint(id)

	var req UpdateBookingStatusRequest
	if err := c.Bind(&req); err != nil {
		return utils.RespondError(c, http.StatusBadRequest, err, "invalid request body")
	}

	user, ok := c.Get("user").(*models.User)
	if !ok {
		return utils.RespondError(c, http.StatusUnauthorized, nil, "unauthorized")
	}

	db := database.GetDB()
	var expert models.Expert
	if err := db.Where("user_id = ?", user.ID).First(&expert).Error; err != nil {
		return utils.RespondError(c, http.StatusForbidden, err, "expert profile not found")
	}

	if err := services.UpdateBookingStatus(bookingID, expert.ID, models.BookingStatus(req.Status)); err != nil {
		return utils.RespondError(c, http.StatusBadRequest, err, "failed to update booking status")
	}

	return utils.RespondSuccess(c, http.StatusOK, "booking status updated successfully", nil)
}
