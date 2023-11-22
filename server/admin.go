package server

import (
	"qexchange/handlers"
	"qexchange/middlewares"
	"qexchange/services"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func AdminRoutes(e *echo.Echo, db *gorm.DB) {
	AdminService := services.NewAdminService(db)
	e.PUT("/admin/update-to-admin", handlers.UpgradeToAdmin(AdminService), middlewares.AuthMiddleware(db))
	e.PUT("/admin/update-auth-level", handlers.UpdateAuthenticationLevel(AdminService), middlewares.AuthMiddleware(db), middlewares.AdminCheckMiddleware)
	e.PUT("/admin/block-user", handlers.BlockUser(AdminService), middlewares.AuthMiddleware(db), middlewares.AdminCheckMiddleware)
	e.PUT("/admin/unblock-user", handlers.UnblockUser(AdminService), middlewares.AuthMiddleware(db), middlewares.AdminCheckMiddleware)
	e.GET("/admin/user-info", handlers.GetUserInfo(AdminService), middlewares.AuthMiddleware(db), middlewares.AdminCheckMiddleware)
}
