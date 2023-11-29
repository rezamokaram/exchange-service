package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"qexchange/handlers"
	"qexchange/models"
	"qexchange/server"
	"strconv"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSupportService(t *testing.T) {
	e := echo.New()
	server.UserRoutes(e, testDB)
	server.SupportRoutes(e, testDB)
	token := LoginAndGetToken(e, t, mockValidUser)
	adminToken := LoginAndGetToken(e, t, mockAdminUser)

	t.Run("Open new ticket", func(t *testing.T) {
		// Sub-test for "Invalid request format"
		t.Run("invalid request format", func(t *testing.T) {
			requestBody := []byte(`{invalidJson:}`)
			req := httptest.NewRequest(http.MethodPost, "/support/open-ticket", bytes.NewReader(requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set("Authorization", token)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code, "Expected status code to be 400 Bad Request")

			var errResp models.Response
			_ = json.NewDecoder(rec.Body).Decode(&errResp)
			assert.Contains(t, errResp.Message, "invalid request format", "Expected error message to contain 'invalid request format'")
		})

		// Sub-test for "Subject and message are required"
		t.Run("missing subject and message", func(t *testing.T) {
			requestBody, _ := json.Marshal(handlers.TicketRequest{
				Subject: "", // Empty subject
				Msg:     "", // Empty message
			})
			req := httptest.NewRequest(http.MethodPost, "/support/open-ticket", bytes.NewReader(requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set("Authorization", token)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code, "Expected status code to be 400 Bad Request")

			var errResp models.Response
			_ = json.NewDecoder(rec.Body).Decode(&errResp)
			assert.Contains(t, errResp.Message, "subject and message are required", "Expected error message to contain 'subject and message are required'")
		})

		// Sub-test for "Ticket opened successfully"
		t.Run("ticket opened successfully", func(t *testing.T) {
			requestBody, _ := json.Marshal(handlers.TicketRequest{
				Subject: "Valid Subject",
				Msg:     "Valid message content",
			})
			req := httptest.NewRequest(http.MethodPost, "/support/open-ticket", bytes.NewReader(requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set("Authorization", token)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code, "Expected status code to be 200 OK")

			var successResp models.Response
			_ = json.NewDecoder(rec.Body).Decode(&successResp)
			assert.Equal(t, "ticket opened successfully", successResp.Message, "Expected success message to match")

			// Check the database to ensure the ticket was actually created
			var ticket models.SupportTicket
			err := testDB.Where("subject = ? AND username = ?", "Valid Subject", mockValidUser.Username).First(&ticket).Error
			if assert.NoError(t, err) {
				assert.NotEqual(t, 0, ticket.ID, "Expected ticket ID to be non-zero")
				assert.Equal(t, "Valid Subject", ticket.Subject, "Expected ticket subject to match")
				assert.Equal(t, "user1", ticket.Username, "Expected ticket username to match")
			}
		})

		// Sub-test for "Wrong trade id"
		t.Run("wrong trade id", func(t *testing.T) {
			invalidTradeID := uint(99999) // Assuming 99999 is a non-existent trade ID
			requestBody, _ := json.Marshal(handlers.TicketRequest{
				Subject: "Subject With Invalid Trade ID",
				Msg:     "Message content",
				TradeId: &invalidTradeID,
			})
			req := httptest.NewRequest(http.MethodPost, "/support/open-ticket", bytes.NewReader(requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set("Authorization", token)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code, "Expected status code to be 400 Bad Request")

			var errResp models.Response
			_ = json.NewDecoder(rec.Body).Decode(&errResp)
			assert.Contains(t, errResp.Message, "wrong trade id", "Expected error message to contain 'wrong trade id'")
		})
	})

	t.Run("close existing ticket", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/support/close-ticket", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Authorization", token)
		req.URL.RawQuery = fmt.Sprintf("ticket_id=%d", 1)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Expected status code to be 200 OK")

		var closeResp models.Response
		_ = json.NewDecoder(rec.Body).Decode(&closeResp)
		assert.Equal(t, "ticket closed successfully", closeResp.Message, "Expected success message to match")

		// Check the database to ensure the ticket status is updated
		var updatedTicket models.SupportTicket
		err := testDB.Where("id = ?", 1).First(&updatedTicket).Error
		if assert.NoError(t, err) {
			assert.Equal(t, models.ClosedTicket, updatedTicket.Status, "Expected ticket status to be 'Closed'")
		}
	})

	t.Run("Send message to support ticket", func(t *testing.T) {
		newTicket := models.SupportTicket{
			UserID:    2,
			Username:  "user1",
			Subject:   "Ticket for Message",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := testDB.Create(&newTicket).Error; err != nil {
			t.Fatalf("Failed to create ticket: %v", err)
		}

		newMessage := handlers.MessageRequest{
			Msg:      "New message to the ticket",
			TicketID: &newTicket.ID,
		}
		requestBody, _ := json.Marshal(newMessage)
		req := httptest.NewRequest(http.MethodPost, "/support/send-message", bytes.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Authorization", token)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Expected status code to be 200 OK")

		var messageResp models.Response
		_ = json.NewDecoder(rec.Body).Decode(&messageResp)
		assert.Equal(t, "message sent successfully", messageResp.Message, "Expected success message to match")

		// Check the database to ensure the message was added to the ticket
		var ticketMessages []models.TicketMessage
		err := testDB.Where("support_ticket_id = ?", newTicket.ID).Find(&ticketMessages).Error
		if assert.NoError(t, err) {
			assert.Greater(t, len(ticketMessages), 0, "Expected at least one message in the ticket")
			assert.Equal(t, ticketMessages[0].Msg, "New message to the ticket", "Expected message content to match")
		}
	})

	t.Run("Get all messages for specific ticket", func(t *testing.T) {
		// Create a ticket
		newTicket := models.SupportTicket{
			UserID:    2,
			Username:  "user1",
			Subject:   "Ticket for Message",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := testDB.Create(&newTicket).Error; err != nil {
			t.Fatalf("Failed to create ticket: %v", err)
		}

		// Create messages for the ticket
		for i := 0; i < 3; i++ {
			msg := models.TicketMessage{
				SupportTicketID: newTicket.ID,
				Msg:             fmt.Sprintf("Message %d for ticket", i+1),
				SenderUsername:  "user1",
				CreatedAt:       time.Now(),
			}
			if err := testDB.Create(&msg).Error; err != nil {
				t.Fatalf("Failed to create message %d: %v", i+1, err)
			}
		}

		// Make a request to get ticket messages
		req := httptest.NewRequest(http.MethodGet, "/support/get-ticket-messages", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Authorization", token)
		req.URL.RawQuery = fmt.Sprintf("ticket_id=%d", newTicket.ID)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Expected status code to be 200 OK")

		var ticketWithMessages models.SupportTicket
		_ = json.NewDecoder(rec.Body).Decode(&ticketWithMessages)
		assert.Len(t, ticketWithMessages.Messages, 3, "Expected 3 messages in the response")
	})

	t.Run("Get active tickets (Admin only)", func(t *testing.T) {
		if err := ClearDatabaseTables(testDB); err != nil {
			t.Fatalf("Failed to clear database tables: %v", err)
		}

		// Create active tickets
		for i := 0; i < 3; i++ {
			activeTicket := models.SupportTicket{
				UserID:    2,
				Username:  "user2",
				Subject:   "Active Ticket " + strconv.Itoa(i+1),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if err := testDB.Create(&activeTicket).Error; err != nil {
				t.Fatalf("Failed to create active ticket %d: %v", i+1, err)
			}
		}

		// Ensure that the request is made with admin-level authorization
		req := httptest.NewRequest(http.MethodGet, "/support/admin/get-active-tickets", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Authorization", adminToken)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Expected status code to be 200 OK")

		var activeTickets []models.SupportTicket
		_ = json.NewDecoder(rec.Body).Decode(&activeTickets)
		assert.Len(t, activeTickets, 3, "Expected 3 active tickets in the response")
	})

}
