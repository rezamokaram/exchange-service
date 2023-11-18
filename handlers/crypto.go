package handlers

import (
	"net/http"
	"qexchange/models/cryptocurrency"
	"qexchange/models"
	"qexchange/services"

	"github.com/labstack/echo/v4"
)

type GetCryptoRequest struct {
	Id int `json:"id"`
}

type CryptoRequest struct {
	Name         string `json:"name"`
	Symbol       string `json:"symbol"`
	CurrentPrice int    `json:"current_price"`
}

func GetCrypto(service services.CryptoService) echo.HandlerFunc {
	return func(c echo.Context) error {

		request := new(GetCryptoRequest)
		if err := c.Bind(request); err != nil {
			response := models.NewErrorRespone("", err)
			return c.JSON(http.StatusBadRequest, response)
		}

		cryptoResponse, statusCode, err := service.GetCrypto(request.Id)
		if err != nil {
			response := models.NewErrorRespone("", err)
			return c.JSON(statusCode, response)
		}

		return c.JSON(http.StatusOK, cryptoResponse)
	}
}

func SetCrypto(service services.CryptoService) echo.HandlerFunc {
	return func(c echo.Context) error {

		request := new(cryptocurrency.Crypto)
		if err := c.Bind(request); err != nil {
			response := models.NewErrorRespone("", err)
			return c.JSON(http.StatusBadRequest, response)
		}

		statusCode, err := service.SetCrypto(*request)
		if err != nil {
			response := models.NewErrorRespone("", err)
			return c.JSON(statusCode, response)
		}

		return c.JSON(http.StatusOK, models.NewRespone("the crypto successfuly added"))
	}
}

func UpdateCrypto(service services.CryptoService) echo.HandlerFunc {
	return func(c echo.Context) error {

		request := new(cryptocurrency.Crypto)
		if err := c.Bind(request); err != nil {
			response := models.NewErrorRespone("", err)
			return c.JSON(http.StatusBadRequest, response)
		}

		statusCode, err := service.UpdateCrypto(*request)
		if err != nil {
			response := models.NewErrorRespone("", err)
			return c.JSON(statusCode, response)
		}

		return c.JSON(http.StatusOK, models.NewRespone("the crypto successfuly updated"))
	}
}

func GetAllCrypto(service services.CryptoService) echo.HandlerFunc {
	return func(c echo.Context) error {

		cryptoList, statusCode, err := service.GetAllCrypto()
		if err != nil {
			response := models.NewErrorRespone("", err)
			return c.JSON(statusCode, response)
		}

		return c.JSON(http.StatusOK, cryptoList)
	}
}
