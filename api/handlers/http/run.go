package http

import (
	"fmt"

	"github.com/RezaMokaram/ExchangeService/api/handlers/http/crypto"
	"github.com/RezaMokaram/ExchangeService/api/handlers/http/middlewares"
	"github.com/RezaMokaram/ExchangeService/api/handlers/http/user"
	"github.com/RezaMokaram/ExchangeService/app"
	"github.com/RezaMokaram/ExchangeService/config"

	"github.com/gofiber/fiber/v2"
)

func Run(appContainer app.App, cfg config.ServerConfig) error {
	router := fiber.New()

	api := router.Group("/api/v1", middlewares.SetUserContext)

	user.RegisterAuthAPI(appContainer, cfg, api)
	crypto.RegisterCryptoAPI(appContainer, cfg, api)

	return router.Listen(fmt.Sprintf(":%d", cfg.HttpPort))
}
