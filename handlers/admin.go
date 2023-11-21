package handlers

import (
	"net/http"
	"qexchange/models"
	"qexchange/services"

	"github.com/labstack/echo/v4"
)

type UpdateUserToAdminRequest struct {
	AdminPassword string `json:"admin_password"`
}

type UpdateAuthenticationLevelRequest struct {
	Username     string `json:"username"`
	NewAuthLevel int    `json:"new_auth_level"`
}

type BlockUserRequest struct {
	Username  string `json:"username"`
	Temporary bool   `json:"temporary"`
}

type UnblockUserRequest struct {
	Username string `json:"username"`
}

func UpgradeToAdmin(service services.AdminService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var adminRequest struct {
			AdminPassword string `json:"admin_password"`
		}

		if err := c.Bind(&adminRequest); err != nil {
			return c.JSON(http.StatusBadRequest, models.NewErrorRespone("invalid request", err))
		}

		user := c.Get("user").(models.User)

		if err := service.UpgradeToAdmin(user, adminRequest.AdminPassword); err != nil {
			return c.JSON(http.StatusInternalServerError, models.NewErrorRespone("failed to upgrade user to admin", err))
		}

		return c.JSON(http.StatusOK, models.NewRespone("user role upgraded to admin"))
	}
}

func UpdateAuthenticationLevel(service services.AdminService) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := new(UpdateAuthenticationLevelRequest)
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, models.NewErrorRespone("invalid request", err))
		}
		err := service.UpdateAuthenticationLevel(request.Username, request.NewAuthLevel)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.NewErrorRespone("failed to update user Authentication Level", err))
		}
		return c.JSON(http.StatusOK, models.NewRespone("user Authentication Level updated "))
	}
}

func BlockUser(service services.AdminService) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := new(BlockUserRequest)
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, models.NewErrorRespone("invalid request", err))
		}

		if statusCode, err := service.BlockUser(request.Username, request.Temporary); err != nil {
			return c.JSON(statusCode, models.NewErrorRespone("failed to update user Blocked Level", err))
		}

		return c.JSON(http.StatusOK, models.NewRespone("user Blocked level updated"))
	}
}

func UnblockUser(service services.AdminService) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := new(UnblockUserRequest)
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, models.NewErrorRespone("invalid request", err))
		}

		if statusCode, err := service.UnblockUser(request.Username); err != nil {
			return c.JSON(statusCode, models.NewErrorRespone("failed to update user Blocked Level", err))
		}

		return c.JSON(http.StatusOK, models.NewRespone("user Blocked level updated"))
	}
}
