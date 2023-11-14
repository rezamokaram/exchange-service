package server

import (
	"qexchange/handlers"
	"qexchange/services"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)


func BankRoutes(e *echo.Echo, db *gorm.DB) {
	bankService := services.NewBankService(db) //going to implement buisness logic
	e.POST("/bank/charge", handlers.BankCharge(bankService))
}