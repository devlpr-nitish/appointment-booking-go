package routes

import (
	"github.com/devlpr-nitish/appointment-booking-go/internal/database"
	"github.com/devlpr-nitish/appointment-booking-go/internal/handlers"
	"github.com/devlpr-nitish/appointment-booking-go/internal/middleware"
	"github.com/devlpr-nitish/appointment-booking-go/internal/repositories"
	"github.com/devlpr-nitish/appointment-booking-go/internal/services"
	"github.com/devlpr-nitish/appointment-booking-go/internal/workers"
	"github.com/labstack/echo/v4"
)

var Hub *handlers.Hub

func NegotiationRoutes(e *echo.Echo) {
	// Initialize DB
	db := database.GetDB()

	// Initialize Repositories
	categoryRepo := repositories.NewCategoryRepository(db)
	requestRepo := repositories.NewRequestRepository(db)
	offerRepo := repositories.NewOfferRepository(db)

	// Initialize Services
	reqService := services.NewRequestService(requestRepo)
	offerService := services.NewOfferService(offerRepo, requestRepo)

	// Initialize Handlers
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)
	reqHandler := handlers.NewRequestHandler(reqService, Hub)
	offerHandler := handlers.NewOfferHandler(offerService, Hub)

	// Initialize WebSocket Hub
	Hub = handlers.NewHub()
	go Hub.Run()
	wsHandler := handlers.NewWebSocketHandler(Hub)

	// Start Background Workers
	workers.StartRequestCleaner(requestRepo, Hub)

	// Routes
	api := e.Group("/api")

	// WebSocket
	e.GET("/ws", wsHandler.HandleConnection)

	// Categories
	api.GET("/categories", categoryHandler.GetAllCategories)
	api.POST("/categories", categoryHandler.CreateCategory, middleware.AuthMiddleware)

	// Requests
	api.POST("/requests", reqHandler.CreateRequest, middleware.AuthMiddleware)
	api.GET("/requests/:id", reqHandler.GetRequest, middleware.AuthMiddleware) // User or Expert
	api.POST("/requests/:id/cancel", reqHandler.CancelRequest, middleware.AuthMiddleware) // User cancels
	api.GET("/expert/requests", reqHandler.GetExpertRequests, middleware.AuthMiddleware)

	// Offers
	api.GET("/requests/:id/offers", offerHandler.GetRequestOffers, middleware.AuthMiddleware)
	api.POST("/offers", offerHandler.CreateOffer, middleware.AuthMiddleware)
	api.POST("/offers/:id/accept", offerHandler.AcceptOffer, middleware.AuthMiddleware)
}
