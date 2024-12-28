package handlers

import (
	"net/http"
	"github.com/RezaMokaram/ExchangeService/models"
	userModels "github.com/RezaMokaram/ExchangeService/models/user"
	"github.com/RezaMokaram/ExchangeService/services"

	"github.com/labstack/echo/v4"
)

// ChargeAccountRequest represents the request body for charging an account
type ChargeAccountRequest struct {
	Amount int `json:"amount"`
}

// AddBankAccountRequest represents the request body for adding a bank account
type AddBankAccountRequest struct {
	BankName      string `json:"bank_name" example:"eghtesad novin"`
	AccountNumber string `json:"account_number" example:"123456"`
	CardNumber    string `json:"card_number" example:"654321"`
	ExpireDate    string `json:"expire_date" example:"04/10"`
	Cvv2          string `json:"cvv2" example:"123"`
}

// WithdrawFromAccountRequest represents the request body for withdrawing money from bank account
type WithdrawFromAccountRequest struct {
	Amount int  `json:"amount" example:"5000"`
	BankID uint `json:"bank_id" example:"1"`
}

// ChargeAccountResponse represents the response for charging an account
type ChargeAccountResponse struct {
	PaymentUrl string `json:"payment_url"`
}

// AddBankAccount handles adding a new bank account
// @Summary Add Bank Account
// @Description Adds a bank account for a user
// @Tags Bank
// @Accept  json
// @Param Authorization header string true "Authorization token"
// @Security BasicAuth
// @Produce  json
// @Produce json
// @Param   body  body      AddBankAccountRequest  true  "Add Bank Account"
// @Success 200   {object}  models.Response
// @Failure 400   {object}  models.Response
// @Router /bank/add_account [post]
func AddBankAccount(bankService services.BankService) echo.HandlerFunc {
	return func(c echo.Context) error {

		request := new(AddBankAccountRequest)
		err := c.Bind(request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.NewErrorResponse("bank action failed", err.Error()))
		}

		user, bind := c.Get("user").(userModels.User)
		if !bind {
			response := models.NewErrorResponse("bank action failed", "bad user data")
			return c.JSON(http.StatusBadRequest, response)
		}

		statusCode, err := bankService.AddBankAccount(
			user,
			request.BankName,
			request.AccountNumber,
			request.CardNumber,
			request.ExpireDate,
			request.Cvv2,
		)
		if err != nil {
			return c.JSON(statusCode, models.NewErrorResponse("bank action failed", err.Error()))
		}

		// If successful, return a 200 OK response
		return c.JSON(http.StatusOK, models.NewResponse("bank info added successfully"))
	}
}

// ChargeAccount handles charging a user's account
// @Summary Charge Account
// @Description Charges a user's account
// @Tags Bank
// @Accept  json
// @Param Authorization header string true "Authorization token"
// @Security BasicAuth
// @Produce json
// @Param   body  body      ChargeAccountRequest   true  "Charge Account"
// @Success 200   {object}  ChargeAccountResponse
// @Failure 400   {object}  models.Response
// @Router /bank/payment/charge [post]
func ChargeAccount(bankService services.BankService) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Parse the request body for amount
		user, bind := c.Get("user").(userModels.User)
		if !bind {
			response := models.NewErrorResponse("bank action failed", "bad user data")
			return c.JSON(http.StatusBadRequest, response)
		}

		request := new(ChargeAccountRequest)
		err := c.Bind(request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.NewErrorResponse("account info not provided", err.Error()))
		}

		// Call the bank service to charge the bank account
		paymentURL, statusCode, err := bankService.ChargeAccount(request.Amount, user)
		if err != nil {
			return c.JSON(statusCode, models.NewErrorResponse("charging account failed", err.Error()))
		}

		response := ChargeAccountResponse{
			PaymentUrl: paymentURL,
		}

		return c.JSON(http.StatusOK, response)
	}
}

