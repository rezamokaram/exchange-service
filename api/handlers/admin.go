package handlers

import (
	"net/http"

	"github.com/RezaMokaram/ExchangeService/internal"
	"github.com/RezaMokaram/ExchangeService/models"
	adminModels "github.com/RezaMokaram/ExchangeService/models/admin"
	userModels "github.com/RezaMokaram/ExchangeService/models/user"

	"github.com/labstack/echo/v4"
)

// UpgradeToAdmin upgrades a user to an admin role
// @Summary Upgrade user to admin
// @Description Upgrades a regular user to an admin
// @Tags Admin
// @Accept  json
// @Produce json
// @Param   body  body      UpdateUserToAdminRequest  true  "Admin Password"
// @Param Authorization header string true "Authorization token"
// @Security BasicAuth
// @Success 200   {object}  models.Response
// @Failure 400   {object}  models.Response
// @Router /admin/update-to-admin [put]
func UpgradeToAdmin(service internal.AdminService) echo.HandlerFunc {
	return func(c echo.Context) error {
		adminRequest := new(adminModels.UpdateUserToAdminRequest)

		if err := c.Bind(&adminRequest); err != nil {
			return c.JSON(http.StatusBadRequest, models.NewErrorResponse("invalid request", err.Error()))
		}

		user, bind := c.Get("user").(userModels.User)
		if !bind {
			response := models.NewErrorResponse("Upgrade To Admin failed", "bad user data")
			return c.JSON(http.StatusBadRequest, response)
		}

		if err := service.UpgradeToAdmin(user, adminRequest.AdminPassword); err != nil {
			response := models.NewErrorResponse("failed to upgrade user to admin", err.Error())
			return c.JSON(http.StatusInternalServerError, response)
		}

		return c.JSON(http.StatusOK, models.NewResponse("user role upgraded to admin"))
	}
}

// UpdateAuthenticationLevel updates a user's authentication level
// @Summary Update user authentication level
// @Description Updates a user's authentication level
// @Tags Admin
// @Accept  json
// @Produce json
// @Param   body  body      UpdateAuthenticationLevelRequest  true  "Update Auth Level"
// @Param Authorization header string true "Authorization token"
// @Security BasicAuth
// @Security AdminAuth
// @Success 200   {object}  models.Response
// @Failure 400   {object}  models.Response
// @Router /admin/update-auth-level [put]
func UpdateAuthenticationLevel(service internal.AdminService) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := new(adminModels.UpdateAuthenticationLevelRequest)
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, models.NewErrorResponse("invalid request", err.Error()))
		}
		err := service.UpdateAuthenticationLevel(request.Username, request.NewAuthLevel)
		if err != nil {
			response := models.NewErrorResponse("failed to update user Authentication Level", err.Error())
			return c.JSON(http.StatusInternalServerError, response)
		}
		return c.JSON(http.StatusOK, models.NewResponse("user Authentication Level updated "))
	}
}

// BlockUser blocks a user account
// @Summary Block user
// @Description Blocks a user account
// @Tags Admin
// @Accept  json
// @Produce json
// @Param   body  body      BlockUserRequest  true  "Block User"
// @Param Authorization header string true "Authorization token"
// @Security BasicAuth
// @Security AdminAuth
// @Success 200   {object}  models.Response
// @Failure 400   {object}  models.Response
// @Router /admin/block-user [put]
func BlockUser(service internal.AdminService) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := new(adminModels.BlockUserRequest)
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, models.NewErrorResponse("invalid request", err.Error()))
		}

		if statusCode, err := service.BlockUser(request.Username, request.Temporary); err != nil {
			return c.JSON(statusCode, models.NewErrorResponse("failed to update user Blocked Level", err.Error()))
		}

		return c.JSON(http.StatusOK, models.NewResponse("user Blocked level updated"))
	}
}

// UnblockUser unblocks a user account
// @Summary Unblock user
// @Description Unblocks a user account
// @Tags Admin
// @Accept  json
// @Produce json
// @Param   body  body      UnblockUserRequest  true  "Unblock User"
// @Param Authorization header string true "Authorization token"
// @Security BasicAuth
// @Security AdminAuth
// @Success 200   {object}  models.Response
// @Failure 400   {object}  models.Response
// @Router /admin/unblock-user [put]
func UnblockUser(service internal.AdminService) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := new(adminModels.UnblockUserRequest)
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, models.NewErrorResponse("invalid request", err.Error()))
		}

		if statusCode, err := service.UnblockUser(request.Username); err != nil {
			return c.JSON(statusCode, models.NewErrorResponse("failed to update user Blocked Level", err.Error()))
		}

		return c.JSON(http.StatusOK, models.NewResponse("user Blocked level updated"))
	}
}

// GetUserInfo retrieves information for a specific user
// @Summary Get user information
// @Description Retrieves information for a specific user
// @Tags Admin
// @Accept  json
// @Produce json
// @Param   username  query  string  true  "Username"
// @Param Authorization header string true "Authorization token"
// @Security BasicAuth
// @Security AdminAuth
// @Success 200   {object}  models.UserInfo
// @Failure 400   {object}  models.Response
// @Router /admin/user-info [get]
func GetUserInfo(service internal.AdminService) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.QueryParam("username")

		userInfo, statusCode, err := service.GetUserInfo(username)
		if err != nil {
			return c.JSON(statusCode, models.NewErrorResponse("failed to get user information", err.Error()))
		}

		return c.JSON(http.StatusOK, userInfo)
	}
}
