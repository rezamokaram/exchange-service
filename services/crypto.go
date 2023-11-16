package services

import (
	"errors"
	// "fmt"
	"net/http"
	// "qexchange/models"
	// "qexchange/utils"

	// "golang.org/x/crypto/bcrypt"
	"qexchange/database"
	"qexchange/models"

	"gorm.io/gorm"
)

type CryptoService interface {
	GetCrypto(
		id int,
	) (models.CryptoResponse, int, error) // returns answer, statusCode, error

	SetCrypto(
		crypto models.Crypto,
	) (int, error) // returns statusCode, error

	UpdateCrypto(
		crypto models.Crypto,
	) (int, error) // returns statusCode, error

	GetAllCrypto() ([]models.CryptoResponse, int, error)
}

type cryptoService struct {
	db        *gorm.DB
	dbService database.DataBaseService
}

func NewCryptoService(db *gorm.DB) CryptoService {
	return &cryptoService{
		db:        db,
		dbService: database.NewDBService(db),
	}
}

func (s *cryptoService) GetCrypto(
	id int,
) (models.CryptoResponse, int, error) {
	var crypto models.Crypto
	result := s.db.Where("id = ?", id).First(&crypto)
	if result.Error != nil {
		return models.CryptoResponse{}, http.StatusBadRequest, errors.New("there is no crypto with this id")
	}

	return models.NewCryptoResponse(crypto), http.StatusOK, nil
}

func (s *cryptoService) SetCrypto(
	crypto models.Crypto,
) (int, error) {
	var cryptoSearch models.Crypto
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
	crypto models.Crypto,
) (int, error) {
	var cryptoSearch models.Crypto
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

	return http.StatusOK, nil
}

func (s *cryptoService) GetAllCrypto() ([]models.CryptoResponse, int, error) {
	var cryptoList []models.Crypto
	result := s.db.Find(&cryptoList)
	if result.Error != nil {
		return make([]models.CryptoResponse, 0), http.StatusBadRequest, errors.New("there is no crypto with this id")
	}

	cryptoResponseList := make([]models.CryptoResponse, 0)
	for _, cr := range cryptoList {
		cryptoResponseList = append(cryptoResponseList, models.NewCryptoResponse(cr))
	}
	return cryptoResponseList, http.StatusOK, nil
}
