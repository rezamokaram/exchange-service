package handlers

import (
	"net/http"
	"qexchange/models"
	"qexchange/services"

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

// ChargeAccountResponse represents the response for charging an account
type ChargeAccountResponse struct {
	PaymentUrl string `json:"payment_url"`
}

// AddBankAccount handles adding a new bank account
// @Summary Add Bank Account
// @Description Adds a bank account for a user
// @Accept  json
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
			return c.JSON(http.StatusBadRequest, models.NewErrorResponse("", err.Error()))
		}

		user, bind := c.Get("user").(models.User)
		if !bind {
			response := models.NewErrorResponse("", "bad user data")
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
			return c.JSON(statusCode, models.NewErrorResponse("", err.Error()))
		}

		// If successful, return a 200 OK response
		return c.JSON(http.StatusOK, models.NewResponse("bank info added successfully"))
	}
}

// ChargeAccount handles charging a user's account
// @Summary Charge Account
// @Description Charges a user's account
// @Accept  json
// @Produce json
// @Param   body  body      ChargeAccountRequest   true  "Charge Account"
// @Success 200   {object}  ChargeAccountResponse
// @Failure 400   {object}  models.Response
// @Router /bank/payment/charge [post]
func ChargeAccount(bankService services.BankService) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Parse the request body for amount
		user, bind := c.Get("user").(models.User)
		if !bind {
			response := models.NewErrorResponse("", "bad user data")
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
