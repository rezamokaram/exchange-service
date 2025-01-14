package handlers

// import (
// 	"net/http"
// 	"github.com/RezaMokaram/ExchangeService/models"
// 	cryptoModels "github.com/RezaMokaram/ExchangeService/models/crypto"
// 	"github.com/RezaMokaram/ExchangeService/internal"

// 	"github.com/labstack/echo/v4"
// )

// // GetCryptoRequest represents the request body for getting a cryptoModels
// type GetCryptoRequest struct {
// 	Id int `json:"id"`
// }

// // CryptoRequest represents the request body for setting or updating a cryptoModels
// type CryptoRequest struct {
// 	Name         string `json:"name" example:"Bitcoin"`
// 	Symbol       string `json:"symbol" example:"BTC"`
// 	CurrentPrice int    `json:"current_price" example:"500"`
// }

// // GetCrypto handles the retrieval of a specific cryptoModels
// // @Summary Get cryptoModels
// // @Description Gets details of a specific cryptoModels by ID
// // @Tags Crypto
// // @Accept  json
// // @Produce  json
// // @Param   body  body     GetCryptoRequest  true  "Crypto ID"
// // @Success 200  {object}  cryptoModels.CryptoResponse
// // @Failure 400  {object}  models.Response
// // @Router /crypto [get]
// func GetCrypto(service internal.CryptoService) echo.HandlerFunc {
// 	return func(c echo.Context) error {

// 		request := new(GetCryptoRequest)
// 		if err := c.Bind(request); err != nil {
// 			response := models.NewErrorResponse("trade failed", err.Error())
// 			return c.JSON(http.StatusBadRequest, response)
// 		}

// 		cryptoResponse, statusCode, err := service.GetCrypto(request.Id)
// 		if err != nil {
// 			response := models.NewErrorResponse("trade failed", err.Error())
// 			return c.JSON(statusCode, response)
// 		}

// 		return c.JSON(http.StatusOK, cryptoResponse)
// 	}
// }

// // SetCrypto handles adding a new cryptoModels
// // @Summary Add new cryptoModels
// // @Description Adds a new cryptoModels to the system
// // @Tags Crypto
// // @Accept  json
// // @Param Authorization header string true "Authorization token"
// // @Security BasicAuth
// // @Produce  json
// // @Param   body  body      cryptoModels.MakeCryptoRequest  true  "Crypto Information"
// // @Success 200   {object}  models.Response
// // @Failure 400   {object}  models.Response
// // @Router /crypto [post]
// func SetCrypto(service internal.CryptoService) echo.HandlerFunc {
// 	return func(c echo.Context) error {

// 		request := new(cryptoModels.MakeCryptoRequest)
// 		if err := c.Bind(request); err != nil {
// 			response := models.NewErrorResponse("trade failed", err.Error())
// 			return c.JSON(http.StatusBadRequest, response)
// 		}

// 		if !request.IsValid() {
// 			response := models.NewErrorResponse("trade failed", "Bad Json Fields")
// 			return c.JSON(http.StatusBadRequest, response)
// 		}

// 		statusCode, err := service.SetCrypto(*request)
// 		if err != nil {
// 			response := models.NewErrorResponse("trade failed", err.Error())
// 			return c.JSON(statusCode, response)
// 		}

// 		return c.JSON(http.StatusOK, models.NewResponse("the crypto successfully added"))
// 	}
// }

// // UpdateCrypto handles updating an existing cryptoModels
// // @Summary Update cryptoModels
// // @Description Updates details of an existing cryptoModels
// // @Tags Crypto
// // @Accept  json
// // @Param Authorization header string true "Authorization token"
// // @Security BasicAuth
// // @Produce  json
// // @Param   body  body      cryptoModels.UpdateCryptoRequest  true  "Crypto Update Information"
// // @Success 200   {object}  models.Response
// // @Failure 400   {object}  models.Response
// // @Router /crypto [put]
// func UpdateCrypto(service internal.CryptoService) echo.HandlerFunc {
// 	return func(c echo.Context) error {

// 		request := new(cryptoModels.UpdateCryptoRequest)
// 		if err := c.Bind(request); err != nil {
// 			response := models.NewErrorResponse("trade failed", err.Error())
// 			return c.JSON(http.StatusBadRequest, response)
// 		}

// 		if !request.IsValid() {
// 			response := models.NewErrorResponse("trade failed", "Bad Json Fields")
// 			return c.JSON(http.StatusBadRequest, response)
// 		}

// 		statusCode, err := service.UpdateCrypto(*request)
// 		if err != nil {
// 			response := models.NewErrorResponse("trade failed", err.Error())
// 			return c.JSON(statusCode, response)
// 		}

// 		return c.JSON(http.StatusOK, models.NewResponse("the crypto successfully updated"))
// 	}
// }

// // GetAllCrypto handles retrieving all cryptocurrencies
// // @Summary Get all cryptocurrencies
// // @Description Retrieves details of all cryptocurrencies
// // @Tags Crypto
// // @Accept  json
// // @Produce  json
// // @Success 200  {array}  cryptoModels.CryptoResponse
// // @Failure 400  {object}  models.Response
// // @Router /crypto/get-all [get]
// func GetAllCrypto(service internal.CryptoService) echo.HandlerFunc {
// 	return func(c echo.Context) error {

// 		cryptoList, statusCode, err := service.GetAllCrypto()
// 		if err != nil {
// 			response := models.NewErrorResponse("trade failed", err.Error())
// 			return c.JSON(statusCode, response)
// 		}

// 		return c.JSON(http.StatusOK, cryptoList)
// 	}
// }
