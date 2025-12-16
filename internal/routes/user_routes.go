package routes

import (
	"github.com/labstack/echo/v4"
)

func UserRoutes(e *echo.Echo) {
	_ = e.Group("/users")
	// g.GET("", handlers.GetUsers)
	// g.POST("", handlers.CreateUser)
}
