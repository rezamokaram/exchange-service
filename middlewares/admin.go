package middlewares

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func AdminCheckMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims, ok := user.Claims.(jwt.MapClaims)

		if !ok || !claims["adm"].(bool) {
			return echo.NewHTTPError(http.StatusForbidden, "Access denied")
		}

		return next(c)
	}
}
