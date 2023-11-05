package server

import (
	"log"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var e *echo.Echo

func init() {
	e = echo.New()
}

func RunServer(db *gorm.DB) {
	log.Fatal(e.Start(":8080"))
}
