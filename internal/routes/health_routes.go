package routes

import (
	"github.com/devlpr-nitish/appointment-booking-go/internal/handlers"
	"github.com/labstack/echo/v4"
)

func HealthRoutes(e *echo.Echo) {
	e.GET("/health", handlers.HealthCheck)
}
