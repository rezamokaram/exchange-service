package server

import (
	"log"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func NewServer() *echo.Echo {
	return echo.New()
}

func RunServer(e *echo.Echo, db *gorm.DB) {
	// user routes
	UserRoutes(e, db)
	PriceRoutes(e, db)
	
	log.Fatal(e.Start(":8080"))
}
