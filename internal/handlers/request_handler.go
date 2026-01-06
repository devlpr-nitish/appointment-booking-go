package handlers

import (
	"net/http"

	"github.com/devlpr-nitish/appointment-booking-go/internal/models"
	"github.com/devlpr-nitish/appointment-booking-go/internal/services"
	"github.com/devlpr-nitish/appointment-booking-go/internal/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type RequestHandler struct {
	service services.RequestService
	hub     *Hub
}

func NewRequestHandler(service services.RequestService, hub *Hub) *RequestHandler {
	return &RequestHandler{
		service: service,
		hub:     hub,
	}
}

type CreateRequestInput struct {
	CategoryID  uuid.UUID `json:"category_id" validate:"required"`
	Amount      float64   `json:"amount" validate:"required,gt=0"`
	Description string    `json:"description"`
}

func (h *RequestHandler) CreateRequest(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		return utils.RespondError(c, http.StatusUnauthorized, nil, "unauthorized")
	}

	// Only users (not experts) can create requests
	if user.Role != models.RoleUser {
		return utils.RespondError(c, http.StatusForbidden, nil, "only users can create requests")
	}

	var input CreateRequestInput
	if err := c.Bind(&input); err != nil {
		return utils.RespondError(c, http.StatusBadRequest, err, "invalid request body")
	}

	if err := c.Validate(&input); err != nil {
		return utils.RespondError(c, http.StatusBadRequest, err, "validation failed")
	}

	req, err := h.service.CreateRequest(user.ID, input.CategoryID, input.Amount, input.Description)
	if err != nil {
		return utils.RespondError(c, http.StatusInternalServerError, err, "failed to create request")
	}

	// Broadcast NEW_REQUEST event via WebSocket to notify experts
	if h.hub != nil {
		h.hub.BroadcastEvent("NEW_REQUEST", map[string]interface{}{
			"request_id":  req.ID,
			"category_id": req.CategoryID,
			"amount":      req.InitialAmount,
			"user_id":     req.UserID,
		})
	}

	return utils.RespondSuccess(c, http.StatusCreated, "request created successfully", req)
}

func (h *RequestHandler) GetExpertRequests(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		return utils.RespondError(c, http.StatusUnauthorized, nil, "unauthorized")
	}

	// Fetch Expert Profile to get CategoryID
	expert, err := services.GetExpertProfile(user.ID)
	if err != nil {
		return utils.RespondError(c, http.StatusForbidden, err, "user is not a registered expert")
	}

	if expert.CategoryID == nil {
		return utils.RespondError(c, http.StatusBadRequest, nil, "expert has no assigned category")
	}

	requests, err := h.service.GetOpenRequestsByCategory(*expert.CategoryID)
	if err != nil {
		return utils.RespondError(c, http.StatusInternalServerError, err, "failed to fetch requests")
	}

	return utils.RespondSuccess(c, http.StatusOK, "requests retrieved successfully", requests)
}

func (h *RequestHandler) GetRequest(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return utils.RespondError(c, http.StatusBadRequest, err, "invalid request id")
	}

	req, err := h.service.GetRequestByID(id)
	if err != nil {
		return utils.RespondError(c, http.StatusNotFound, err, "request not found")
	}

	return utils.RespondSuccess(c, http.StatusOK, "request retrieved successfully", req)
}
