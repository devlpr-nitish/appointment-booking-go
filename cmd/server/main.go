package main

import (
	"log"

	"github.com/devlpr-nitish/appointment-booking-go/internal/config"
	"github.com/devlpr-nitish/appointment-booking-go/internal/database"
	"github.com/devlpr-nitish/appointment-booking-go/internal/routes"
	"github.com/devlpr-nitish/appointment-booking-go/internal/utils"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	cfg := config.LoadConfig()
	e := echo.New()

	// Register validator
	e.Validator = utils.NewValidator()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization},
	}))

	db := database.Connect(cfg)

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB: %v", err)
	}
	defer sqlDB.Close()

	routes.Routes(e)

	log.Printf("Server is running on http://localhost:%s", cfg.AppPort)
	e.Logger.Fatal(e.Start(":" + cfg.AppPort))
}
