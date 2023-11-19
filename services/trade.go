package services

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"qexchange/models/cryptocurrency"
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

	CheckStopLoss(
		crypto cryptocurrency.Crypto,
	)

	CheckTakeProfit(
		crypto cryptocurrency.Crypto,
	)

	SetFutureOrder(
		req 	trade.FutureOrderRequest,
		user 	models.User,
	) (int, error)

	CheckFutureOrder(
		oldCrypto cryptocurrency.Crypto,
		newCrypto cryptocurrency.Crypto,
	)
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
	var crypto cryptocurrency.Crypto
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

	var crypto cryptocurrency.Crypto
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
	return allClosedTrades, http.StatusOK, nil
}

func (s *tradeService) GetAllOpenTrades(
	user models.User,
) ([]trade.OpenTrade ,int, error) {
	var allOpenTrades []trade.OpenTrade
	result := s.db.Where("user_id = ?", user.ID).Find(&allOpenTrades)
	if result.Error != nil {
		return make([]trade.OpenTrade,0), http.StatusInternalServerError, result.Error//errors.New("data base error")
	}

	return allOpenTrades, http.StatusOK, nil
}

func (s *tradeService) CheckStopLoss(
	crypto cryptocurrency.Crypto,
) {
	fmt.Println("Stop Loss Processing ...")
	var allTriggeredTrades []trade.OpenTrade
	result := s.db.Where("stop_loss >= ?", crypto.SellFee).Find(&allTriggeredTrades)
	if result.Error != nil {
		return
	}


	fmt.Println(len(allTriggeredTrades), " trade detected for closing ...")
	var wg sync.WaitGroup
	for _, triggeredTrade := range allTriggeredTrades {
		wg.Add(1)
		go func (s *tradeService, toCloseTrade trade.OpenTrade) {
			defer wg.Done()
			var user models.User
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
	crypto cryptocurrency.Crypto,
) {
	fmt.Println("Take Profit Processing ...")
	var allTriggeredTrades []trade.OpenTrade
	result := s.db.Where("take_profit <= ? and take_profit > 0", crypto.SellFee).Find(&allTriggeredTrades)
	if result.Error != nil {
		return
	}


	fmt.Println(len(allTriggeredTrades), " trade detected for closing ...")
	var wg sync.WaitGroup
	for _, triggeredTrade := range allTriggeredTrades {
		wg.Add(1)
		go func (s *tradeService, toCloseTrade trade.OpenTrade) {
			defer wg.Done()
			var user models.User
			res := s.db.Where("id = ?", toCloseTrade.UserID).First(&user)
			if res.Error != nil {
				return
			}
			s.CloseTradeWithTrade(toCloseTrade, user, crypto, toCloseTrade.Amount)
		}(s, triggeredTrade)
	}
	wg.Wait()
}

func (s *tradeService) CloseTradeWithTrade( // faster
	openTrade trade.OpenTrade,
	user 	models.User,
	crypto cryptocurrency.Crypto,
	amount float64,
) (int, error) {

	if openTrade.Amount == amount {
		result := s.db.Exec("DELETE FROM open_trade WHERE id = ?", openTrade.ID)
		if result.Error != nil {
			return http.StatusBadRequest, errors.New("requested amount is too large")
		}

	} else {
		openTrade.Amount -= amount
		s.db.Save(&openTrade)
	}

	cost := amount * float64(crypto.SellFee)
	bankService := NewBankService(s.db)
	statusCode, err := bankService.AddToUserBalanace(user, int(cost))
	if err != nil {
		return statusCode, errors.New("error in banking operations")
	}

	newClosedTrade := openTrade.ToCloseTrade(crypto.SellFee)
	s.db.Save(&newClosedTrade)

	return http.StatusOK, nil
}

func (s *tradeService) SetFutureOrder(
	req trade.FutureOrderRequest,
	user 	models.User,
) (int, error) {
	futureOrder := req.ToFutureOrder(user.ID)
	result := s.db.Save(&futureOrder)
	if result.Error != nil {
		return http.StatusInternalServerError, result.Error
	}

	return http.StatusOK, nil
}

func (s *tradeService) CheckFutureOrder(
	oldCrypto cryptocurrency.Crypto,
	newCrypto cryptocurrency.Crypto,
) {
	fmt.Printf("checking for future orders...\n")
	if newCrypto.CurrentPrice == oldCrypto.CurrentPrice {
		fmt.Printf("there is nothing to do!")
		return
	}
	allTriggeredFutureOrders := make([]trade.FutureOrder, 0)
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
		go func (s *tradeService, userID uint, fo trade.FutureOrder)  {
			defer wg.Done()
			toOpenTradeRequest := fo.ToOpenTradeRequest()
			var user models.User
			result := s.db.Where("id = ?", userID).First(&user)
			if result.Error != nil {
				// this error handling is a bad solution
				// in fact we need to send otp
				fmt.Printf("error while finding user in database")
			}

			status,err := s.OpenTrade(toOpenTradeRequest, user)
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

// func (s *tradeService) OpenTradeWithFutureOrder(
// 	futureOrder trade.FutureOrder,
// 	crypto 		cryptocurrency.Crypto,
// ) (int, error) {
// 	var user models.User
// 	result := s.db.Where("id = ?", futureOrder.UserID).First(&user)
// 	if result.Error != nil {
// 		return http.StatusBadRequest, errors.New("there is no user with this id")
// 	}

// 	var profile models.Profile
// 	result = s.db.Where("id = ?", futureOrder.UserID).First(&profile)
// 	if result.Error != nil {
// 		return http.StatusBadRequest, errors.New("there is no profile with this id")
// 	}

// 	cost := futureOrder.Amount * float64(crypto.BuyFee)
// 	if cost > float64(profile.Balance) {
// 		return http.StatusBadRequest, errors.New("you do not have enough money in your account" + strconv.Itoa(profile.Balance)) // TODO
// 	}

// 	bankService := NewBankService(s.db)
// 	statusCode, err := bankService.SubtractFromUserBalanace(user, int(cost))
// 	if err != nil {
// 		return statusCode, errors.New("error in banking operations")
// 	}

// 	// start of opening trade
// 	transaction := futureOrder.ToTransaction(user.ID, crypto.BuyFee)
// 	result = s.db.Save(&transaction)
// 	if result.Error != nil {
// 		return http.StatusInternalServerError, errors.New("database error")
// 	}
	
// 	newTrade := futureOrder.ToOpenTrade(crypto)
// 	result = s.db.Save(&newTrade)
// 	if result.Error != nil {
// 		return http.StatusInternalServerError, errors.New("database error")
// 	}

// 	//
// 	result = s.db.Exec("DELETE FROM future_order WHERE id = ?", futureOrder.ID)
// 	if result.Error != nil {
// 		fmt.Printf("error in delete from future order table\n")
// 		// i don`t know what should we do
// 		// probably we need to repeat query
// 	}

// 	return http.StatusOK, nil
// }