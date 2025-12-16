package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func RespondSuccess(c echo.Context, status int, message string, data interface{}) error {
	return c.JSON(status, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func RespondError(c echo.Context, status int, err error, details interface{}) error {
	return c.JSON(status, APIResponse{
		Success: false,
		Message: http.StatusText(status),
		Error: map[string]interface{}{
			"reason":  err.Error(),
			"details": details,
		},
	})
}
