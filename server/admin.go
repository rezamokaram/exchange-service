package server

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"qexchange/handlers"
	"qexchange/middlewares"
	"qexchange/services"
)

func AdminRoutes(e *echo.Echo, db *gorm.DB) {
	AdminService := services.NewAdminService(db)
	e.PUT("/admin/update-to-admin", handlers.UpgradeToAdmin(AdminService), middlewares.AuthMiddleware(db))
	e.PUT("/admin/update-auth-level", handlers.UpdateAuthenticationLevel(AdminService), middlewares.AuthMiddleware(db), middlewares.AdminCheckMiddleware)
}
