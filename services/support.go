package services

import (
	"errors"
	"net/http"
	"qexchange/models"
	"qexchange/models/trade"

	"gorm.io/gorm"
)

type SupportService interface {
	SendTicket(user models.User, subject, description string, tradeID *uint) (int, error)
}

type supportService struct {
	db *gorm.DB
}

func NewSupportService(db *gorm.DB) SupportService {
	return &supportService{
		db: db,
	}
}

func (s *supportService) SendTicket(user models.User, subject, description string, tradeID *uint) (int, error) {
	hasTradeID := false
	if tradeID != nil {
		hasTradeID = true // to know whether include the tradeId later or not
		var openTrade trade.OpenTrade
		if errors.Is(s.db.First(&openTrade, tradeID).Error, gorm.ErrRecordNotFound) {
			return http.StatusBadRequest, errors.New("wrong trade id")
		}
	}

	var newSupportTicket models.SupportTicket
	newSupportTicket.UserID = user.ID
	newSupportTicket.Subject = subject
	newSupportTicket.Description = description
	if hasTradeID {
		newSupportTicket.TradeId = tradeID
	}

	if s.db.Save(&newSupportTicket).Error != nil {
		return http.StatusBadRequest, errors.New("failed saving the ticket in database")
	}

	return http.StatusOK, nil
}
