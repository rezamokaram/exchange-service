package handlers

import (
	"errors"
	"net/http"
	"qexchange/models"
	"qexchange/services"

	"github.com/labstack/echo/v4"
)

type ChargeAccountRequest struct {
	Amount int `json:"amount"`
}

type AddBankAccountRequest struct {
	BankName      string `json:"bank_name"`
	AccountNumber string `json:"account_number"`
	CardNumber    string `json:"card_number"`
	ExpireDate    string `json:"expire_date"`
	Cvv2          string `json:"cvv2"`
}

type ChargeAccountResponse struct {
	PaymentUrl string `json:"payment_url"`
}

func AddBankAccount(bankService services.BankService) echo.HandlerFunc {
	return func(c echo.Context) error {

		request := new(AddBankAccountRequest)
		err := c.Bind(request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.NewErrorRespone("", err))
		}

		user, bind := c.Get("user").(models.User)
		if !bind {
			response := models.NewErrorRespone("", errors.New("bad user data"))
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
			return c.JSON(statusCode, models.NewErrorRespone("", err))
		}

		// If successful, return a 200 OK response
		return c.JSON(http.StatusOK, models.NewRespone("bank info added successfully"))
	}
}

func ChargeAccount(bankService services.BankService) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Parse the request body for amount
		user := c.Get("user").(models.User)
		request := new(ChargeAccountRequest)
		err := c.Bind(request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.NewErrorRespone("account info not provided", err))
		}

		// Call the bank service to charge the bank account
		paymentURL, statusCode, err := bankService.ChargeAccount(request.Amount, user)
		if err != nil {
			return c.JSON(statusCode, UserResponse{Error: err.Error(), Message: "charging account failed"})
		}

		response := ChargeAccountResponse{
			PaymentUrl: paymentURL,
		}

		return c.JSON(http.StatusOK, response)
	}
}

func VerifyPayment(bankService services.BankService) echo.HandlerFunc {
	return func(c echo.Context) error {
		authority := c.QueryParam("Authority")
		status := c.QueryParam("Status")

		statusCode, err := bankService.VerifyPayment(authority, status)
		if err != nil {
			return c.JSON(statusCode, UserResponse{Error: err.Error(), Message: "payment verification failed"})
		}

		// return correct response here

		return c.JSON(http.StatusOK, UserResponse{Message: "success"})
	}
}
