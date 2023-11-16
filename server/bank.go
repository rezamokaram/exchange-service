package server

import (
	"qexchange/handlers"
	"qexchange/middlewares"
	"qexchange/services"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func BankRoutes(e *echo.Echo, db *gorm.DB) {
	bankService := services.NewBankService(db)
	e.POST("/payment/charge", handlers.ChargeAccount(bankService), middlewares.AuthMiddleware(db))
	e.GET("/payment/verify", handlers.VerifyPayment(bankService))
}
