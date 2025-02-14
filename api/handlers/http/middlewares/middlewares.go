package middlewares

import (
	"log/slog"
	"time"

	"github.com/RezaMokaram/ExchangeService/api/handlers/http/common"
	"github.com/RezaMokaram/ExchangeService/pkg/context"
	"github.com/RezaMokaram/ExchangeService/pkg/jwt"
	"github.com/RezaMokaram/ExchangeService/pkg/logger"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NewAuthMiddleware(secret []byte) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:  jwtware.SigningKey{Key: secret},
		Claims:      &jwt.UserClaims{},
		TokenLookup: "header:Authorization",
		SuccessHandler: func(ctx *fiber.Ctx) error {
			userClaims := common.UserClaims(ctx)
			if userClaims == nil {
				return fiber.ErrUnauthorized
			}

			logger := context.GetLogger(ctx.UserContext())
			context.SetLogger(ctx.UserContext(), logger.With("user_id", userClaims.UserID))

			return ctx.Next()
		},
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		},
		AuthScheme: "Bearer",
	})
}

func SetUserContext(c *fiber.Ctx) error {
	c.SetUserContext(context.NewAppContext(c.UserContext(), context.WithLogger(logger.NewLogger())))
	return c.Next()
}

func SetTransaction(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tx := db.Begin()

		context.SetDB(c.UserContext(), tx, true)

		err := c.Next()

		if c.Response().StatusCode() >= 300 {
			return context.Rollback(c.UserContext())
		}

		if err := context.CommitOrRollback(c.UserContext(), true); err != nil {
			return err
		}

		return err
	}
}

func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Log Request Body
		var reqBody []byte
		if c.Request().Body() != nil {
			reqBody = c.Request().Body()
		}

		// Call next handler
		err := c.Next()

		// Log Response Body
		resBody := c.Response().Body()
		status := c.Response().StatusCode()

		slog.Info("HTTP Request",
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.String("query", c.OriginalURL()),
			slog.String("request_body", string(reqBody)),
			slog.String("response_body", string(resBody)),
			slog.Int("status", status),
			slog.String("remote_ip", c.IP()),
			slog.Duration("duration", time.Since(start)),
		)

		return err
	}
}
