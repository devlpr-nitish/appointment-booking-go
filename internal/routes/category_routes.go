package routes

import (
	"github.com/devlpr-nitish/appointment-booking-go/internal/database"
	"github.com/devlpr-nitish/appointment-booking-go/internal/handlers"
	"github.com/devlpr-nitish/appointment-booking-go/internal/repositories"
	"github.com/labstack/echo/v4"
)

func CategoryRoutes(e *echo.Echo) {
	repo := repositories.NewCategoryRepository(database.GetDB())
	h := handlers.NewCategoryHandler(repo)

	g := e.Group("/categories")
	g.GET("", h.GetAllCategories)
	g.GET("/search", h.SearchCategories)
	g.POST("", h.CreateCategory) // FindOrCreate (idempotent)
}
