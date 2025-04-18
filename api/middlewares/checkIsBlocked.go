package middlewares

import (
	"net/http"

	"github.com/rezamokaram/exchange-service/internal"
	"github.com/rezamokaram/exchange-service/models"
	userModels "github.com/rezamokaram/exchange-service/models/user"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CheckIsBlocked(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, bind := c.Get("user").(userModels.User)
			if !bind {
				response := models.NewErrorResponse("Check Blocked", "bad user data")
				return c.JSON(http.StatusBadRequest, response)
			}

			var profile userModels.Profile
			if db.Where("user_id = ?", user.ID).First(&profile).Error != nil {
				response := models.NewErrorResponse("Check Blocked", "profile not found")
				return c.JSON(http.StatusBadRequest, response)
			}

			if profile.BlockedLevel != internal.Unblocked {
				return c.JSON(http.StatusUnauthorized, models.NewErrorResponse("Check Blocked", "User is blocked"))
			}

			return next(c)
		}
	}
}
