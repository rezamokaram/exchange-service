package server

import (
	"qexchange/handlers"
	"qexchange/services"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func UserRoutes(e *echo.Echo, db *gorm.DB) {
	userService := services.NewUserService(db)
	e.POST(
		"/user/register",
		handlers.UserRegister(userService),
	)

	e.POST(
		"/user/login",
		handlers.UserLogin(userService),
	)
}
