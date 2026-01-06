package handlers

import (
	"net/http"

	"github.com/devlpr-nitish/appointment-booking-go/internal/models"
	"github.com/devlpr-nitish/appointment-booking-go/internal/services"
	"github.com/devlpr-nitish/appointment-booking-go/internal/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type OfferHandler struct {
	service services.OfferService
	hub     *Hub
}

func NewOfferHandler(service services.OfferService, hub *Hub) *OfferHandler {
	return &OfferHandler{
		service: service,
		hub:     hub,
	}
}

type CreateOfferInput struct {
	RequestID uuid.UUID `json:"request_id" validate:"required"`
	Amount    float64   `json:"amount" validate:"required,gt=0"`
}

func (h *OfferHandler) CreateOffer(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		return utils.RespondError(c, http.StatusUnauthorized, nil, "unauthorized")
	}

	// Resolve ExpertID from UserID
	expert, err := services.GetExpertProfile(user.ID)
	if err != nil {
		return utils.RespondError(c, http.StatusForbidden, err, "user is not a registered expert")
	}

	// Determine request ID (bind from body)
	var input CreateOfferInput
	if err := c.Bind(&input); err != nil {
		return utils.RespondError(c, http.StatusBadRequest, err, "invalid request body")
	}

	offer, err := h.service.CreateOffer(input.RequestID, expert.ID, input.Amount)
	if err != nil {
		return utils.RespondError(c, http.StatusInternalServerError, err, "failed to create offer")
	}

	// Broadcast NEW_OFFER event via WebSocket to notify the user
	if h.hub != nil {
		h.hub.BroadcastEvent("NEW_OFFER", map[string]interface{}{
			"offer_id":   offer.ID,
			"request_id": offer.RequestID,
			"expert_id":  offer.ExpertID,
			"amount":     offer.Amount,
		})
	}

	return utils.RespondSuccess(c, http.StatusCreated, "offer created successfully", offer)
}

func (h *OfferHandler) GetRequestOffers(c echo.Context) error {
	requestIDStr := c.Param("id")
	requestID, err := uuid.Parse(requestIDStr)
	if err != nil {
		return utils.RespondError(c, http.StatusBadRequest, err, "invalid request id")
	}

	offers, err := h.service.GetOffersByRequestID(requestID)
	if err != nil {
		return utils.RespondError(c, http.StatusInternalServerError, err, "failed to get offers")
	}

	return utils.RespondSuccess(c, http.StatusOK, "offers retrieved successfully", offers)
}

func (h *OfferHandler) AcceptOffer(c echo.Context) error {
	offerIDStr := c.Param("id")
	offerID, err := uuid.Parse(offerIDStr)
	if err != nil {
		return utils.RespondError(c, http.StatusBadRequest, err, "invalid offer id")
	}

	if err := h.service.AcceptOffer(offerID); err != nil {
		return utils.RespondError(c, http.StatusInternalServerError, err, "failed to accept offer")
	}

	// Broadcast OFFER_ACCEPTED event via WebSocket
	if h.hub != nil {
		h.hub.BroadcastEvent("OFFER_ACCEPTED", map[string]interface{}{
			"offer_id": offerID,
		})
	}

	return utils.RespondSuccess(c, http.StatusOK, "offer accepted successfully", nil)
}
