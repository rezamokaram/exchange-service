package user

import (
	"github.com/RezaMokaram/ExchangeService/api/handlers/http/middlewares"
	"github.com/RezaMokaram/ExchangeService/app"
	"github.com/RezaMokaram/ExchangeService/config"

	"github.com/gofiber/fiber/v2"
)

func RegisterAuthAPI(appContainer app.App, cfg config.ServerConfig, router fiber.Router) {
	userSvcGetter := userServiceGetter(appContainer, cfg)
	router.Post("/sign-up", middlewares.SetTransaction(appContainer.DB()), SignUp(userSvcGetter))
	router.Get("/send-otp", middlewares.SetTransaction(appContainer.DB()), SendSignInOTP(userSvcGetter))
	router.Post("/sign-in", middlewares.SetTransaction(appContainer.DB()), SignIn(userSvcGetter))
	router.Get("/test", middlewares.NewAuthMiddleware([]byte(cfg.Secret)), TestHandler)
}
