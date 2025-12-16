package routes

import (
	"github.com/labstack/echo/v4"
)

func ReviewRoutes(e *echo.Echo) {
	_ = e.Group("/reviews")
	// g.POST("", handlers.CreateReview)
}
