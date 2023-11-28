package handlers

import (
	"net/http"
	"qexchange/models"
	"qexchange/models/cryptocurrency"
	"qexchange/services"

	"github.com/labstack/echo/v4"
)

// GetCryptoRequest represents the request body for getting a cryptocurrency
type GetCryptoRequest struct {
	Id int `json:"id"`
}

// CryptoRequest represents the request body for setting or updating a cryptocurrency
type CryptoRequest struct {
	Name         string `json:"name" example:"Bitcoin"`
	Symbol       string `json:"symbol" example:"BTC"`
	CurrentPrice int    `json:"current_price" example:"500"`
}

// GetCrypto handles the retrieval of a specific cryptocurrency
// @Summary Get cryptocurrency
// @Description Gets details of a specific cryptocurrency by ID
// @Tags Crypto
// @Accept  json
// @Produce  json
// @Param   body  body     GetCryptoRequest  true  "Crypto ID"
// @Success 200  {object}  cryptocurrency.CryptoResponse
// @Failure 400  {object}  models.Response
// @Router /crypto [get]
func GetCrypto(service services.CryptoService) echo.HandlerFunc {
	return func(c echo.Context) error {

		request := new(GetCryptoRequest)
		if err := c.Bind(request); err != nil {
			response := models.NewErrorResponse("", err.Error())
			return c.JSON(http.StatusBadRequest, response)
		}

		cryptoResponse, statusCode, err := service.GetCrypto(request.Id)
		if err != nil {
			response := models.NewErrorResponse("", err.Error())
			return c.JSON(statusCode, response)
		}

		return c.JSON(http.StatusOK, cryptoResponse)
	}
}

// SetCrypto handles adding a new cryptocurrency
// @Summary Add new cryptocurrency
// @Description Adds a new cryptocurrency to the system
// @Tags Crypto
// @Accept  json
// @Param Authorization header string true "Authorization token"
// @Security BasicAuth
// @Produce  json
// @Param   body  body      cryptocurrency.MakeCryptoRequest  true  "Crypto Information"
// @Success 200   {object}  models.Response
// @Failure 400   {object}  models.Response
// @Router /crypto [post]
func SetCrypto(service services.CryptoService) echo.HandlerFunc {
	return func(c echo.Context) error {

		request := new(cryptocurrency.MakeCryptoRequest)
		if err := c.Bind(request); err != nil {
			response := models.NewErrorResponse("", err.Error())
			return c.JSON(http.StatusBadRequest, response)
		}

		if !request.IsValid() {
			response := models.NewErrorResponse("", "Bad Json Fields")
			return c.JSON(http.StatusBadRequest, response)
		}

		statusCode, err := service.SetCrypto(*request)
		if err != nil {
			response := models.NewErrorResponse("", err.Error())
			return c.JSON(statusCode, response)
		}

		return c.JSON(http.StatusOK, models.NewResponse("the crypto successfully added"))
	}
}

// UpdateCrypto handles updating an existing cryptocurrency
// @Summary Update cryptocurrency
// @Description Updates details of an existing cryptocurrency
// @Tags Crypto
// @Accept  json
// @Param Authorization header string true "Authorization token"
// @Security BasicAuth
// @Produce  json
// @Param   body  body      cryptocurrency.UpdateCryptoRequest  true  "Crypto Update Information"
// @Success 200   {object}  models.Response
// @Failure 400   {object}  models.Response
// @Router /crypto [put]
func UpdateCrypto(service services.CryptoService) echo.HandlerFunc {
	return func(c echo.Context) error {

		request := new(cryptocurrency.UpdateCryptoRequest)
		if err := c.Bind(request); err != nil {
			response := models.NewErrorResponse("", err.Error())
			return c.JSON(http.StatusBadRequest, response)
		}

		if !request.IsValid() {
			response := models.NewErrorResponse("", "Bad Json Fields")
			return c.JSON(http.StatusBadRequest, response)
		}

		statusCode, err := service.UpdateCrypto(*request)
		if err != nil {
			response := models.NewErrorResponse("", err.Error())
			return c.JSON(statusCode, response)
		}

		return c.JSON(http.StatusOK, models.NewResponse("the crypto successfully updated"))
	}
}

// GetAllCrypto handles retrieving all cryptocurrencies
// @Summary Get all cryptocurrencies
// @Description Retrieves details of all cryptocurrencies
// @Tags Crypto
// @Accept  json
// @Produce  json
// @Success 200  {array}  cryptocurrency.CryptoResponse
// @Failure 400  {object}  models.Response
// @Router /crypto/get-all [get]
func GetAllCrypto(service services.CryptoService) echo.HandlerFunc {
	return func(c echo.Context) error {

		cryptoList, statusCode, err := service.GetAllCrypto()
		if err != nil {
			response := models.NewErrorResponse("", err.Error())
			return c.JSON(statusCode, response)
		}

		return c.JSON(http.StatusOK, cryptoList)
	}
}
