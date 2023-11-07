package handlers

import (
	"net/http"
	"qexchange/services"

	"github.com/labstack/echo/v4"
)

func UserRegister(service services.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		// hanlde here
		// should call service.Register()

		return c.JSON(http.StatusOK, "user")
	}
}

func UserLogin(service services.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		// hanlde here
		// should call service.Login()
		return c.JSON(http.StatusOK, "user")
	}
}
