package handlers

import (
	"net/http"
	"qexchange/services"

	"github.com/labstack/echo/v4"
)

func BankCharge(bankService services.BankService) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Parse the request body for the bank account number and amount
		bankAccountNumber := c.FormValue("bankAccountNumber")
		amount := c.FormValue("amount")

		// Call the bank service to charge the bank account
		statusCode, err := bankService.ChargeBankAccount(bankAccountNumber, amount)
		if err != nil {
			return c.JSON(statusCode, map[string]string{"error": err.Error()})
		}

		// If successful, return a 200 OK response
		return c.JSON(http.StatusOK, map[string]string{"status": "success"})
	}
}