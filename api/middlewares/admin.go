package middlewares

import (
	"github.com/RezaMokaram/ExchangeService/models"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func AdminCheckMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		tokenString := req.Header.Get("Authorization")

		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, models.NewErrorResponse("access denied", "Unauthorized"))
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			return c.JSON(http.StatusUnauthorized, models.NewErrorResponse("access denied", "Invalid token"))
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {

			return c.JSON(http.StatusUnauthorized, models.NewErrorResponse("access denied", "Invalid token claims"))
		}
		if adm, ok := claims["adm"].(bool); ok && adm {

			return next(c)
		}
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse("access denied", "user does not have admin access"))
	}
}
