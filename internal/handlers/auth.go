package handlers

import (
	"net/http"

	"github.com/devlpr-nitish/appointment-booking-go/internal/services"
	"github.com/devlpr-nitish/appointment-booking-go/internal/utils"
	"github.com/labstack/echo/v4"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

func Register(c echo.Context) error {
	var req RegisterRequest

	if err := c.Bind(&req); err != nil {
		return utils.RespondError(c, http.StatusBadRequest, err, "Invalid request format")
	}

	if req.Email == "" || req.Password == "" {
		return utils.RespondError(c, http.StatusBadRequest, echo.NewHTTPError(http.StatusBadRequest, "Missing required field"), "email and password are required")
	}

	// Validate role
	if req.Role != "" && req.Role != "user" && req.Role != "expert" {
		return utils.RespondError(c, http.StatusBadRequest, echo.NewHTTPError(http.StatusBadRequest, "Invalid role"), "role must be either 'user' or 'expert'")
	}

	// Default to "user" if no role specified
	if req.Role == "" {
		req.Role = "user"
	}

	user, err := services.RegisterUser(req.Email, req.Password, req.Name, req.Role)

	if err != nil {
		return utils.RespondError(c, http.StatusInternalServerError, err, "Registration failed")
	}

	return utils.RespondSuccess(c, http.StatusCreated, "user registered successfully", user)
}

func Login(c echo.Context) error {
	var req LoginRequest

	if err := c.Bind(&req); err != nil {
		return utils.RespondError(c, http.StatusBadRequest, err, "Invalid request format")
	}

	if req.Email == "" || req.Password == "" {
		return utils.RespondError(c, http.StatusBadRequest, echo.NewHTTPError(http.StatusBadRequest, "Missing required field"), "email and password are required")
	}

	token, err := services.LoginUser(req.Email, req.Password)

	if err != nil {
		return utils.RespondError(c, http.StatusUnauthorized, err, "invalid email or password")
	}

	return utils.RespondSuccess(c, http.StatusOK, "user loggedin successfully", map[string]string{"token": token})
}
