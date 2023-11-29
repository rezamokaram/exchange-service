package server

import (
	"qexchange/handlers"
	"qexchange/middlewares"
	"qexchange/services"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func PriceRoutes(e *echo.Echo, db *gorm.DB) {
	cryptoService := services.NewCryptoService(db)

	e.GET(
		"/crypto",
		handlers.GetCrypto(cryptoService),
	)

	e.POST(
		"/crypto",
		handlers.SetCrypto(cryptoService),
		middlewares.AuthMiddleware(db),
		middlewares.CheckIsBlocked(db),
		middlewares.AdminCheckMiddleware,
	)

	e.PUT(
		"/crypto",
		handlers.UpdateCrypto(cryptoService),
		middlewares.AuthMiddleware(db),
		middlewares.CheckIsBlocked(db),
		middlewares.AdminCheckMiddleware,
	)

	e.GET(
		"/crypto/get-all",
		handlers.GetAllCrypto(cryptoService),
	)
}
