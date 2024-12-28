package handlers

import (
	"net/http"
	"github.com/RezaMokaram/ExchangeService/models"
	userModels "github.com/RezaMokaram/ExchangeService/models/user"
	svc "github.com/RezaMokaram/ExchangeService/internal/user"

	"github.com/labstack/echo/v4"
)

func UserRegister(service svc.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := new(userModels.RegisterRequest)
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, models.NewErrorResponse("registration failed", err.Error()))
		}

		if err := request.IsValid(); err != nil {
			return c.JSON(
				http.StatusBadRequest,
				models.NewErrorResponse(
					"registration failed",
					err.Error(),
				),
			)
		}

		statusCode, err := service.Register(request.Username, request.Password, request.PasswordRepeat, request.Email)
		if err != nil {
			return c.JSON(statusCode, models.NewErrorResponse("registration failed", err.Error()))
		}

		return c.JSON(http.StatusOK, models.NewResponse("user created successfully"))
	}
}

func UserLogin(service svc.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := new(userModels.LoginRequest)
		err := c.Bind(request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.NewErrorResponse("login failed", err.Error()))
		}

		if err := request.IsValid(); err != nil {
			return c.JSON(
				http.StatusBadRequest,
				models.NewErrorResponse(
					"login failed",
					err.Error(),
				),
			)
		}

		status, token, err := service.Login(request.Username, request.Password)
		if err != nil {
			return c.JSON(status, models.NewErrorResponse("login failed", err.Error()))
		}

		tokenResponse := userModels.LoginResponse{
			Token: token,
		}

		return c.JSON(http.StatusOK, tokenResponse)
	}
}
