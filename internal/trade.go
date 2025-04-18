package internal

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	cryptoModels "github.com/rezamokaram/exchange-service/models/crypto"
	tradeModels "github.com/rezamokaram/exchange-service/models/trade"
	userModels "github.com/rezamokaram/exchange-service/models/user"

	"gorm.io/gorm"
)

type TradeService interface {
	OpenTrade(
		request tradeModels.OpenTradeRequest,
		user userModels.User,
	) (int, error)

	CloseTrade(
		request tradeModels.ClosedTradeRequest,
		user userModels.User,
	) (int, error)

	CloseTradeWithTrade(
		openTrade tradeModels.OpenTrade,
		user userModels.User,
		crypto cryptoModels.Crypto,
		amount int,
	) (int, error)

	GetAllClosedTrades(
		user userModels.User,
	) ([]tradeModels.ClosedTrade, int, error)

	GetAllOpenTrades(
		user userModels.User,
	) ([]tradeModels.OpenTrade, int, error)

	CheckStopLoss(
		crypto cryptoModels.Crypto,
	)

	CheckTakeProfit(
		crypto cryptoModels.Crypto,
	)

	SetFutureOrder(
		req tradeModels.FutureOrderRequest,
		user userModels.User,
	) (int, error)

	DeleteFutureOrder(
		req tradeModels.DeleteFutureOrderRequest,
		user userModels.User,
	) (int, error)

	CheckFutureOrder(
		oldCrypto cryptoModels.Crypto,
		newCrypto cryptoModels.Crypto,
	)

	GetAllFutureOrders(
		user userModels.User,
	) ([]tradeModels.FutureOrder, int, error)

	FilterClosedTrades(
		user userModels.User,
		req tradeModels.FilterTradesRequest,
	) (tradeModels.FilterTradesResponse, int, error)
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
	request tradeModels.OpenTradeRequest,
	user userModels.User,
) (int, error) {
	var crypto cryptoModels.Crypto
	result := s.db.Where("id = ?", request.CryptoID).First(&crypto)
	if result.Error != nil {
		return http.StatusBadRequest, errors.New("there is no crypto with this id")
	}

	var profile userModels.Profile
	result = s.db.Where("id = ?", user.ID).First(&profile)
	if result.Error != nil {
		return http.StatusBadRequest, errors.New("there is no profile with this id")
	}

	cost := request.Amount * crypto.BuyFee
	if cost > profile.Balance {
		return http.StatusBadRequest, errors.New("you do not have enough money in your account" + strconv.Itoa(profile.Balance))
	}

	bankService := NewBankService(s.db)
	description := fmt.Sprintf("Trade Service: for opening a trade, crypto = %v with crypto id = %v and amount = %v at %v", crypto.Name, crypto.ID, request.Amount, time.Now())
	statusCode, err := bankService.SubtractFromUserBalance(user, int(cost), 1, description)
	if err != nil {
		return statusCode, err
	}

	newTrade := request.ToOpenTrade(user.ID, crypto.BuyFee)
	result = s.db.Save(&newTrade)
	if result.Error != nil {
		return http.StatusInternalServerError, errors.New("database error")
	}

	return http.StatusOK, nil
}

func (s *tradeService) CloseTrade(
	request tradeModels.ClosedTradeRequest,
	user userModels.User,
) (int, error) {
	var openTrade tradeModels.OpenTrade
	result := s.db.Where("id = ?", request.OpenTradeID).First(&openTrade)
	if result.Error != nil {
		return http.StatusBadRequest, errors.New("there is no trade with this id")
	}

	if openTrade.UserID != user.ID {
		return http.StatusBadRequest, errors.New("this trade belong to another user")
	}

	var crypto cryptoModels.Crypto
	result = s.db.Where("id = ?", openTrade.CryptoID).First(&crypto)
	if result.Error != nil {
		return http.StatusInternalServerError, errors.New("database error")
	}

	return s.CloseTradeWithTrade(openTrade, user, crypto, request.Amount)
}

