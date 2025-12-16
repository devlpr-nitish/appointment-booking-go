package middleware

import (
	"net/http"
	"strings"

	"github.com/devlpr-nitish/appointment-booking-go/internal/database"
	"github.com/devlpr-nitish/appointment-booking-go/internal/models"
	"github.com/devlpr-nitish/appointment-booking-go/internal/utils"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return utils.RespondError(c, http.StatusUnauthorized, nil, "missing authorization header")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return utils.RespondError(c, http.StatusUnauthorized, nil, "invalid authorization header format")
		}

		claims, err := utils.ValidateJWT(parts[1])
		if err != nil {
			return utils.RespondError(c, http.StatusUnauthorized, err, "invalid or expired token")
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return utils.RespondError(c, http.StatusUnauthorized, nil, "invalid token claims")
		}

		userID := uint(userIDFloat)

		var user models.User
		if err := database.GetDB().First(&user, userID).Error; err != nil {
			return utils.RespondError(c, http.StatusUnauthorized, err, "user not found")
		}

		c.Set("user", &user)
		return next(c)
	}
}

func ExpertOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(*models.User)
		if !ok {
			return utils.RespondError(c, http.StatusUnauthorized, nil, "unauthorized")
		}

		// Normalize role case comparison just in case, though usually strict
		if user.Role != models.RoleExpert {
			return utils.RespondError(c, http.StatusForbidden, nil, "access denied: experts only")
		}

		return next(c)
	}
}
