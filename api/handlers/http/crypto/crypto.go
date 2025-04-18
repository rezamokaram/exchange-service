package crypto

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rezamokaram/exchange-service/api/handlers/http/common"
	"github.com/rezamokaram/exchange-service/api/pb"
	"github.com/rezamokaram/exchange-service/api/service"
)

func CreateCrypto(svcGetter common.ServiceGetter[*service.CryptoService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		svc := svcGetter(c.UserContext())
		var req pb.CreateCryptoRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		_, err := svc.CreateCrypto(c.UserContext(), &req)
		if err != nil {
			if errors.Is(err, service.ErrUserCreationValidation) {
				return fiber.NewError(fiber.StatusBadRequest, err.Error())
			}

			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		// return c.JSON(resp)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "internal error",
		})
	}
}