func (s *tradeService) CloseTradeWithTrade(
	openTrade tradeModels.OpenTrade,
	user userModels.User,
	crypto cryptoModels.Crypto,
	amount int,
) (int, error) {

	if openTrade.UserID != user.ID {
		return http.StatusBadRequest, errors.New("this trade belong to another user")
	}

	if amount > openTrade.Amount {
		return http.StatusBadRequest, errors.New("requested amount is too large")
	}

	if openTrade.Amount == amount {
		result := s.db.Exec("DELETE FROM open_trade WHERE id = ?", openTrade.ID)
		if result.Error != nil {
			return http.StatusBadRequest, errors.New("requested amount is too large")
		}

	} else {
		openTrade.Amount -= amount
		s.db.Save(&openTrade)
	}

	cost := amount * crypto.SellFee
	bankService := NewBankService(s.db)
	description := fmt.Sprintf("Trade Service: for closing a trade, crypto = %v with crypto id = %v and amount = %v at %v", crypto.Name, crypto.ID, openTrade.Amount, time.Now())
	statusCode, err := bankService.AddToUserBalance(user, cost, 1, description)
	if err != nil {
		return statusCode, errors.New("error in banking operations")
	}

	newClosedTrade := openTrade.ToCloseTrade(crypto.SellFee, amount)
	s.db.Save(&newClosedTrade)

	return http.StatusOK, nil
}

func (s *tradeService) GetAllClosedTrades(
	user userModels.User,
) ([]tradeModels.ClosedTrade, int, error) {
	var allClosedTrades []tradeModels.ClosedTrade
	result := s.db.Where("user_id = ?", user.ID).Find(&allClosedTrades)
	if result.Error != nil {
		return make([]tradeModels.ClosedTrade, 0), http.StatusInternalServerError, result.Error //errors.New("data base error")
	}
	return allClosedTrades, http.StatusOK, nil
}

func (s *tradeService) GetAllOpenTrades(
	user userModels.User,
) ([]tradeModels.OpenTrade, int, error) {
	var allOpenTrades []tradeModels.OpenTrade
	result := s.db.Where("user_id = ?", user.ID).Find(&allOpenTrades)
	if result.Error != nil {
		return make([]tradeModels.OpenTrade, 0), http.StatusInternalServerError, result.Error //errors.New("data base error")
	}

	return allOpenTrades, http.StatusOK, nil
}

func (s *tradeService) CheckStopLoss(
	crypto cryptoModels.Crypto,
) {
	fmt.Println("Stop Loss Processing ...")
	var allTriggeredTrades []tradeModels.OpenTrade
	result := s.db.Where("stop_loss >= ?", crypto.SellFee).Find(&allTriggeredTrades)
	if result.Error != nil {
		return
	}

	fmt.Println(len(allTriggeredTrades), " trade detected for closing ...")
	var wg sync.WaitGroup
	for _, triggeredTrade := range allTriggeredTrades {
		wg.Add(1)
		go func(s *tradeService, toCloseTrade tradeModels.OpenTrade) {
			defer wg.Done()
			var user userModels.User
			res := s.db.Where("id = ?", toCloseTrade.UserID).First(&user)
			if res.Error != nil {
				return
			}
			s.CloseTradeWithTrade(toCloseTrade, user, crypto, toCloseTrade.Amount)
		}(s, triggeredTrade)
	}
	wg.Wait()
}

func (s *tradeService) CheckTakeProfit(
	crypto cryptoModels.Crypto,
) {
	fmt.Println("Take Profit Processing ...")
	var allTriggeredTrades []tradeModels.OpenTrade
	result := s.db.Where("take_profit <= ? and take_profit > 0", crypto.SellFee).Find(&allTriggeredTrades)
	if result.Error != nil {
		return
	}

	fmt.Println(len(allTriggeredTrades), " trade detected for closing ...")
	var wg sync.WaitGroup
	for _, triggeredTrade := range allTriggeredTrades {
		wg.Add(1)
		go func(s *tradeService, toCloseTrade tradeModels.OpenTrade) {
			defer wg.Done()
			var user userModels.User
			res := s.db.Where("id = ?", toCloseTrade.UserID).First(&user)
			if res.Error != nil {
				return
			}
			s.CloseTradeWithTrade(toCloseTrade, user, crypto, toCloseTrade.Amount)
		}(s, triggeredTrade)
	}
	wg.Wait()
}

func (s *tradeService) SetFutureOrder(
	req tradeModels.FutureOrderRequest,
	user userModels.User,
) (int, error) {
	futureOrder := req.ToFutureOrder(user.ID)
	result := s.db.Save(&futureOrder)
	if result.Error != nil {
		return http.StatusInternalServerError, result.Error
	}

	return http.StatusOK, nil
}

