package server

// import (
// 	"fmt"
// 	"log"

// 	"github.com/rezamokaram/exchange-service/config"
// 	_ "github.com/rezamokaram/exchange-service/docs"

// 	"github.com/labstack/echo/v4"
// 	echoSwagger "github.com/swaggo/echo-swagger"
// 	"gorm.io/gorm"
// )

// func NewServer() *echo.Echo {
// 	return echo.New()
// }

// func RunServer(e *echo.Echo, db *gorm.DB, cfg *config.ServerConfig) {
// 	e.GET("/swagger/*", echoSwagger.WrapHandler)
// 	// UserRoutes(e, db)
// 	// PriceRoutes(e, db)
// 	TradeRoutes(e, db)
// 	BankRoutes(e, db)
// 	AdminRoutes(e, db)
// 	SupportRoutes(e, db)
// 	log.Fatal(e.Start(fmt.Sprintf("%s:%v", cfg.HttpHost, cfg.HttpPort)))
// }
