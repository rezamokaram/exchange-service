package services

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"qexchange/models"
	"time"

	"gorm.io/gorm"
)

const (
	callbackURL     = "http://localhost:8080/payment/verify"
	zarinpalRequest = "https://sandbox.banktest.ir/zarinpal/api.zarinpal.com/pg/v4/payment/request.json"
	zarinpalVerify  = "https://sandbox.banktest.ir/zarinpal/api.zarinpal.com/pg/v4/payment/verify.json"
	zarinpalGateURL = "https://sandbox.banktest.ir/zarinpal/www.zarinpal.com/pg/StartPay/"
)

type BankService interface {
	AddBankAccount(user models.User, bank_name, account_number, card_number, expire_date, cvv2 string) (int, error)
	ChargeAccount(amount int, user models.User) (string, int, error) // returns payment_url, status code, error
	VerifyPayment(authority, status string) (int, error)             // returns status code, error
	AddToUserBalance(user models.User, amount, service int, description string) (int, error)
	SubtractFromUserBalance(user models.User, amount, service int, description string) (int, error)
	WithdrawFromAccount(user models.User, amount int, BankID uint) (int, error)
	GetAllTransactions(user models.User) ([]models.Transaction, int, error)
	GetAllPayments(user models.User) ([]models.PaymentInfo, int, error)
}

type bankService struct {
	db *gorm.DB
}

func NewBankService(db *gorm.DB) BankService {
	return &bankService{
		db: db,
	}
}

type ZarinpalData struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Authority string `json:"authority"`
	FeeType   string `json:"fee_type"`
	Fee       int    `json:"fee"`
}

type ZarinpalResponse struct {
	Data   ZarinpalData  `json:"data"`
	Errors []interface{} `json:"errors"`
}

func (s *bankService) AddBankAccount(user models.User, bank_name, account_number, card_number, expire_date, cvv2 string) (int, error) {

	if bank_name == "" || account_number == "" || card_number == "" || expire_date == "" || cvv2 == "" {
		return http.StatusNotFound, errors.New("account data is not provided")
	}

	// Create a new BankingInfo instance
	var bankInfo models.BankingInfo
	bankInfo.UserID = user.ID
	bankInfo.BankName = bank_name
	bankInfo.AccountNumber = account_number
	bankInfo.CardNumber = card_number
	bankInfo.Cvv2 = cvv2
	bankInfo.ExpireDate = expire_date

	if err := s.db.Save(&bankInfo).Error; err != nil {
		return http.StatusInternalServerError, errors.New("failed to save bank account")
	}

	return http.StatusOK, nil
}

func (s *bankService) ChargeAccount(amount int, user models.User) (string, int, error) {
	bankRequestData := map[string]interface{}{
		"merchant_id":  os.Getenv("MerchantID"),
		"amount":       amount,
		"callback_url": callbackURL,
		"description":  "Payment to charge Qexchgange account",
	}

	jsonData, err := json.Marshal(bankRequestData)
	if err != nil {
		return "", http.StatusInternalServerError, errors.New("failed parsing request data")
	}

	// this line disables ssl check
	// should be removed in production and use http instead of client
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	res, err := client.Post(zarinpalRequest, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	var result ZarinpalResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", http.StatusInternalServerError, err
	}

	// add payment to database
	newPayment := models.PaymentInfo{
		UserID:    user.ID,
		Amount:    int64(amount),
		Status:    "Wait",
		Authority: result.Data.Authority,
	}

	createdPayment := s.db.Create(&newPayment)
	if createdPayment.Error != nil {
		return "", http.StatusInternalServerError, errors.New("failed to insert payment into database")
	}

	PaymentUrl := fmt.Sprintf("%v%v", zarinpalGateURL, result.Data.Authority)

	return PaymentUrl, http.StatusOK, nil
}

