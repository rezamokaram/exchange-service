package database

import (
	"qexchange/models"

	"gorm.io/gorm"
)

type DataBaseService interface {
	AddToUserBalanace(user models.User, amount int) error        // adds amount to user balance
	SubtractFromUserBalanace(user models.User, amount int) error // subtracts amount from user balance
}

type DBService struct {
	db *gorm.DB
}

func NewDBService(db *gorm.DB) DataBaseService {
	return &DBService{
		db: db,
	}
}

func (s *DBService) AddToUserBalanace(user models.User, amount int) error {
	newBalance := user.Profile.Balance + amount
	user.Profile.Balance = newBalance
	result := s.db.Save(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *DBService) SubtractFromUserBalanace(user models.User, amount int) error {
	newBalance := user.Profile.Balance - amount
	user.Profile.Balance = newBalance
	result := s.db.Save(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
