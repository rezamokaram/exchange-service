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
		crypto cryptocurrency.MakeCryptoRequest,
	) (int, error) // returns statusCode, error

	UpdateCrypto(
		crypto cryptocurrency.UpdateCryptoRequest,
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
	req cryptocurrency.MakeCryptoRequest,
) (int, error) {
	if s.db.Where("name = ?", req.Name).First(&cryptocurrency.Crypto{}).Error == nil {
		return http.StatusBadRequest, errors.New("the crypto already exist")
	}

	crypto := req.ToCrypto()
	result := s.db.Save(&crypto)
	if result.Error != nil {
		return http.StatusInternalServerError, result.Error
	}

	return http.StatusOK, nil
}

func (s *cryptoService) UpdateCrypto(
	req cryptocurrency.UpdateCryptoRequest,
) (int, error) {
	var oldCrypto cryptocurrency.Crypto
	result := s.db.Where("id = ?", req.Id).First(&oldCrypto)
	if result.Error != nil {
		return http.StatusBadRequest, errors.New("there is no crypto with this id")
	}

	crypto := req.UpdateCrypto(oldCrypto)
	result = s.db.Save(&crypto)
	if result.Error != nil {
		return http.StatusBadRequest, result.Error
	}
	
	if crypto.CurrentPrice == oldCrypto.CurrentPrice {
		return http.StatusOK, nil
	}

	tradeService := NewTradeService(s.db)

	tradeService.CheckFutureOrder(oldCrypto, crypto)

	if oldCrypto.CurrentPrice > crypto.CurrentPrice {
		tradeService.CheckStopLoss(crypto)
	} else if oldCrypto.CurrentPrice < crypto.CurrentPrice {
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
