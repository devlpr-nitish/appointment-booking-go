package routes

import (
	"github.com/devlpr-nitish/appointment-booking-go/internal/handlers"
	"github.com/devlpr-nitish/appointment-booking-go/internal/middleware"
	"github.com/labstack/echo/v4"
)

func ExpertRoutes(e *echo.Echo) {
	g := e.Group("/expert")

	// Public routes (no auth required)
	g.GET("/get-experts", handlers.GetExperts)
	g.GET("/search", handlers.GetExpertByCatergoryName)
	g.GET("/get-expert-by-id/:id", handlers.GetExpertById)
	g.GET("/available-slots", handlers.GetAvailableSlots)

	// Protected routes (auth required)
	g.Use(middleware.AuthMiddleware)

	g.POST("/profile", handlers.CreateExpertProfile)
	g.GET("/profile", handlers.GetExpertProfile)
	g.PATCH("/profile", handlers.UpdateExpertProfile)

	// Availability routes
	g.POST("/availability", handlers.CreateAvailability)
	g.GET("/availability", handlers.GetAvailability)
	g.PATCH("/availability/:id", handlers.UpdateAvailability)
	g.DELETE("/availability/:id", handlers.DeleteAvailability)
}
