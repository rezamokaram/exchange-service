package handlers

import (
	"errors"
	"net/http"
	"qexchange/models"
	"qexchange/services"

	"github.com/labstack/echo/v4"
)

type TicketRequest struct {
	Msg     string `json:"message"`
	Subject string `json:"subject"`
	TradeId *uint  `json:"trade_id,omitempty"`
}

func OpenTicket(service services.SupportService) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := new(TicketRequest)
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, models.NewErrorRespone("", errors.New("invalid request format")))
		}

		// Validate that the subject and description are not empty
		if request.Subject == "" || request.Msg == "" {
			response := models.NewErrorRespone("", errors.New("subject and message are required"))
			return c.JSON(http.StatusBadRequest, response)
		}

		user, bind := c.Get("user").(models.User)
		if !bind {
			response := models.NewErrorRespone("", errors.New("bad user data"))
			return c.JSON(http.StatusBadRequest, response)
		}

		statusCode, err := service.OpenTicket(user, request.Subject, request.Msg, request.TradeId)
		if err != nil {
			response := models.NewErrorRespone("", err)
			return c.JSON(statusCode, response)
		}

		return c.JSON(statusCode, models.NewRespone("ticket opened successfully"))
	}
}

func GetActiveTickets(service services.SupportService) echo.HandlerFunc {
	return func(c echo.Context) error {
		tickets, statusCode, err := service.GetActiveTickets()
		if err != nil {
			response := models.NewErrorRespone("", err)
			return c.JSON(statusCode, response)
		}

		return c.JSON(statusCode, tickets)
	}
}

func GetTicketMessages(service services.SupportService) echo.HandlerFunc {
	return func(c echo.Context) error {
		ticket, statusCode, err := service.GetTicketMessages()
		if err != nil {
			response := models.NewErrorRespone("", err)
			return c.JSON(statusCode, response)
		}

		return c.JSON(statusCode, ticket)
	}
}
