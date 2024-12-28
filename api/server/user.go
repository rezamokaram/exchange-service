package server

import (
	"github.com/RezaMokaram/ExchangeService/api/handlers"
	svc "github.com/RezaMokaram/ExchangeService/internal/user"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func UserRoutes(e *echo.Echo, db *gorm.DB) {
	userService := svc.NewUserService(db)
	e.POST(
		"/user/register",
		handlers.UserRegister(userService),
	)

	e.POST(
		"/user/login",
		handlers.UserLogin(userService),
	)
}
