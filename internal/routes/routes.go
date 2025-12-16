package routes

import (
	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Echo) {
	HealthRoutes(e)
	AuthRoutes(e)
	UserRoutes(e)
	ExpertRoutes(e)
	BookingRoutes(e)
	PaymentRoutes(e)
	ReviewRoutes(e)
	AvailabilityRoutes(e)

	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "Welcome to Appointment Booking API",
		})
	})
}
