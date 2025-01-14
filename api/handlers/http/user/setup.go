package user

import (
	"fmt"

	"github.com/RezaMokaram/ExchangeService/api/handlers/http/middlewares"
	"github.com/RezaMokaram/ExchangeService/app"
	"github.com/RezaMokaram/ExchangeService/config"

	"github.com/gofiber/fiber/v2"
)

func Run(appContainer app.App, cfg config.ServerConfig) error {
	router := fiber.New()

	api := router.Group("/api/v1", middlewares.SetUserContext)

	registerAuthAPI(appContainer, cfg, api)

	return router.Listen(fmt.Sprintf(":%d", cfg.HttpPort))
}

func registerAuthAPI(appContainer app.App, cfg config.ServerConfig, router fiber.Router) {
	userSvcGetter := userServiceGetter(appContainer, cfg)
	router.Post("/sign-up", middlewares.SetTransaction(appContainer.DB()), SignUp(userSvcGetter))
	router.Get("/send-otp", middlewares.SetTransaction(appContainer.DB()), SendSignInOTP(userSvcGetter))
	router.Post("/sign-in", middlewares.SetTransaction(appContainer.DB()), SignIn(userSvcGetter))
	router.Get("/test", middlewares.NewAuthMiddleware([]byte(cfg.Secret)), TestHandler)
}
