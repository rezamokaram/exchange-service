package middlewares

import (
	"net/http"
	"qexchange/models"
	"qexchange/services"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CheckAuthLevel(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, bind := c.Get("user").(models.User)
			if !bind {
				return c.JSON(http.StatusBadRequest, models.NewErrorResponse("Check Auth", "bad user data"))
			}

			var profile models.Profile
			if db.Where("user_id = ?", user.ID).First(&profile).Error != nil {
				return c.JSON(http.StatusBadRequest, models.NewErrorResponse("Check Auth", "profile not found"))
			}

			if profile.AuthenticationLevel == services.Unauthenticated {
				return c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Check Auth", "User is not authorized"))
			}

			return next(c)
		}
	}
}
