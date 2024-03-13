package services

import (
	"errors"
	"net/http"
	"time"

	userModels "qexchange/models/user"
	"qexchange/models"
	"qexchange/models/trade"

	"gorm.io/gorm"
)

type SupportService interface {
	OpenTicket(user userModels.User, subject, ticketMsg string, tradeID *uint) (int, error)
	SendMessage(user userModels.User, message string, ticketID uint) (int, error)
	GetActiveTickets() ([]models.SupportTicket, int, error)
	GetAllTickets(user userModels.User) ([]models.SupportTicket, int, error)
	GetTicketMessages(ticketID uint) (models.SupportTicket, int, error)
	CloseTicket(ticketID uint) (int, error)
}

type supportService struct {
	db *gorm.DB
}

func NewSupportService(db *gorm.DB) SupportService {
	return &supportService{
		db: db,
	}
}

func (s *supportService) OpenTicket(user userModels.User, subject, ticketMsg string, tradeID *uint) (int, error) {
	hasTradeID := false
	if tradeID != nil {
		hasTradeID = true // to know whether include the tradeId later or not
		var openTrade trade.OpenTrade
		if errors.Is(s.db.First(&openTrade, tradeID).Error, gorm.ErrRecordNotFound) {
			return http.StatusBadRequest, errors.New("wrong trade id")
		}
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var newSupportTicket models.SupportTicket
		newSupportTicket.UserID = user.ID
		newSupportTicket.Username = user.Username
		newSupportTicket.Subject = subject
		newSupportTicket.CreatedAt = time.Now()
		newSupportTicket.UpdatedAt = time.Now()

		var tMsg models.TicketMessage
		tMsg.Msg = ticketMsg
		tMsg.SenderUsername = user.Username
		tMsg.CreatedAt = time.Now()

		// Save the new SupportTicket first to get its ID
		if err := s.db.Save(&newSupportTicket).Error; err != nil {
			return err
		}

		tMsg.SupportTicketID = newSupportTicket.ID

		if hasTradeID {
			newSupportTicket.TradeId = tradeID
		}

		if err := s.db.Save(&tMsg).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return http.StatusInternalServerError, errors.New("failed saving the ticket in database")
	}

	return http.StatusOK, nil
}

func (s *supportService) SendMessage(user userModels.User, message string, ticketID uint) (int, error) {
	var ticket models.SupportTicket
	if s.db.Where("id = ?", ticketID).First(&ticket).Error != nil {
		return http.StatusBadRequest, errors.New("wrong ticket_id")
	}

	var ticketMessage models.TicketMessage
	ticketMessage.Msg = message
	ticketMessage.SenderUsername = user.Username
	ticketMessage.SupportTicketID = ticketID

	if s.db.Save(&ticketMessage).Error != nil {
		return http.StatusInternalServerError, errors.New("failed saving the message in database")
	}

	return http.StatusOK, nil
}

func (s *supportService) GetActiveTickets() ([]models.SupportTicket, int, error) {
	var tickets []models.SupportTicket

	if s.db.Where("status IN (?)", []int{models.OpenTicket, models.PendingTicket}).Omit("Messages").Find(&tickets).Error != nil {
		return nil, http.StatusInternalServerError, errors.New("could not get the tickets")
	}

	return tickets, http.StatusOK, nil
}

func (s *supportService) GetAllTickets(user userModels.User) ([]models.SupportTicket, int, error) {
	var tickets []models.SupportTicket

	if s.db.Where("user_id = ?", user.ID).Omit("Messages").Find(&tickets).Error != nil {
		return nil, http.StatusInternalServerError, errors.New("could not get the tickets")
	}

	return tickets, http.StatusOK, nil
}

func (s *supportService) GetTicketMessages(ticketID uint) (models.SupportTicket, int, error) {
	var ticket models.SupportTicket

	if s.db.Where("id = ?", ticketID).Preload("Messages").First(&ticket).Error != nil {
		return models.SupportTicket{}, http.StatusInternalServerError, errors.New("could not get the tickets")
	}

	return ticket, http.StatusOK, nil
}

func (s *supportService) CloseTicket(ticketID uint) (int, error) {
	var ticket models.SupportTicket

	if s.db.Where("id = ?", ticketID).Preload("Messages").First(&ticket).Error != nil {
		return http.StatusInternalServerError, errors.New("could not get the tickets")
	}

	if ticket.Status == models.ClosedTicket {
		return http.StatusBadRequest, errors.New("ticket already closed")
	}

	ticket.Status = models.ClosedTicket

	if s.db.Save(&ticket).Error != nil {
		return http.StatusInternalServerError, errors.New("failed updating ticket in database")
	}

	return http.StatusOK, nil
}
