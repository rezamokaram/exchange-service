package handlers

import (
	"net/http"
	"qexchange/models"
	userModels "qexchange/models/user"
	"qexchange/services"

	"github.com/labstack/echo/v4"
)

// UserRegister handles the registration of a new user
// @Summary User registration
// @Description Register a new user
// @Tags User
// @Accept  json
// @Produce json
// @Param   body  body      userModels.RegisterRequest  true  "User Registration"
// @Success 200   {object}  models.Response
// @Failure 400   {object}  models.Response
// @Failure 500   {object}  models.Response
// @Router /user/register [post]
func UserRegister(service services.UserService) echo.HandlerFunc {
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

// UserLogin handles user login
// @Summary User login
// @Description Logs in a user
// @Tags User
// @Accept  json
// @Produce json
// @Param   body  body      userModels.LoginRequest     true  "User Login"
// @Success 200   {object}  userModels.LoginResponse
// @Failure 400   {object}  models.Response
// @Failure 500   {object}  models.Response
// @Router /user/login [post]
func UserLogin(service services.UserService) echo.HandlerFunc {
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
