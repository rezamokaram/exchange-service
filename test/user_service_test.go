package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"qexchange/database"
	"qexchange/handlers"
	"qexchange/models"
	"qexchange/server"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUserRegister(t *testing.T) {
	e := echo.New()
	db, _ := database.CreateTestDatabase()
	server.UserRoutes(e, db)

	t.Run("create new user", func(t *testing.T) {
		newUser := handlers.RegisterRequest{
			Username:       "testuser",
			Email:          "test@example.com",
			Password:       "password123",
			PasswordRepeat: "password123",
		}

		requestBody, err := json.Marshal(newUser)
		if err != nil {
			t.Fatalf("Failed to marshal request: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/user/register", bytes.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		// c := e.NewContext(req, rec)

		e.ServeHTTP(rec, req)

		if assert.Equal(t, http.StatusOK, rec.Code, "Expected status code to be 200 OK") {
			// Validate Database Entry
			var user models.User
			result := db.Where("email = ?", "test@example.com").First(&user)
			fmt.Println("user.Username ===>", user.Username)
			if assert.NoError(t, result.Error) {
				assert.Equal(t, "testuser", user.Username, "Expected username to match")
			}
		}
	})

	t.Run("passwords mismatch", func(t *testing.T) {
		mismatchUser := handlers.RegisterRequest{
			Username:       "mismatchuser",
			Email:          "mismatch@example.com",
			Password:       "password123",
			PasswordRepeat: "differentpassword",
		}

		requestBody, _ := json.Marshal(mismatchUser)
		req := httptest.NewRequest(http.MethodPost, "/user/register", bytes.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var errResp models.Response
		err := json.NewDecoder(rec.Body).Decode(&errResp)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// Assert the error message
		expectedErrMsg := ": passwords do not match"
		assert.Equal(t, http.StatusBadRequest, rec.Code, "Expected status code to be 400 Bad Request")
		assert.Equal(t, expectedErrMsg, errResp.Message, "Expected error message to match")
	})

	t.Run("duplicate username", func(t *testing.T) {
		duplicateUsername := handlers.RegisterRequest{
			Username:       "testuser", // assuming "testuser" was already created in previous test
			Email:          "unique@example.com",
			Password:       "password123",
			PasswordRepeat: "password123",
		}

		requestBody, _ := json.Marshal(duplicateUsername)
		req := httptest.NewRequest(http.MethodPost, "/user/register", bytes.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var errResp models.Response
		err := json.NewDecoder(rec.Body).Decode(&errResp)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		expectedErrMsg := ": a user with this username already exists"
		assert.Equal(t, http.StatusBadRequest, rec.Code, "Expected status code to be 400 Bad Request")
		assert.Equal(t, expectedErrMsg, errResp.Message, "Expected error message to match")
	})

	t.Run("duplicate emai", func(t *testing.T) {
		duplicateEmail := handlers.RegisterRequest{
			Username:       "uniqueuser",
			Email:          "test@example.com", // assuming "test@example.com" was used in the first test
			Password:       "password123",
			PasswordRepeat: "password123",
		}

		requestBody, _ := json.Marshal(duplicateEmail)
		req := httptest.NewRequest(http.MethodPost, "/user/register", bytes.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var errResp models.Response
		err := json.NewDecoder(rec.Body).Decode(&errResp)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		expectedErrMsg := ": a user with this email already exists"
		assert.Equal(t, http.StatusBadRequest, rec.Code, "Expected status code to be 400 Bad Request")
		assert.Equal(t, expectedErrMsg, errResp.Message, "Expected error message to match")
	})
}
