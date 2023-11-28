package handlers

import (
	"net/http"
	"qexchange/models"
	"qexchange/services"

	"github.com/labstack/echo/v4"
)

// UpdateUserToAdminRequest represents the request body for upgrading a user to admin
type UpdateUserToAdminRequest struct {
	AdminPassword string `json:"admin_password" example:"secret"`
}

// UpdateAuthenticationLevelRequest represents the request for updating a user's authentication level
type UpdateAuthenticationLevelRequest struct {
	Username     string `json:"username" example:"user2"`
	NewAuthLevel int    `json:"new_auth_level" example:"0"`
}

// BlockUserRequest represents the request for blocking a user
type BlockUserRequest struct {
	Username  string `json:"username" example:"user1"`
	Temporary bool   `json:"temporary" example:"false"`
}

// UnblockUserRequest represents the request for unblocking a user
type UnblockUserRequest struct {
	Username string `json:"username" example:"user2"`
}

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
func UpgradeToAdmin(service services.AdminService) echo.HandlerFunc {
	return func(c echo.Context) error {
		adminRequest := new(UpdateUserToAdminRequest)

		if err := c.Bind(&adminRequest); err != nil {
			return c.JSON(http.StatusBadRequest, models.NewErrorResponse("invalid request", err.Error()))
		}

		user, bind := c.Get("user").(models.User)
		if !bind {
			response := models.NewErrorResponse("", "bad user data")
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
func UpdateAuthenticationLevel(service services.AdminService) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := new(UpdateAuthenticationLevelRequest)
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
func BlockUser(service services.AdminService) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := new(BlockUserRequest)
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
func UnblockUser(service services.AdminService) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := new(UnblockUserRequest)
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
func GetUserInfo(service services.AdminService) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.QueryParam("username")

		userInfo, statusCode, err := service.GetUserInfo(username)
		if err != nil {
			return c.JSON(statusCode, models.NewErrorResponse("failed to get user information", err.Error()))
		}

		return c.JSON(http.StatusOK, userInfo)
	}
}
