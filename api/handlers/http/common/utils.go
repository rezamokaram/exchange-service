package common

import (
	"context"

	"github.com/rezamokaram/exchange-service/pkg/jwt"

	"github.com/gofiber/fiber/v2"
	jwt2 "github.com/golang-jwt/jwt/v5"
)

func UserClaims(ctx *fiber.Ctx) *jwt.UserClaims {
	if u := ctx.Locals("user"); u != nil {
		userClaims, ok := u.(*jwt2.Token).Claims.(*jwt.UserClaims)
		if ok {
			return userClaims
		}
	}
	return nil
}

type ServiceGetter[T any] func(context.Context) T
