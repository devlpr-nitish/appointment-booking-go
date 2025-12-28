package routes

import (
	"github.com/devlpr-nitish/appointment-booking-go/internal/handlers"
	"github.com/devlpr-nitish/appointment-booking-go/internal/middleware"
	"github.com/labstack/echo/v4"
)

func BookingRoutes(e *echo.Echo) {
	g := e.Group("/bookings")
	g.Use(middleware.AuthMiddleware)

	g.POST("/create-booking", handlers.CreateBooking)
}
