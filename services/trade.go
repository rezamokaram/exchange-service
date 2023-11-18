package services

import (
	"errors"
	"net/http"
	"strconv"

	"qexchange/models"
	"qexchange/models/trade"
	// "qexchange/services"

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

	var profile models.Profile
	result = s.db.Where("id = ?", user.ID).First(&profile)
	if result.Error != nil {
		return http.StatusBadRequest, errors.New("there is no profile with this id")
	}

	cost := request.Amount * float64(crypto.BuyFee)
	if cost > float64(profile.Balance) {
		return http.StatusBadRequest, errors.New("you do not have enough money in your account" + strconv.Itoa(profile.Balance))
	}

	bankService := NewBankService(s.db)
	statusCode, err := bankService.SubtractFromUserBalanace(user, int(cost))
	if err != nil {
		return statusCode, errors.New("error in banking operations")
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
	if result.Error != nil {
		return http.StatusInternalServerError, errors.New("database error")
	}

	if request.Amount > openTrade.Amount {
		return http.StatusBadRequest, errors.New("requested amount is too large")
	}



	if openTrade.Amount == request.Amount {
		// s.db.Delete(&openTrade)
		result = s.db.Exec("DELETE FROM open_trade WHERE id = ?", openTrade.ID)
		if result.Error != nil {
			return http.StatusBadRequest, errors.New("requested amount is too large")
		}

	} else {
		openTrade.Amount -= request.Amount
		s.db.Save(&openTrade)
	}

	cost := request.Amount * float64(crypto.SellFee)
	bankService := NewBankService(s.db)
	statusCode, err := bankService.AddToUserBalanace(user, int(cost))
	if err != nil {
		return statusCode, errors.New("error in banking operations")
	}

	newClosedTrade := openTrade.ToCloseTrade(crypto.SellFee)
	s.db.Save(&newClosedTrade)

	return http.StatusOK, nil
}

func (s *tradeService) GetAllClosedTrades(
	user models.User,
) ([]trade.ClosedTrade ,int, error) {
	var allClosedTrades []trade.ClosedTrade
	result := s.db.Where("user_id = ?", user.ID).Find(&allClosedTrades)
	if result.Error != nil {
		return make([]trade.ClosedTrade,0), http.StatusInternalServerError, result.Error//errors.New("data base error")
	}
	return allClosedTrades, http.StatusAccepted, nil
}

func (s *tradeService) GetAllOpenTrades(
	user models.User,
) ([]trade.OpenTrade ,int, error) {
	var allOpenTrades []trade.OpenTrade
	result := s.db.Where("user_id = ?", user.ID).Find(&allOpenTrades)
	if result.Error != nil {
		return make([]trade.OpenTrade,0), http.StatusInternalServerError, result.Error//errors.New("data base error")
	}

	return allOpenTrades, http.StatusAccepted, nil
}