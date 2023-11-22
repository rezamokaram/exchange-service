package server

import (
	"qexchange/handlers"
	"qexchange/middlewares"
	"qexchange/services"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func TradeRoutes(e *echo.Echo, db *gorm.DB) {
	tradeService := services.NewTradeService(db)
	e.POST("/open-trade", handlers.OpenTrade(tradeService), middlewares.AuthMiddleware(db))
	e.POST("/close-trade", handlers.CloseTrade(tradeService), middlewares.AuthMiddleware(db))
	e.GET("/open-trade/get-all", handlers.GetAllOpenTrades(tradeService), middlewares.AuthMiddleware(db))
	e.GET("/close-trade/get-all", handlers.GetAllClosedTrades(tradeService), middlewares.AuthMiddleware(db))
	e.GET("/close-trade/filter-all", handlers.FilterTrades(tradeService), middlewares.AuthMiddleware(db))
	e.POST("/future-order", handlers.SetFutureOrder(tradeService), middlewares.AuthMiddleware(db))
	e.GET("/future-order/get-all", handlers.GetAllFutureOrders(tradeService), middlewares.AuthMiddleware(db))
}