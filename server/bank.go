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
	e.POST("/bank/add_account", handlers.AddBankAccount(bankService), middlewares.AuthMiddleware(db))
	e.POST("/bank/payment/withdraw", handlers.WithdrawFromAccount(bankService), middlewares.AuthMiddleware(db))
	e.POST("/bank/payment/charge", handlers.ChargeAccount(bankService), middlewares.AuthMiddleware(db))
	e.GET("/bank/payment/verify", handlers.VerifyPayment(bankService))
}
