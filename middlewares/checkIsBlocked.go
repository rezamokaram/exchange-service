package middlewares

import (
	"net/http"
	"qexchange/models"
	"qexchange/services"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CheckIsBlocked(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, bind := c.Get("user").(models.User)
			if !bind {
				response := models.NewErrorRespone("bad user data", nil)
				return c.JSON(http.StatusBadRequest, response)
			}

			var profile models.Profile
			if db.Where("user_id = ?", user.ID).First(&profile).Error != nil {
				response := models.NewErrorRespone("profile not found", nil)
				return c.JSON(http.StatusBadRequest, response)
			}

			if profile.BlockedLevel != services.Unblocked {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "User is blocked"})
			}

			return next(c)
		}
	}
}
