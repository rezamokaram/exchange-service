package server

import (
	"qexchange/handlers"
	"qexchange/services"
	"qexchange/middlewares"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func PriceRoutes(e *echo.Echo, db *gorm.DB) {
	cryptoService := services.NewCryptoService(db)
	e.GET("/crypto", handlers.GetCrypto(cryptoService))
	e.POST("/crypto", handlers.SetCrypto(cryptoService), middlewares.AuthMiddleware(db))
	e.PUT("/crypto", handlers.UpdateCrypto(cryptoService), middlewares.AuthMiddleware(db))
	e.GET("/crypto/get-all", handlers.GetAllCrypto(cryptoService))
}
