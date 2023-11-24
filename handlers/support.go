package handlers

import (
	"errors"
	"net/http"
	"qexchange/models"
	"qexchange/services"
	"strconv"

	"github.com/labstack/echo/v4"
)

type TicketRequest struct {
	Msg     string `json:"message"`
	Subject string `json:"subject"`
	TradeId *uint  `json:"trade_id,omitempty"`
}

type MessageRequest struct {
	Msg      string `json:"message"`
	TicketID *uint  `json:"ticket_id"`
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
		ticketID := c.QueryParam("ticket_id")

		ticketIDStr, err := strconv.Atoi(ticketID)

		if ticketID == "" || err != nil {
			return c.JSON(http.StatusBadRequest, models.NewErrorRespone("", errors.New("wrong ticket_id")))
		}

		ticket, statusCode, err := service.GetTicketMessages(uint(ticketIDStr))
		if err != nil {
			response := models.NewErrorRespone("", err)
			return c.JSON(statusCode, response)
		}

		return c.JSON(statusCode, ticket)
	}
}

func GetAllTickets(service services.SupportService) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, bind := c.Get("user").(models.User)
		if !bind {
			response := models.NewErrorRespone("", errors.New("bad user data"))
			return c.JSON(http.StatusBadRequest, response)
		}

		tickets, statusCode, err := service.GetAllTickets(user)
		if err != nil {
			response := models.NewErrorRespone("", err)
			return c.JSON(statusCode, response)
		}

		return c.JSON(statusCode, tickets)
	}
}

func SendMessage(service services.SupportService) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := new(MessageRequest)
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, models.NewErrorRespone("", errors.New("invalid request format")))
		}

		// Validate that the subject and description are not empty
		if request.Msg == "" || request.TicketID == nil {
			response := models.NewErrorRespone("", errors.New("message and ticket_id are required"))
			return c.JSON(http.StatusBadRequest, response)
		}

		user, bind := c.Get("user").(models.User)
		if !bind {
			response := models.NewErrorRespone("", errors.New("bad user data"))
			return c.JSON(http.StatusBadRequest, response)
		}

		statusCode, err := service.SendMessage(user, request.Msg, *request.TicketID)
		if err != nil {
			response := models.NewErrorRespone("", err)
			return c.JSON(statusCode, response)
		}

		return c.JSON(statusCode, models.NewRespone("message sent successfully"))
	}
}

func CloseTicket(service services.SupportService) echo.HandlerFunc {
	return func(c echo.Context) error {
		ticketID := c.QueryParam("ticket_id")

		ticketIDStr, err := strconv.Atoi(ticketID)

		if ticketID == "" || err != nil {
			return c.JSON(http.StatusBadRequest, models.NewErrorRespone("", errors.New("wrong ticket_id")))
		}

		statusCode, err := service.CloseTicket(uint(ticketIDStr))
		if err != nil {
			response := models.NewErrorRespone("", err)
			return c.JSON(statusCode, response)
		}

		return c.JSON(statusCode, models.NewRespone("ticket closed successfully"))
	}
}