// VerifyPayment handles the verification of a payment transaction
// @Summary Verify Payment
// @Description Verifies a payment transaction
// @Tags Bank
// @Accept  json
// @Produce json
// @Success 200   {object}  models.Response
// @Router /bank/payment/verify [get]
func VerifyPayment(bankService services.BankService) echo.HandlerFunc {
	return func(c echo.Context) error {
		authority := c.QueryParam("Authority")
		status := c.QueryParam("Status")

		statusCode, err := bankService.VerifyPayment(authority, status)
		if err != nil {
			return c.JSON(statusCode, models.NewErrorResponse("payment verification failed", err.Error()))
		}

		// return correct response here
		return c.JSON(http.StatusOK, models.NewResponse("success"))
	}
}

// WithdrawFromAccount handles withdrawing money from a user's account balance
// @Summary Withdraw from Account
// @Description Allows a user to withdraw a specified amount from their account balance
// @Tags Bank
// @Accept  json
// @Param Authorization header string true "Authorization token"
// @Security BasicAuth
// @Produce json
// @Param   body  body      WithdrawFromAccountRequest  true  "Withdraw from Account"
// @Success 200   {object}  models.Response
// @Failure 400   {object}  models.Response
// @Router /bank/payment/withdraw [post]
func WithdrawFromAccount(bankService services.BankService) echo.HandlerFunc {
	return func(c echo.Context) error {

		request := new(WithdrawFromAccountRequest)
		err := c.Bind(request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.NewErrorResponse("bank action failed", err.Error()))
		}

		user, bind := c.Get("user").(userModels.User)
		if !bind {
			response := models.NewErrorResponse("bank action failed", "bad user data")
			return c.JSON(http.StatusBadRequest, response)
		}

		if request.Amount == 0 || request.BankID == 0 {
			response := models.NewErrorResponse("bank action failed", "amount or bank_id not provided")
			return c.JSON(http.StatusBadRequest, response)
		}

		statusCode, err := bankService.WithdrawFromAccount(user, request.Amount, request.BankID)
		if err != nil {
			return c.JSON(statusCode, models.NewErrorResponse("bank action failed", err.Error()))
		}

		return c.JSON(
			http.StatusOK,
			models.NewResponse("balance updated successfully"),
		)
	}
}

// GetAllTransactions gets all transactions
// @Summary Get all transactions
// @Description Retrieves all transactions
// @Tags Bank
// @Produce  json
// @Param Authorization header string true "Authorization token"
// @Security BasicAuth
// @Success 200  {array}   models.Transaction
// @Failure 400  {object}  models.Response
// @Router /bank/transaction/get-all [get]
func GetAllTransactions(bankService services.BankService) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, bind := c.Get("user").(userModels.User)
		if !bind {
			response := models.NewErrorResponse("bank action failed", "bad user data")
			return c.JSON(http.StatusBadRequest, response)
		}

		allTransactions, statusCode, err := bankService.GetAllTransactions(user)
		if err != nil {
			return c.JSON(statusCode, models.NewErrorResponse("bank action failed", err.Error()))
		}

		return c.JSON(
			http.StatusOK,
			allTransactions,
		)
	}
}

// GetAllPayments gets all payments
// @Summary Get all payments
// @Description Retrieves all payments
// @Tags Bank
// @Produce  json
// @Param Authorization header string true "Authorization token"
// @Security BasicAuth
// @Success 200  {array}   models.PaymentInfo
// @Failure 400  {object}  models.Response
// @Router /bank/payment/get-all [get]
func GetAllPayments(bankService services.BankService) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, bind := c.Get("user").(userModels.User)
		if !bind {
			response := models.NewErrorResponse("bank action failed", "bad user data")
			return c.JSON(http.StatusBadRequest, response)
		}

		allPayments, statusCode, err := bankService.GetAllPayments(user)
		if err != nil {
			return c.JSON(statusCode, models.NewErrorResponse("bank action failed", err.Error()))
		}

		return c.JSON(
			http.StatusOK,
			allPayments,
		)
	}
}
