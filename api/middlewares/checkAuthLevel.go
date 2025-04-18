package middlewares

import (
	"net/http"

	"github.com/rezamokaram/exchange-service/internal"
	"github.com/rezamokaram/exchange-service/models"
	userModels "github.com/rezamokaram/exchange-service/models/user"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CheckAuthLevel(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, bind := c.Get("user").(userModels.User)
			if !bind {
				return c.JSON(http.StatusBadRequest, models.NewErrorResponse("Check Auth", "bad user data"))
			}

			var profile userModels.Profile
			if db.Where("user_id = ?", user.ID).First(&profile).Error != nil {
				return c.JSON(http.StatusBadRequest, models.NewErrorResponse("Check Auth", "profile not found"))
			}

			if profile.AuthenticationLevel == internal.Unauthenticated {
				return c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Check Auth", "User is not authorized"))
			}

			return next(c)
		}
	}
}
