package handlers

import (
	"net/http"
	"qexchange/models"
	userModels "qexchange/models/user"
	"qexchange/services"
	"strconv"

	"github.com/labstack/echo/v4"
)

// TicketRequest represents the request body for opening a support ticket
type TicketRequest struct {
	Msg     string `json:"message" example:"I have a problem"`
	Subject string `json:"subject" example:"problem with system"`
	TradeId *uint  `json:"trade_id,omitempty" example:"1"`
}

// MessageRequest represents the request body for sending a message to a support ticket
type MessageRequest struct {
	Msg      string `json:"message" example:"I have a problem"`
	TicketID *uint  `json:"ticket_id" example:"1"`
}

// OpenTicket handles opening a new support ticket
// @Summary Open a support ticket
// @Description Opens a new support ticket
// @Tags Support
// @Accept  json
// @Produce json
// @Param   body  body      TicketRequest  true  "Open Ticket"
// @Param Authorization header string true "Authorization token"
// @Security BasicAuth
// @Success 200   {object}  models.Response
// @Failure 400   {object}  models.Response
// @Router /support/open-ticket [post]
func OpenTicket(service services.SupportService) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := new(TicketRequest)
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, models.NewErrorResponse("invalid request format", err.Error()))
		}

		// Validate that the subject and description are not empty
		if request.Subject == "" || request.Msg == "" {
			response := models.NewErrorResponse("failed to open ticket", "subject and message are required")
			return c.JSON(http.StatusBadRequest, response)
		}

		user, bind := c.Get("user").(userModels.User)
		if !bind {
			response := models.NewErrorResponse("failed to open ticket", "bad user data")
			return c.JSON(http.StatusBadRequest, response)
		}

		statusCode, err := service.OpenTicket(user, request.Subject, request.Msg, request.TradeId)
		if err != nil {
			response := models.NewErrorResponse("failed to open ticket", err.Error())
			return c.JSON(statusCode, response)
		}

		return c.JSON(statusCode, models.NewResponse("ticket opened successfully"))
	}
}

// GetActiveTickets retrieves all active tickets (Admin only)
// @Summary Admin: Get active support tickets
// @Description Retrieves all active support tickets (Admin only)
// @Tags Support
// @Accept  json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Security BasicAuth
// @Security AdminAuth
// @Success 200   {object}  []models.SupportTicket
// @Failure 400   {object}  models.Response
// @Router /support/admin/get-active-tickets [get]
func GetActiveTickets(service services.SupportService) echo.HandlerFunc {
	return func(c echo.Context) error {
		tickets, statusCode, err := service.GetActiveTickets()
		if err != nil {
			response := models.NewErrorResponse("failed to get tickets", err.Error())
			return c.JSON(statusCode, response)
		}

		return c.JSON(statusCode, tickets)
	}
}

// GetTicketMessages retrieves messages for a specific support ticket
// @Summary Get messages for a support ticket
// @Description Retrieves all messages for a specific support ticket
// @Tags Support
// @Accept  json
// @Produce json
// @Param   ticket_id  query  int  true  "Ticket ID"
// @Param Authorization header string true "Authorization token"
// @Security BasicAuth
// @Success 200   {object}  models.SupportTicket
// @Failure 400   {object}  models.Response
// @Router /support/get-ticket-messages [get]
func GetTicketMessages(service services.SupportService) echo.HandlerFunc {
	return func(c echo.Context) error {
		ticketID := c.QueryParam("ticket_id")

		ticketIDStr, err := strconv.Atoi(ticketID)

		if ticketID == "" || err != nil {
			return c.JSON(http.StatusBadRequest, models.NewErrorResponse("wrong ticket_id", err.Error()))
		}

		ticket, statusCode, err := service.GetTicketMessages(uint(ticketIDStr))
		if err != nil {
			response := models.NewErrorResponse("failed to get ticket", err.Error())
			return c.JSON(statusCode, response)
		}

		return c.JSON(statusCode, ticket)
	}
}

// GetAllTickets retrieves all tickets for the authenticated user
// @Summary Get all tickets for a user
// @Description Retrieves all support tickets for the authenticated user
// @Tags Support
// @Accept  json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Security BasicAuth
// @Success 200   {object}  []models.SupportTicket
// @Failure 400   {object}  models.Response
// @Router /support/get-all-tickets [get]
func GetAllTickets(service services.SupportService) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, bind := c.Get("user").(userModels.User)
		if !bind {
			response := models.NewErrorResponse("failed to get tickets", "bad user data")
			return c.JSON(http.StatusBadRequest, response)
		}

		tickets, statusCode, err := service.GetAllTickets(user)
		if err != nil {
			response := models.NewErrorResponse("failed to get tickets", err.Error())
			return c.JSON(statusCode, response)
		}

		return c.JSON(statusCode, tickets)
	}
}

// SendMessage handles sending a message to an existing support ticket
// @Summary Send a message to a support ticket
// @Description Sends a message to a specific support ticket
// @Tags Support
// @Accept  json
// @Produce json
// @Param   body  body      MessageRequest  true  "Send Message"
// @Param Authorization header string true "Authorization token"
// @Security BasicAuth
// @Success 200   {object}  models.Response
// @Failure 400   {object}  models.Response
// @Router /support/send-message [post]
func SendMessage(service services.SupportService) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := new(MessageRequest)
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, models.NewErrorResponse("invalid request format", err.Error()))
		}

		// Validate that the subject and description are not empty
		if request.Msg == "" || request.TicketID == nil {
			response := models.NewErrorResponse("failed to send message", "message and ticket_id are required")
			return c.JSON(http.StatusBadRequest, response)
		}

		user, bind := c.Get("user").(userModels.User)
		if !bind {
			response := models.NewErrorResponse("failed to send message", "bad user data")
			return c.JSON(http.StatusBadRequest, response)
		}

		statusCode, err := service.SendMessage(user, request.Msg, *request.TicketID)
		if err != nil {
			response := models.NewErrorResponse("failed to send message", err.Error())
			return c.JSON(statusCode, response)
		}

		return c.JSON(statusCode, models.NewResponse("message sent successfully"))
	}
}

// CloseTicket handles closing an existing support ticket
// @Summary Close a support ticket
// @Description Closes an existing support ticket
// @Tags Support
// @Accept  json
// @Produce json
// @Param   ticket_id  query  int  true  "Ticket ID"
// @Param Authorization header string true "Authorization token"
// @Security BasicAuth
// @Success 200   {object}  models.Response
// @Failure 400   {object}  models.Response
// @Router /support/close-ticket [patch]
func CloseTicket(service services.SupportService) echo.HandlerFunc {
	return func(c echo.Context) error {
		ticketID := c.QueryParam("ticket_id")

		ticketIDStr, err := strconv.Atoi(ticketID)

		if ticketID == "" || err != nil {
			return c.JSON(http.StatusBadRequest, models.NewErrorResponse("wrong ticket_id", err.Error()))
		}

		statusCode, err := service.CloseTicket(uint(ticketIDStr))
		if err != nil {
			response := models.NewErrorResponse("failed to close the ticket", err.Error())
			return c.JSON(statusCode, response)
		}

		return c.JSON(statusCode, models.NewResponse("ticket closed successfully"))
	}
}
