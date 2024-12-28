package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"github.com/RezaMokaram/ExchangeService/handlers"
	"github.com/RezaMokaram/ExchangeService/models"
	cryptoModels "github.com/RezaMokaram/ExchangeService/models/crypto"
	"github.com/RezaMokaram/ExchangeService/server"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetCrypto(t *testing.T) {
	e := echo.New()
	server.PriceRoutes(e, testDB)

	t.Run("crypto not found", func(t *testing.T) {
		request := handlers.GetCryptoRequest{Id: 999} // Non-existent ID
		requestBody, _ := json.Marshal(request)
		req := httptest.NewRequest(http.MethodGet, "/crypto", bytes.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var errResp models.Response
		err := json.NewDecoder(rec.Body).Decode(&errResp)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		assert.Equal(t, http.StatusBadRequest, rec.Code, "Expected status code to be 400 Bad Request")
		assert.Contains(t, errResp.Message, "there is no crypto with this id", "Expected error message to contain 'there is no crypto with this id'")
	})

	t.Run("valid crypto request", func(t *testing.T) {
		testCrypto := cryptoModels.Crypto{
			Name:         "Bitcoin",
			Symbol:       "BTC",
			CurrentPrice: 500,
			BuyFee:       515,
			SellFee:      485,
		}

		request := handlers.GetCryptoRequest{Id: int(1)}
		requestBody, _ := json.Marshal(request)
		req := httptest.NewRequest(http.MethodGet, "/crypto", bytes.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var cryptoResponse cryptoModels.CryptoResponse
		err := json.NewDecoder(rec.Body).Decode(&cryptoResponse)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		assert.Equal(t, testCrypto.Name, cryptoResponse.Name, "Expected crypto name to match")
		assert.Equal(t, testCrypto.Symbol, cryptoResponse.Symbol, "Expected crypto symbol to match")
		assert.Equal(t, testCrypto.BuyFee, cryptoResponse.BuyFee, "Expected crypto BuyFee to match")
		assert.Equal(t, testCrypto.SellFee, cryptoResponse.SellFee, "Expected crypto SellFee to match")
	})
}

func TestSetCrypto(t *testing.T) {
	e := echo.New()
	server.PriceRoutes(e, testDB)
	server.UserRoutes(e, testDB)
	token := LoginAndGetToken(e, t, mockAdminUser)

	t.Run("crypto already exists", func(t *testing.T) {
		existingCryptoRequest := cryptoModels.MakeCryptoRequest{
			Name:         "Bitcoin",
			Symbol:       "BTC",
			CurrentPrice: 600,
		}

		requestBody, _ := json.Marshal(existingCryptoRequest)
		req := httptest.NewRequest(http.MethodPost, "/crypto", bytes.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Authorization", token)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var errResp models.Response
		err := json.NewDecoder(rec.Body).Decode(&errResp)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		assert.Equal(t, http.StatusBadRequest, rec.Code, "Expected status code to be 400 Bad Request")
		assert.Contains(t, errResp.Message, "the crypto already exist", "Expected error message to contain 'the crypto already exist'")
	})

	t.Run("invalid crypto data", func(t *testing.T) {
		invalidCryptoRequest := cryptoModels.MakeCryptoRequest{
			Name: "", // Missing name
		}

		requestBody, _ := json.Marshal(invalidCryptoRequest)
		req := httptest.NewRequest(http.MethodPost, "/crypto", bytes.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Authorization", token)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var errResp models.Response
		err := json.NewDecoder(rec.Body).Decode(&errResp)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		assert.Equal(t, http.StatusBadRequest, rec.Code, "Expected status code to be 400 Bad Request")
		assert.Contains(t, errResp.Message, "Bad Json Fields", "Expected error message to contain 'Bad Json Fields'")
	})

	t.Run("successfully add crypto", func(t *testing.T) {
		newCryptoRequest := cryptoModels.MakeCryptoRequest{
			Name:         "NewCrypto",
			Symbol:       "NCR",
			CurrentPrice: 1000,
		}

		requestBody, _ := json.Marshal(newCryptoRequest)
		req := httptest.NewRequest(http.MethodPost, "/crypto", bytes.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Authorization", token)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var successResp models.Response
		err := json.NewDecoder(rec.Body).Decode(&successResp)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		assert.Equal(t, http.StatusOK, rec.Code, "Expected status code to be 200 OK")
		assert.Contains(t, successResp.Message, "the crypto successfully added", "Expected success message to contain 'the crypto successfully added'")

		var addedCrypto cryptoModels.Crypto
		err = testDB.Where("name = ?", newCryptoRequest.Name).First(&addedCrypto).Error
		if err != nil {
			t.Fatalf("Failed to find the newly added crypto in the database: %v", err)
		}

		expectedSellFee := cryptoModels.CalculateSellFee(newCryptoRequest.CurrentPrice)
		expectedBuyFee := cryptoModels.CalculateBuyFee(newCryptoRequest.CurrentPrice)

		assert.Equal(t, newCryptoRequest.Name, addedCrypto.Name, "Expected crypto name in the database to match the request")
		assert.Equal(t, newCryptoRequest.Symbol, addedCrypto.Symbol, "Expected crypto symbol in the database to match the request")
		assert.Equal(t, newCryptoRequest.CurrentPrice, addedCrypto.CurrentPrice, "Expected crypto current price in the database to match the request")
		assert.Equal(t, addedCrypto.BuyFee, expectedBuyFee, "Expected crypto buy fee in the database to match the expected:", expectedBuyFee)
		assert.Equal(t, addedCrypto.SellFee, expectedSellFee, "Expected crypto sell fee in the database to match the expected:", expectedSellFee)
	})
}

func TestUpdateCrypto(t *testing.T) {
	e := echo.New()
	server.PriceRoutes(e, testDB)
	server.UserRoutes(e, testDB)
	token := LoginAndGetToken(e, t, mockAdminUser)

	t.Run("crypto not found for update", func(t *testing.T) {
		updateRequest := cryptoModels.UpdateCryptoRequest{
			Id:           999, // Non-existent ID
			Name:         "NonExistentCrypto",
			Symbol:       "NEX",
			CurrentPrice: 400,
		}

		requestBody, _ := json.Marshal(updateRequest)
		req := httptest.NewRequest(http.MethodPut, "/crypto", bytes.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Authorization", token)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var errResp models.Response
		err := json.NewDecoder(rec.Body).Decode(&errResp)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		assert.Equal(t, http.StatusBadRequest, rec.Code, "Expected status code to be 400 Bad Request")
		assert.Contains(t, errResp.Message, "there is no crypto with this id", "Expected error message to contain 'there is no crypto with this id'")
	})

	t.Run("successfully update crypto", func(t *testing.T) {
		updatedCryptoRequest := cryptoModels.UpdateCryptoRequest{
			Id:           1, // Assuming 'Bitcoin' has ID 1
			Name:         "BitcoinUpdated",
			Symbol:       "BTCU",
			CurrentPrice: 550, // Updated price
		}

		requestBody, _ := json.Marshal(updatedCryptoRequest)
		req := httptest.NewRequest(http.MethodPut, "/crypto", bytes.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Authorization", token)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var successResp models.Response
		err := json.NewDecoder(rec.Body).Decode(&successResp)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		assert.Equal(t, http.StatusOK, rec.Code, "Expected status code to be 200 OK")
		assert.Contains(t, successResp.Message, "the crypto successfully updated", "Expected success message to contain 'the crypto successfully updated'")

		// Verifying the updated entry in the database
		var updatedCrypto cryptoModels.Crypto
		err = testDB.Where("id = ?", updatedCryptoRequest.Id).First(&updatedCrypto).Error
		if err != nil {
			t.Fatalf("Failed to find the updated crypto in the database: %v", err)
		}

		assert.Equal(t, updatedCryptoRequest.Name, updatedCrypto.Name, "Expected crypto name in the database to match the updated name")
		assert.Equal(t, updatedCryptoRequest.Symbol, updatedCrypto.Symbol, "Expected crypto symbol in the database to match the updated symbol")
		assert.Equal(t, updatedCryptoRequest.CurrentPrice, updatedCrypto.CurrentPrice, "Expected crypto current price in the database to match the updated price")

		expectedSellFee := cryptoModels.CalculateSellFee(updatedCrypto.CurrentPrice)
		expectedBuyFee := cryptoModels.CalculateBuyFee(updatedCrypto.CurrentPrice)

		assert.Equal(t, updatedCrypto.BuyFee, expectedBuyFee, "Expected crypto buy fee in the database to match the expected:", expectedBuyFee)
		assert.Equal(t, updatedCrypto.SellFee, expectedSellFee, "Expected crypto sell fee in the database to match the expected:", expectedSellFee)
	})
}

func TestGetAllCrypto(t *testing.T) {
	e := echo.New()
	server.PriceRoutes(e, testDB)

	t.Run("successfully get all cryptos", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/crypto/get-all", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var cryptoList []cryptoModels.CryptoResponse
		err := json.NewDecoder(rec.Body).Decode(&cryptoList)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		assert.Equal(t, http.StatusOK, rec.Code, "Expected status code to be 200 OK")
		assert.NotEmpty(t, cryptoList, "Expected to receive a list of cryptocurrencies")
	})
}
