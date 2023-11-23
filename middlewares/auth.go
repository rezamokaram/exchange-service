package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"qexchange/models"
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
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Authorization token is missing"})
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
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token"})
			}

			// Validate token claims
			claims, ok := token.Claims.(*CustomClaims)
			if !ok || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token claims"})
			}

			// Check token expiration
			if time.Now().Unix() > claims.Exp {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Token has expired"})
			}

			// Fetch user from database using claims
			var user models.User
			if err := db.First(&user, claims.ID).Preload("Profile").Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return c.JSON(http.StatusUnauthorized, map[string]string{"message": "User not found"})
				}
				return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
			}

			// Set user object in context
			c.Set("user", user)

			// Proceed to the next handler
			return next(c)
		}
	}
}
