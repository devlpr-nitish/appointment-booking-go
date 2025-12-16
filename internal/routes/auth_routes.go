package routes

import (
	"github.com/devlpr-nitish/appointment-booking-go/internal/handlers"
	"github.com/labstack/echo/v4"
)

func AuthRoutes(e *echo.Echo) {
	g := e.Group("/auth")
	g.POST("/register", handlers.Register)
	g.POST("/login", handlers.Login)
}
