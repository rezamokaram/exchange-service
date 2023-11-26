package handlers

import (
	"net/http"
	"qexchange/models"
	"qexchange/services"

	"github.com/labstack/echo/v4"
)

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Username       string `json:"username" example:"newUser"`
	Email          string `json:"email" example:"newUser@example.com"`
	Password       string `json:"password" example:"123456"`
	PasswordRepeat string `json:"passwordrepeat" example:"123456"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Username string `json:"username" example:"newUser"`
	Password string `json:"password" example:"123456"`
}

// TokenResponse represents the response containing a JWT token
type TokenResponse struct {
	Token string `json:"token"`
}

// UserRegister handles the registration of a new user
// @Summary User registration
// @Description Register a new user
// @Accept  json
// @Produce json
// @Param   body  body      RegisterRequest  true  "User Registration"
// @Success 200   {object}  models.Response
// @Failure 400   {object}  models.Response
// @Failure 500   {object}  models.Response
// @Router /user/register [post]
func UserRegister(service services.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		// parse body
		request := new(RegisterRequest)
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, models.NewErrorResponse("", err.Error()))
		}

		// call the register service
		statusCode, err := service.Register(request.Username, request.Password, request.PasswordRepeat, request.Email)
		if err != nil {
			return c.JSON(statusCode, models.NewErrorResponse("", err.Error()))
		}

		return c.JSON(http.StatusOK, models.NewResponse("user created successfuly"))
	}
}

// UserLogin handles user login
// @Summary User login
// @Description Logs in a user
// @Accept  json
// @Produce json
// @Param   body  body      LoginRequest     true  "User Login"
// @Success 200   {object}  TokenResponse
// @Failure 400   {object}  models.Response
// @Failure 500   {object}  models.Response
// @Router /user/login [post]
func UserLogin(service services.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := new(LoginRequest)
		err := c.Bind(request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.NewErrorResponse("", err.Error()))
		}
		status, token, err := service.Login(request.Username, request.Password)
		if err != nil {
			return c.JSON(status, models.NewErrorResponse("", err.Error()))
		}

		tokenrespone := TokenResponse{
			Token: token,
		}

		return c.JSON(http.StatusOK, tokenrespone)
	}
}
