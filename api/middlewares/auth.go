package middlewares

import (
	"fmt"
	"github.com/RezaMokaram/ExchangeService/models"
	userModels "github.com/RezaMokaram/ExchangeService/models/user"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// CustomClaims represents custom JWT claims.
type CustomClaims struct {
	ID    uint  `json:"id"`
	Exp   int64 `json:"exp"`
	Admin bool  `json:"adm"`
	jwt.StandardClaims
}

// AuthMiddleware is the authentication middleware.
func AuthMiddleware(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			tokenString := req.Header.Get("Authorization")
			// Check if the token is missing
			if tokenString == "" {
				return c.JSON(http.StatusUnauthorized, models.NewErrorResponse("access denied", "Authorization token is missing"))
			}

			// Parse the token
			token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
				// Verify the signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(os.Getenv("JWT_SECRET")), nil
			})

			if err != nil {
				response := models.NewErrorResponse("access denied", "Invalid token")
				return c.JSON(http.StatusUnauthorized, response)
			}

			// Validate token claims
			claims, ok := token.Claims.(*CustomClaims)
			if !ok || !token.Valid {
				response := models.NewErrorResponse("access denied", "Invalid token claims")
				return c.JSON(http.StatusUnauthorized, response)
			}

			// Check token expiration
			if time.Now().Unix() > claims.Exp {
				response := models.NewErrorResponse("access denied", "Token has expired")
				return c.JSON(http.StatusUnauthorized, response)
			}

			// Fetch user from database using claims
			var user userModels.User
			if err := db.First(&user, claims.ID).Preload("Profile").Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					response := models.NewErrorResponse("access denied", "User not found")
					return c.JSON(http.StatusUnauthorized, response)
				}
				response := models.NewErrorResponse("access denied", "Internal server error")
				return c.JSON(http.StatusInternalServerError, response)
			}

			// Set user object in context
			c.Set("user", user)

			// Proceed to the next handler
			return next(c)
		}
	}
}
