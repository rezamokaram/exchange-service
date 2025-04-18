package crypto

import (
	"github.com/rezamokaram/exchange-service/api/handlers/http/middlewares"
	"github.com/rezamokaram/exchange-service/app"
	"github.com/rezamokaram/exchange-service/config"

	"github.com/gofiber/fiber/v2"
)

func RegisterCryptoAPI(appContainer app.App, cfg config.ServerConfig, router fiber.Router) {
	userSvcGetter := cryptoServiceGetter(appContainer, cfg)
	router.Post("/crypto", middlewares.SetTransaction(appContainer.DB()), CreateCrypto(userSvcGetter))
	// router.Get("/send-otp", middlewares.SetTransaction(appContainer.DB()), SendSignInOTP(userSvcGetter))
	// router.Post("/sign-in", middlewares.SetTransaction(appContainer.DB()), SignIn(userSvcGetter))
	// router.Get("/test", middlewares.NewAuthMiddleware([]byte(cfg.Secret)), TestHandler)
}
