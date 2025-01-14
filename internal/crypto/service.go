package crypto

import (
	"context"
	"errors"

	"github.com/RezaMokaram/ExchangeService/internal/crypto/domain"
	"github.com/RezaMokaram/ExchangeService/internal/crypto/port"
)

var (
	ErrCryptoOnCreate           = errors.New("error on creating new crypto")
	ErrCryptoCreationValidation = errors.New("validation failed")
	ErrCryptoNotFound           = errors.New("crypto not found")
)

type service struct {
	repo port.Repo
}

func NewCryptoService(r port.Repo) port.Service {
	return &service{
		repo: r,
	}
}

func (s *service) CreateCrypto(ctx context.Context, crypto domain.Crypto) (domain.CryptoID, error) {
	cid, err := s.repo.Create(ctx, crypto)
	if err != nil {
		return 0, err
	}

	return cid, nil
}

func (s *service) GetCryptoByFilter(
	ctx context.Context,
	filter *domain.CryptoFilter,
) (*domain.Crypto, error) {
	crypto, err := s.repo.GetByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	if crypto == nil {
		return nil, ErrCryptoNotFound
	}

	return crypto, nil
}

// func (s *service) UpdateCrypto(
// 	req cryptoModels.UpdateCryptoRequest,
// ) (int, error) {
// 	var oldCrypto cryptoModels.Crypto
// 	result := s.db.Where("id = ?", req.Id).First(&oldCrypto)
// 	if result.Error != nil {
// 		return http.StatusBadRequest, errors.New("there is no crypto with this id")
// 	}

// 	crypto := req.UpdateCrypto(oldCrypto)
// 	result = s.db.Save(&crypto)
// 	if result.Error != nil {
// 		return http.StatusBadRequest, result.Error
// 	}

// 	if crypto.CurrentPrice == oldCrypto.CurrentPrice {
// 		return http.StatusOK, nil
// 	}

// TODO: Implement trade service *
// tradeService := NewTradeService(s.db)

// tradeService.CheckFutureOrder(oldCrypto, crypto)

// if oldCrypto.CurrentPrice > crypto.CurrentPrice {
// 	tradeService.CheckStopLoss(crypto)
// } else if oldCrypto.CurrentPrice < crypto.CurrentPrice {
// 	tradeService.CheckTakeProfit(crypto)
// }

// 	return http.StatusOK, nil
// }

// func (s *service) GetAllCrypto() ([]cryptoModels.CryptoResponse, int, error) {
// 	var cryptoList []cryptoModels.Crypto
// 	result := s.db.Find(&cryptoList)
// 	if result.Error != nil {
// 		return make([]cryptoModels.CryptoResponse, 0), http.StatusBadRequest, errors.New("there is no crypto with this id")
// 	}

// 	cryptoResponseList := make([]cryptoModels.CryptoResponse, 0)
// 	for _, cr := range cryptoList {
// 		cryptoResponseList = append(cryptoResponseList, cryptoModels.NewCryptoResponse(cr))
// 	}
// 	return cryptoResponseList, http.StatusOK, nil
// }
