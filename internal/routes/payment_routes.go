package routes

import (
	"github.com/labstack/echo/v4"
)

func PaymentRoutes(e *echo.Echo) {
	_ = e.Group("/payments")
	// g.POST("", handlers.CreatePayment)
}
