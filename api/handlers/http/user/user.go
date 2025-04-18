package user

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/rezamokaram/exchange-service/api/handlers/http/common"
	"github.com/rezamokaram/exchange-service/api/pb"
	"github.com/rezamokaram/exchange-service/api/service"
	"github.com/rezamokaram/exchange-service/pkg/context"

	"github.com/gofiber/fiber/v2"
)

func SendSignInOTP(svcGetter common.ServiceGetter[*service.UserService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		svc := svcGetter(c.UserContext())
		phone := strings.TrimSpace(c.Query("phone"))

		if err := svc.SendSignInOTP(c.UserContext(), phone); err != nil {
			return err
		}

		return nil
	}
}

func SignUp(svcGetter common.ServiceGetter[*service.UserService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		svc := svcGetter(c.UserContext())
		var req pb.UserSignUpRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		resp, err := svc.SignUp(c.UserContext(), &req)
		if err != nil {
			if errors.Is(err, service.ErrUserCreationValidation) {
				return fiber.NewError(fiber.StatusBadRequest, err.Error())
			}

			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		// return c.JSON(resp)
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"data": resp,
		})
	}
}

func SignIn(svcGetter common.ServiceGetter[*service.UserService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		svc := svcGetter(c.UserContext())
		var req pb.UserSignInRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		resp, err := svc.SignIn(c.UserContext(), &req)
		if err != nil {
			if errors.Is(err, service.ErrUserNotFound) {
				return c.SendStatus(fiber.StatusNotFound)
			}

			if errors.Is(err, service.ErrInvalidUserPassword) {
				return fiber.NewError(fiber.StatusUnauthorized, err.Error())
			}

			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(resp)
	}
}

func TestHandler(ctx *fiber.Ctx) error {
	logger := context.GetLogger(ctx.UserContext())

	logger.Info("from test handler", "time", time.Now().Format(time.DateTime))

	return nil
}
