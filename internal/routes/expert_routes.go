package routes

import (
	"github.com/devlpr-nitish/appointment-booking-go/internal/handlers"
	"github.com/devlpr-nitish/appointment-booking-go/internal/middleware"
	"github.com/labstack/echo/v4"
)

func ExpertRoutes(e *echo.Echo) {
	g := e.Group("/expert")
	g.Use(middleware.AuthMiddleware)

	g.POST("/profile", handlers.CreateExpertProfile)
	g.GET("/profile", handlers.GetExpertProfile)
	g.PATCH("/profile", handlers.UpdateExpertProfile)
	g.GET("/get-experts", handlers.GetExperts)
	g.GET("/search", handlers.GetExpertByCatergoryName)
	g.GET("/get-expert-by-id/:id", handlers.GetExpertById)
}