func (s *bankService) VerifyPayment(authority, status string) (int, error) {
	// find the payment
	var payment models.PaymentInfo
	dbResult := s.db.Where("authority = ?", authority).First(&payment)
	if dbResult.Error != nil {
		return http.StatusBadRequest, errors.New("no payment found")
	}

	if payment.Status != "Wait" {
		return http.StatusBadRequest, errors.New("no payment waiting for verification found")
	}

	if status == "NOK" {
		payment.Status = "Failed"
		dbResult = s.db.Save(&payment)
		if dbResult.Error != nil {
			return http.StatusInternalServerError, errors.New("failed updating payment record")
		}
		return http.StatusBadRequest, errors.New("payment verification failed")
	}

	data := map[string]interface{}{
		"merchant_id": os.Getenv("MerchantID"),
		"amount":      payment.Amount,
		"authority":   authority,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		payment.Status = "Failed"
		dbResult = s.db.Save(&payment)
		if dbResult.Error != nil {
			return http.StatusInternalServerError, errors.New("failed updating payment record")
		}
		return http.StatusInternalServerError, errors.New("failed parsing bank data to verify payment")
	}

	// this line disables ssl check
	// should be removed in production and use http instead of client
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	res, err := client.Post(zarinpalVerify, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		payment.Status = "Failed"
		dbResult = s.db.Save(&payment)
		if dbResult.Error != nil {
			return http.StatusInternalServerError, errors.New("failed updating payment record")
		}
		return http.StatusInternalServerError, errors.New("failed to send request to verify payment")
	}
	defer res.Body.Close()

	jsonBody := make(map[string]interface{})
	err = json.NewDecoder(res.Body).Decode(&jsonBody)
	if err != nil {
		payment.Status = "Failed"
		dbResult = s.db.Save(&payment)
		if dbResult.Error != nil {
			return http.StatusInternalServerError, errors.New("failed updating payment record")
		}
		return http.StatusInternalServerError, errors.New("failed to parse verification response")
	}

	if data, ok := jsonBody["data"]; ok {
		if dataMap, ok := data.(map[string]interface{}); ok {
			if code, ok := dataMap["code"]; ok {
				if code == float64(100) {
					// verified => update database
					payment.Status = "Successful"
					dbResult = s.db.Save(&payment)
					if dbResult.Error != nil {
						return http.StatusInternalServerError, errors.New("failed updating payment record")
					}

					// update user balance
					var user models.User
					result := s.db.Where("id = ?", payment.UserID).First(&user)
					if result.Error != nil {
						return http.StatusInternalServerError, errors.New("failed finding user")
					}

					bankService := NewBankService(s.db)
					description := fmt.Sprintf("Bank Service: for payment with id = %v at %v", payment.ID, payment.CreatedAt)
					statusCode, err := bankService.AddToUserBalance(user, int(payment.Amount), 0, description)
					if err != nil {
						return statusCode, errors.New("failed updating user balance")
					}

					return http.StatusOK, nil

				} else if code == float64(101) {
					return http.StatusAlreadyReported, errors.New("payment already verified")
				} else {
					// failed => update database
					payment.Status = "Failed"
					dbResult = s.db.Save(&payment)
					if dbResult.Error != nil {
						return http.StatusInternalServerError, errors.New("failed updating payment record")
					}
					return http.StatusBadRequest, errors.New("code other than 100 or 101 returned")
				}
			} else {
				// failed => update database
				payment.Status = "Failed"
				dbResult = s.db.Save(&payment)
				if dbResult.Error != nil {
					return http.StatusInternalServerError, errors.New("failed updating payment record")
				}
				return http.StatusBadRequest, errors.New("no code in the json")
			}
		} else {
			// failed => update database
			payment.Status = "Failed"
			dbResult = s.db.Save(&payment)
			if dbResult.Error != nil {
				return http.StatusInternalServerError, errors.New("failed updating payment record")
			}
			return http.StatusBadRequest, errors.New("data in json failed")
		}
	} else {
		// failed => update database
		payment.Status = "Failed"
		dbResult = s.db.Save(&payment)
		if dbResult.Error != nil {
			return http.StatusInternalServerError, errors.New("failed updating payment record")
		}
		return http.StatusBadRequest, errors.New("no data in json")
	}
}

func (s *bankService) AddToUserBalance(user models.User, amount, service int, description string) (int, error) {
	var profile models.Profile
	result := s.db.Where("id = ?", user.ID).First(&profile)
	if result.Error != nil {
		return http.StatusBadRequest, errors.New("there is no profile with this id")
	}

	transaction := models.NewTransaction(user.ID, amount, service, true, description)
	result = s.db.Save(&transaction)
	if result.Error != nil {
		return http.StatusInternalServerError, result.Error
	}

	profile.Balance += amount

	result = s.db.Save(&profile)
	if result.Error != nil {
		return http.StatusInternalServerError, result.Error
	}
	return http.StatusAccepted, nil
}

func (s *bankService) SubtractFromUserBalance(user models.User, amount, service int, description string) (int, error) {
	var profile models.Profile
	result := s.db.Where("id = ?", user.ID).First(&profile)
	if result.Error != nil {
		return http.StatusBadRequest, errors.New("there is no profile with this id")
	}

	if profile.Balance < amount {
		return http.StatusBadRequest, errors.New("profile balance is lower than requested amount")
	}

	transaction := models.NewTransaction(user.ID, amount, service, false, description)
	result = s.db.Save(&transaction)
	if result.Error != nil {
		return http.StatusInternalServerError, result.Error
	}

	profile.Balance -= amount

	result = s.db.Save(&profile)
	if result.Error != nil {
		return http.StatusInternalServerError, result.Error
	}
	return http.StatusAccepted, nil
}

func (s *bankService) WithdrawFromAccount(user models.User, amount int, BankID uint) (int, error) {
	// get user with profile
	var userWithProfile models.User
	if err := s.db.Where("username = ?", user.Username).Preload("Profile").First(&userWithProfile).Error; err != nil {
		return http.StatusNotFound, errors.New("user not found")
	}

	if userWithProfile.Profile.Balance < amount {
		return http.StatusBadRequest, errors.New("not enough money in account")
	}

	var bankInfo models.BankingInfo
	if err := s.db.Where("id = ?", BankID).First(&bankInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return http.StatusNotFound, errors.New("banking info id is not valid")
		}
		return http.StatusNotFound, err
	}

	// newBalance := userWithProfile.Profile.Balance - amount
	// userWithProfile.Profile.Balance = newBalance

	// if err := s.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&userWithProfile).Error; err != nil {
	// 	return -1, http.StatusBadRequest, errors.New("failed updating user")
	// }

	description := fmt.Sprintf("Bank Service: withdraw from balance, amount = %v, at %v", amount, time.Now())
	statusCode, err := s.SubtractFromUserBalance(user, amount, 0, description)
	if err != nil {
		return statusCode, err
	}

	return http.StatusOK, nil
}

func (s *bankService) GetAllTransactions(user models.User) ([]models.Transaction, int, error) {
	var allTransactions []models.Transaction
	result := s.db.Where("user_id = ?", user.ID).Find(&allTransactions)
	if result.Error != nil {
		return make([]models.Transaction, 0), http.StatusInternalServerError, result.Error
	}
	return allTransactions, http.StatusOK, nil
}

func (s *bankService) GetAllPayments(user models.User) ([]models.PaymentInfo, int, error) {
	var allPayments []models.PaymentInfo
	result := s.db.Where("user_id = ?", user.ID).Find(&allPayments)
	if result.Error != nil {
		return make([]models.PaymentInfo, 0), http.StatusInternalServerError, result.Error
	}
	return allPayments, http.StatusOK, nil
}
