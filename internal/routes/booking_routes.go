package routes

import (
	"github.com/labstack/echo/v4"
)

func BookingRoutes(e *echo.Echo) {
	_ = e.Group("/bookings")
	// g.POST("", handlers.CreateBooking)
}
