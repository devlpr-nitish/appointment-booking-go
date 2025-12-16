package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/devlpr-nitish/appointment-booking-go/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return utils.RespondError(c, http.StatusUnauthorized, nil, "Missing authorization token")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return utils.RespondError(c, http.StatusUnauthorized, nil, "Invalid token format")
		}

		tokenString := parts[1]
		jwtSecret := []byte(os.Getenv("JWT_SECRET"))

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return utils.RespondError(c, http.StatusUnauthorized, err, "Invalid or expired token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return utils.RespondError(c, http.StatusUnauthorized, nil, "Invalid token claims")
		}

		// Set user_id in context
		if userID, ok := claims["user_id"].(float64); ok {
			c.Set("user_id", uint(userID))
		} else {
			return utils.RespondError(c, http.StatusUnauthorized, nil, "Invalid user ID in token")
		}

		return next(c)
	}
}
