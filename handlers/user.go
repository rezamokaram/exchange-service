package handlers

import (
	"net/http"
	"qexchange/services"

	"github.com/labstack/echo/v4"
)

type RegisterRequest struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	PasswordRepeat string `json:"passwordrepeat"`
}
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func UserRegister(service services.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		// parse body
		request := new(RegisterRequest)
		if err := c.Bind(request); err != nil {
			response := UserResponse{
				Error:   err.Error(),
				Message: "",
			}
			return c.JSON(http.StatusBadRequest, response)
		}

		// call the register service
		statusCode, err := service.Register(request.Username, request.Password, request.PasswordRepeat, request.Email)
		if err != nil {
			response := UserResponse{
				Error:   err.Error(),
				Message: "",
			}
			return c.JSON(statusCode, response)
		}

		return c.JSON(http.StatusOK, UserResponse{
			Error:   "",
			Message: "user created successfuly",
		})
	}
}

func UserLogin(service services.UserService) echo.HandlerFunc {
	return func(c echo.Context) error {
		// hanlde here
		// should call service.Login()
		return c.JSON(http.StatusOK, "user")
	}
}
