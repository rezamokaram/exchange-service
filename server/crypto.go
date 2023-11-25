package server

import (
	"qexchange/handlers"
	"qexchange/services"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func PriceRoutes(e *echo.Echo, db *gorm.DB) {
	cryptoService := services.NewCryptoService(db)
	e.GET("/crypto", handlers.GetCrypto(cryptoService))
	e.POST("/crypto", handlers.SetCrypto(cryptoService))
	e.PUT("/crypto", handlers.UpdateCrypto(cryptoService))
	e.GET("/crypto/getall", handlers.GetAllCrypto(cryptoService))
}
