package routes

import (
	"github.com/labstack/echo/v4"
)

func AvailabilityRoutes(e *echo.Echo) {
	_ = e.Group("/availability")
	// g.GET("", handlers.GetAvailability)
}
