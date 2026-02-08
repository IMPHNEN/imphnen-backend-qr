package middleware

import (
	"github.com/IMPHNEN/imphnen-backend-qr/internal/utils"
	"github.com/labstack/echo/v4"
)

func RBACMiddleware(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, ok := c.Get("role").(string)
			if !ok {
				return utils.ErrorResponse(c, 403, "access denied", "forbidden")
			}

			for _, allowed := range allowedRoles {
				if role == allowed {
					return next(c)
				}
			}

			return utils.ErrorResponse(c, 403, "insufficient permissions", "forbidden")
		}
	}
}
