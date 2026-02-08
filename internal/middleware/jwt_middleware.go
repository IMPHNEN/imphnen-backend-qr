package middleware

import (
	"strings"

	"github.com/IMPHNEN/imphnen-backend-qr/internal/utils"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return utils.ErrorResponse(c, 401, "missing authorization header", "unauthorized")
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return utils.ErrorResponse(c, 401, "invalid authorization format", "unauthorized")
			}

			claims, err := utils.ValidateToken(secret, parts[1])
			if err != nil {
				return utils.ErrorResponse(c, 401, "invalid or expired token", "unauthorized")
			}

			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)
			c.Set("role", claims.Role)

			return next(c)
		}
	}
}
