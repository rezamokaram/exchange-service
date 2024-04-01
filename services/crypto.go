package services

import (
	"errors"
	"net/http"

	cryptoModels "qexchange/models/crypto"

	"gorm.io/gorm"
)

type CryptoService interface {
	GetCrypto(
		id int,
	) (cryptoModels.CryptoResponse, int, error) // returns answer, statusCode, error

	SetCrypto(
		crypto cryptoModels.MakeCryptoRequest,
	) (int, error) // returns statusCode, error

	UpdateCrypto(
		crypto cryptoModels.UpdateCryptoRequest,
	) (int, error) // returns statusCode, error

	GetAllCrypto() ([]cryptoModels.CryptoResponse, int, error)
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
) (cryptoModels.CryptoResponse, int, error) {
	var crypto cryptoModels.Crypto
	result := s.db.Where("id = ?", id).First(&crypto)
	if result.Error != nil {
		return cryptoModels.CryptoResponse{}, http.StatusBadRequest, errors.New("there is no crypto with this id")
	}

	return cryptoModels.NewCryptoResponse(crypto), http.StatusOK, nil
}

func (s *cryptoService) SetCrypto(
	req cryptoModels.MakeCryptoRequest,
) (int, error) {
	if s.db.Where("name = ?", req.Name).First(&cryptoModels.Crypto{}).Error == nil {
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
	req cryptoModels.UpdateCryptoRequest,
) (int, error) {
	var oldCrypto cryptoModels.Crypto
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

func (s *cryptoService) GetAllCrypto() ([]cryptoModels.CryptoResponse, int, error) {
	var cryptoList []cryptoModels.Crypto
	result := s.db.Find(&cryptoList)
	if result.Error != nil {
		return make([]cryptoModels.CryptoResponse, 0), http.StatusBadRequest, errors.New("there is no crypto with this id")
	}

	cryptoResponseList := make([]cryptoModels.CryptoResponse, 0)
	for _, cr := range cryptoList {
		cryptoResponseList = append(cryptoResponseList, cryptoModels.NewCryptoResponse(cr))
	}
	return cryptoResponseList, http.StatusOK, nil
}
