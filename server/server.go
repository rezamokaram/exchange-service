package server

import (
	"log"

	_ "qexchange/docs"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"
)

func NewServer() *echo.Echo {
	return echo.New()
}

func RunServer(e *echo.Echo, db *gorm.DB) {
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	UserRoutes(e, db)
	PriceRoutes(e, db)
	TradeRoutes(e, db)
	BankRoutes(e, db)
	AdminRoutes(e, db)
	SupportRoutes(e, db)
	log.Fatal(e.Start(":8080"))
}
