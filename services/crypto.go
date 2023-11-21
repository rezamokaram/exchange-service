package services

import (
	"errors"
	// "fmt"
	"net/http"
	// "qexchange/models"
	// "qexchange/utils"

	// "golang.org/x/crypto/bcrypt"

	"qexchange/models/cryptocurrency"

	"gorm.io/gorm"
)

type CryptoService interface {
	GetCrypto(
		id int,
	) (cryptocurrency.CryptoResponse, int, error) // returns answer, statusCode, error

	SetCrypto(
		crypto cryptocurrency.Crypto,
	) (int, error) // returns statusCode, error

	UpdateCrypto(
		crypto cryptocurrency.Crypto,
	) (int, error) // returns statusCode, error

	GetAllCrypto() ([]cryptocurrency.CryptoResponse, int, error)
}

type cryptoService struct {
	db *gorm.DB
}

func NewCryptoService(db *gorm.DB) CryptoService {
	return &cryptoService{
		db: db,
	}
}

func (s *cryptoService) GetCrypto(
	id int,
) (cryptocurrency.CryptoResponse, int, error) {
	var crypto cryptocurrency.Crypto
	result := s.db.Where("id = ?", id).First(&crypto)
	if result.Error != nil {
		return cryptocurrency.CryptoResponse{}, http.StatusBadRequest, errors.New("there is no crypto with this id")
	}

	return cryptocurrency.NewCryptoResponse(crypto), http.StatusOK, nil
}

func (s *cryptoService) SetCrypto(
	crypto cryptocurrency.Crypto,
) (int, error) {
	var cryptoSearch cryptocurrency.Crypto
	result := s.db.Where("name = ?", crypto.Name).First(&cryptoSearch)
	if result.Error == nil {
		return http.StatusBadRequest, errors.New("the crypto already exist")
	}
	crypto.BuyFee = crypto.CurrentPrice + (crypto.CurrentPrice / 100) + 10
	crypto.SellFee = crypto.CurrentPrice - ((crypto.CurrentPrice / 100) + 10)
	if crypto.SellFee < 0 {
		crypto.SellFee = 0
	}

	s.db.Save(&crypto)

	return http.StatusOK, nil
}

func (s *cryptoService) UpdateCrypto(
	crypto cryptocurrency.Crypto,
) (int, error) {
	var cryptoSearch cryptocurrency.Crypto
	result := s.db.Where("id = ?", crypto.ID).First(&cryptoSearch)
	if result.Error != nil {
		return http.StatusBadRequest, errors.New("there is no crypto with this id")
	}
	crypto.BuyFee = crypto.CurrentPrice + (crypto.CurrentPrice / 100) + 10
	crypto.SellFee = crypto.CurrentPrice - ((crypto.CurrentPrice / 100) + 10)
	if crypto.SellFee < 0 {
		crypto.SellFee = 0
	}

	s.db.Save(&crypto)

	tradeService := NewTradeService(s.db)

	tradeService.CheckFutureOrder(cryptoSearch, crypto)

	if cryptoSearch.CurrentPrice > crypto.CurrentPrice {
		tradeService.CheckStopLoss(crypto)
	} else if cryptoSearch.CurrentPrice < crypto.CurrentPrice {
		tradeService.CheckTakeProfit(crypto)
	}

	return http.StatusOK, nil
}

func (s *cryptoService) GetAllCrypto() ([]cryptocurrency.CryptoResponse, int, error) {
	var cryptoList []cryptocurrency.Crypto
	result := s.db.Find(&cryptoList)
	if result.Error != nil {
		return make([]cryptocurrency.CryptoResponse, 0), http.StatusBadRequest, errors.New("there is no crypto with this id")
	}

	cryptoResponseList := make([]cryptocurrency.CryptoResponse, 0)
	for _, cr := range cryptoList {
		cryptoResponseList = append(cryptoResponseList, cryptocurrency.NewCryptoResponse(cr))
	}
	return cryptoResponseList, http.StatusOK, nil
}
