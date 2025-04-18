package http

import (
	"fmt"

	"github.com/rezamokaram/exchange-service/api/handlers/http/crypto"
	"github.com/rezamokaram/exchange-service/api/handlers/http/middlewares"
	"github.com/rezamokaram/exchange-service/api/handlers/http/user"
	"github.com/rezamokaram/exchange-service/app"
	"github.com/rezamokaram/exchange-service/config"

	"github.com/gofiber/fiber/v2"
)

func Run(appContainer app.App, cfg config.ServerConfig) error {
	router := fiber.New()

	router.Use(middlewares.RequestLogger())
	api := router.Group("/api/v1", middlewares.SetUserContext)

	user.RegisterAuthAPI(appContainer, cfg, api)
	crypto.RegisterCryptoAPI(appContainer, cfg, api)

	return router.Listen(fmt.Sprintf(":%d", cfg.HttpPort))
}