func (s *tradeService) DeleteFutureOrder(
	req tradeModels.DeleteFutureOrderRequest,
	user userModels.User,
) (int, error) {
	var futureOrder tradeModels.FutureOrder
	result := s.db.Where("id = ?", req.OrderID).First(&futureOrder)
	if result.Error != nil {
		return http.StatusInternalServerError, result.Error
	}

	if futureOrder.UserID != user.ID {
		return http.StatusForbidden, errors.New("this order belong to another user")
	}

	result = s.db.Exec("DELETE FROM future_order WHERE id = ?", futureOrder.ID)
	if result.Error != nil {
		return http.StatusBadRequest, result.Error
	}

	return http.StatusOK, nil
}

func (s *tradeService) CheckFutureOrder(
	oldCrypto cryptoModels.Crypto,
	newCrypto cryptoModels.Crypto,
) {
	fmt.Printf("checking for future orders...\n")
	if newCrypto.CurrentPrice == oldCrypto.CurrentPrice {
		fmt.Printf("there is nothing to do!")
		return
	}
	allTriggeredFutureOrders := make([]tradeModels.FutureOrder, 0)
	if newCrypto.CurrentPrice > oldCrypto.CurrentPrice {
		res := s.db.Where("target_price >= ? and target_price <= ?", oldCrypto.BuyFee, newCrypto.BuyFee).Find(&allTriggeredFutureOrders)
		if res.Error != nil {
			fmt.Printf("server internal error: %v", res.Error)
		}
	} else {
		res := s.db.Where("target_price >= ? and target_price <= ?", newCrypto.BuyFee, oldCrypto.BuyFee).Find(&allTriggeredFutureOrders)
		if res.Error != nil {
			fmt.Printf("server internal error: %v", res.Error)
		}
	}

	var wg sync.WaitGroup
	for _, toOpenFutureOrder := range allTriggeredFutureOrders {
		wg.Add(1)
		go func(s *tradeService, userID uint, fo tradeModels.FutureOrder) {
			defer wg.Done()
			toOpenTradeRequest := fo.ToOpenTradeRequest()
			var user userModels.User
			result := s.db.Where("id = ?", userID).First(&user)
			if result.Error != nil {
				// this error handling is a bad solution
				// in fact we need to send otp
				fmt.Printf("error while finding user in database")
			}

			status, err := s.OpenTrade(toOpenTradeRequest, user)
			if err != nil {
				fmt.Printf("the process failed with status code %v and error message %v \n", status, err.Error())
			}

			result = s.db.Exec("DELETE FROM future_order WHERE id = ?", fo.ID)
			if result.Error != nil {
				fmt.Printf("error in delete from future order table\n")
				// i don`t know what should we do
				// probably we need to repeat query
			}

		}(s, toOpenFutureOrder.UserID, toOpenFutureOrder)
	}
	wg.Wait()
}

func (s *tradeService) GetAllFutureOrders(
	user userModels.User,
) ([]tradeModels.FutureOrder, int, error) {
	var allFutureOrders []tradeModels.FutureOrder
	result := s.db.Where("user_id = ?", user.ID).Find(&allFutureOrders)
	if result.Error != nil {
		return make([]tradeModels.FutureOrder, 0), http.StatusInternalServerError, result.Error
	}

	return allFutureOrders, http.StatusOK, nil
}

func (s *tradeService) FilterClosedTrades(
	user userModels.User,
	req tradeModels.FilterTradesRequest,
) (tradeModels.FilterTradesResponse, int, error) {
	if req.End.IsZero() {
		req.End = time.Now().Add(time.Hour)
	}

	var filterResponse tradeModels.FilterTradesResponse
	var trades []tradeModels.ClosedTrade
	var result *gorm.DB
	if len(req.CryptoList) == 0 {
		result = s.db.Where("created_at >= ? AND created_at <= ?", req.Start, req.End).Find(&trades)
	} else {
		result = s.db.Where("created_at BETWEEN ? AND ? AND crypto_id IN (?)", req.Start, req.End, req.CryptoList).Find(&trades)
	}
	if result.Error != nil {
		return filterResponse, http.StatusInternalServerError, result.Error
	}

	filterResponse.Start = req.Start
	filterResponse.End = req.End
	filterResponse.CryptoList = req.CryptoList
	for _, tr := range trades {
		filterResponse.ProfitOverAll += tr.Profit
	}
	filterResponse.ClosedTrades = trades

	return filterResponse, http.StatusOK, nil
}
