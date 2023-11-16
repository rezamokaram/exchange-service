package services

import (
	"errors"
	"net/http"

	"qexchange/models"
	"qexchange/models/trade"

	"gorm.io/gorm"
)

type TradeService interface {
	OpenTrade(
		request trade.OpenTradeRequest,
		user models.User,
	) (int, error)

	CloseTrade(
		request trade.ClosedTradeRequest,
		user models.User,
	) (int, error)

	GetAllClosedTrades(
		user models.User,
	) ([]trade.ClosedTrade, int, error)

	GetAllOpenTrades(
		user models.User,
	) ([]trade.OpenTrade, int, error)
}

type tradeService struct {
	db *gorm.DB
}

func NewTradeService(db *gorm.DB) TradeService {
	return &tradeService{
		db: db,
	}
}

func (s *tradeService) OpenTrade(
	request trade.OpenTradeRequest,
	user models.User,
) (int, error) {
	var crypto models.Crypto
	result := s.db.Where("id = ?", request.CryptoID).First(&crypto)
	if result.Error != nil {
		return http.StatusBadRequest, errors.New("there is no crypto with this id")
	}
	
	cost := request.Amount * float64(crypto.BuyFee)
	if cost > float64(user.Profile.Balance) {
		return http.StatusBadRequest, errors.New("you do not have enough money in your account")
	}

	// authorization level

	// start of opening trade
	transaction := request.ToTransaction(user.ID, crypto.BuyFee)
	result = s.db.Save(&transaction)
	if result.Error != nil {
		return http.StatusInternalServerError, errors.New("database error")
	}
	
	newTrade := request.ToOpenTrade(user.ID, crypto.BuyFee)
	result = s.db.Save(&newTrade)
	if result.Error != nil {
		return http.StatusInternalServerError, errors.New("database error")
	}

	// update balance

	return http.StatusOK, nil
}

func (s *tradeService) CloseTrade(
	request trade.ClosedTradeRequest,
	user 	models.User,
) (int, error) {
	var openTrade trade.OpenTrade
	result := s.db.Where("id = ?", request.OpenTradeID).First(&openTrade)
	if result.Error != nil {
		return http.StatusBadRequest, errors.New("there is no trade with this id")
	}

	if openTrade.UserID != user.ID {
		return http.StatusBadRequest, errors.New("this trade belong to another user")
	}

	var crypto models.Crypto
	result = s.db.Where("id = ?", openTrade.CryptoID).First(&crypto) 
	if result != nil {
		return http.StatusInternalServerError, errors.New("data base error")
	}

	if request.Amount > openTrade.Amount {
		return http.StatusBadRequest, errors.New("requested amount is too large")
	}

	if openTrade.Amount == request.Amount {
		s.db.Delete(&openTrade)
	} else {
		openTrade.Amount -= request.Amount
		s.db.Save(&openTrade)
	}

	newClosedTrade := openTrade.ToCloseTrade(crypto.SellFee)
	s.db.Save(&newClosedTrade)
	
	// update balance

	return http.StatusOK, nil
}

func (s *tradeService) GetAllClosedTrades(
	user models.User,
) ([]trade.ClosedTrade ,int, error) {
	var allClosedTrades []trade.ClosedTrade
	result := s.db.Where("user_id = ?", user.ID).Find(&allClosedTrades)
	if result != nil {
		return make([]trade.ClosedTrade,0), http.StatusInternalServerError, errors.New("data base error")
	}

	return allClosedTrades, http.StatusAccepted, nil
}

func (s *tradeService) GetAllOpenTrades(
	user models.User,
) ([]trade.OpenTrade ,int, error) {
	var allOpenTrades []trade.OpenTrade
	result := s.db.Where("user_id = ?", user.ID).Find(&allOpenTrades)
	if result != nil {
		return make([]trade.OpenTrade,0), http.StatusInternalServerError, errors.New("data base error")
	}

	return allOpenTrades, http.StatusAccepted, nil
}