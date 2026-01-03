package handlers

import (
	"net/http"
	"strings"

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
