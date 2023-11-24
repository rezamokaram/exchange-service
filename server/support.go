package server

import (
	"qexchange/handlers"
	"qexchange/middlewares"
	"qexchange/services"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SupportRoutes(e *echo.Echo, db *gorm.DB) {
	SupportService := services.NewSupportService(db)
	e.POST("/support/open-ticket", handlers.OpenTicket(SupportService), middlewares.AuthMiddleware(db))
	e.POST("/support/send-message", handlers.SendMessage(SupportService), middlewares.AuthMiddleware(db))
	e.GET("/support/get-ticket-messages", handlers.GetTicketMessages(SupportService), middlewares.AuthMiddleware(db))
	e.GET("/support/get-all-tickets", handlers.GetAllTickets(SupportService), middlewares.AuthMiddleware(db))

	e.GET("/support/admin/get-active-tickets", handlers.GetActiveTickets(SupportService), middlewares.AuthMiddleware(db), middlewares.AdminCheckMiddleware)
}
