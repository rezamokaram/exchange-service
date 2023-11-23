package handlers

import (
	"errors"
	"net/http"
	"qexchange/models"
	"qexchange/services"

	"github.com/labstack/echo/v4"
)

type TicketRequest struct {
	Subject     string `json:"subject"`
	Description string `json:"description"`
	TradeId     *uint  `json:"trade_id,omitempty"`
}

func SendTicket(service services.SupportService) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := new(TicketRequest)
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, models.NewErrorRespone("", errors.New("invalid request format")))
		}

		// Validate that the subject and description are not empty
		if request.Subject == "" || request.Description == "" {
			response := models.NewErrorRespone("subject and description are required", nil)
			return c.JSON(http.StatusBadRequest, response)
		}

		user, bind := c.Get("user").(models.User)
		if !bind {
			response := models.NewErrorRespone("bad user data", nil)
			return c.JSON(http.StatusBadRequest, response)
		}

		statusCode, err := service.SendTicket(user, request.Subject, request.Description, request.TradeId)
		if err != nil {
			response := models.NewErrorRespone("", err)
			return c.JSON(statusCode, response)
		}

		return c.JSON(statusCode, models.NewRespone("ticket made successfully"))
	}
}
