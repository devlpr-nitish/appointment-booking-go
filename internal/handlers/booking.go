package handlers

import (
	"net/http"

	"github.com/devlpr-nitish/appointment-booking-go/internal/models"
	"github.com/devlpr-nitish/appointment-booking-go/internal/services"
	"github.com/devlpr-nitish/appointment-booking-go/internal/utils"
	"github.com/labstack/echo/v4"
)

type CreateBookingRequest struct {
	ExpertID uint `json:"expert_id" validate:"required"`
	SlotID   uint `json:"slot_id" validate:"required"`
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

	booking, err := services.CreateBooking(user.ID, req.ExpertID, req.SlotID)
	if err != nil {
		return utils.RespondError(c, http.StatusInternalServerError, err, "failed to create booking")
	}

	return utils.RespondSuccess(c, http.StatusCreated, "booking created successfully", booking)
}
