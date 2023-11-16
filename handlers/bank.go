package handlers

import (
	"net/http"
	"qexchange/models"
	"qexchange/services"

	"github.com/labstack/echo/v4"
)

type ChargeAccountRequest struct {
	Amount int `json:"amount"`
}

type ChargeAccountResponse struct {
	PaymentUrl string `json:"payment_url"`
}

func ChargeAccount(bankService services.BankService) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Parse the request body for amount
		user := c.Get("user").(models.User)
		request := new(ChargeAccountRequest)
		err := c.Bind(request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.NewErrorRespone("", err))
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
