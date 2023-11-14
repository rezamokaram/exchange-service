package services

import (
	"errors"
	"fmt"
	"net/http"
	"qexchange/models"
	"qexchange/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type BankService interface {
	ChargeBankAccount(bankAccountNumber, amount string) (int, error) //return status code error
}

type bankService struct {
	db *gorm.DB
}

func NewBankService(db *gorm.DB) BankService {
	return &bankService{
		db: db,
	}
}

func (s *bankService) ChargeBankAccount(bankAccountNumber, amount string) (int, error) {
	// Find the banking info
	var bankingInfo models.BankingInfo
	result := s.db.Where("account_number = ?", bankAccountNumber).First(&bankingInfo)
	if result.Error != nil {
		return http.StatusNotFound, errors.New("bank account not found")
	}

	// Charge the bank account
	// Here you might want to make an HTTP request to the bank's API
	// For now, let's just pretend that the bank account was successfully charged

	return http.StatusOK, nil
}

func (s *userService) ChargeProfile(userID uint, amount int) (int, error) {
    // Find the user's profile
    var profile models.Profile
    result := s.db.Where("user_id = ?", userID).First(&profile)
    if result.Error != nil {
        return http.StatusNotFound, errors.New("profile not found")
    }

    // Find the user's banking info
    var bankingInfo models.BankingInfo
    result = s.db.Where("user_id = ?", userID).First(&bankingInfo)
    if result.Error != nil {
        return http.StatusNotFound, errors.New("banking info not found")
    }

    // Charge the bank account
    statusCode, err := s.ChargeBankAccount(bankingInfo.AccountNumber, strconv.Itoa(amount))
    if err != nil {
        return statusCode, err
    }

    // Add the amount to the profile's balance
    profile.Balance += amount
    if err := s.db.Save(&profile).Error; err != nil {
        return http.StatusInternalServerError, err
    }

    return http.StatusOK, nil
}