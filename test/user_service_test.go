package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"qexchange/handlers"
	"qexchange/models"
	"qexchange/server"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUserRegister(t *testing.T) {
	e := echo.New()
	server.UserRoutes(e, testDB)

	t.Run("create new user", func(t *testing.T) {
		newUser := handlers.RegisterRequest{
			Username:       "testuser",
			Email:          "test@example.com",
			Password:       "Password123%",
			PasswordRepeat: "Password123%",
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
			result := testDB.Where("email = ?", "test@example.com").First(&user)
			if assert.NoError(t, result.Error) {
				assert.Equal(t, "testuser", user.Username, "Expected username to match")
			}
		}
	})

	t.Run("passwords mismatch", func(t *testing.T) {
		mismatchUser := handlers.RegisterRequest{
			Username:       "mismatchuser",
			Email:          "mismatch@example.com",
			Password:       "Password123%",
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
		expectedErrMsg := "passwords do not match"
		assert.Equal(t, http.StatusBadRequest, rec.Code, "Expected status code to be 400 Bad Request")
		assert.Contains(t, errResp.Message, expectedErrMsg, "Expected error message to contain 'passwords do not match'")
	})

	t.Run("duplicate username", func(t *testing.T) {
		duplicateUsername := handlers.RegisterRequest{
			Username:       "testuser", // assuming "testuser" was already created in previous test
			Email:          "unique@example.com",
			Password:       "Password123%",
			PasswordRepeat: "Password123%",
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

		expectedErrMsg := "a user with this username already exists"
		assert.Equal(t, http.StatusBadRequest, rec.Code, "Expected status code to be 400 Bad Request")
		assert.Contains(t, errResp.Message, expectedErrMsg, "Expected error message to contain 'a user with this username already exists'")
	})

	t.Run("duplicate email", func(t *testing.T) {
		duplicateEmail := handlers.RegisterRequest{
			Username:       "uniqueuser",
			Email:          "test@example.com", // assuming "test@example.com" was used in the first test
			Password:       "Password123%",
			PasswordRepeat: "Password123%",
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

		expectedErrMsg := "a user with this email already exists"
		assert.Equal(t, http.StatusBadRequest, rec.Code, "Expected status code to be 400 Bad Request")
		assert.Contains(t, errResp.Message, expectedErrMsg, "Expected error message to contain 'a user with this email already exists'")
	})

	t.Run("email not provided", func(t *testing.T) {
		user := handlers.RegisterRequest{
			Username:       "testuser5",
			Email:          "", // Email not provided
			Password:       "Password123%",
			PasswordRepeat: "Password123%",
		}

		requestBody, _ := json.Marshal(user)
		req := httptest.NewRequest(http.MethodPost, "/user/register", bytes.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var errResp models.Response
		err := json.NewDecoder(rec.Body).Decode(&errResp)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		expectedErrMsg := "email not provided"
		assert.Equal(t, http.StatusBadRequest, rec.Code, "Expected status code to be 400 Bad Request")
		assert.Contains(t, errResp.Message, expectedErrMsg, "Expected error message to contain 'email not provided'")
	})

	t.Run("invalid email format", func(t *testing.T) {
		user := handlers.RegisterRequest{
			Username:       "testuser6",
			Email:          "test@example", // Invalid email format
			Password:       "password123",
			PasswordRepeat: "password123",
		}

		requestBody, _ := json.Marshal(user)
		req := httptest.NewRequest(http.MethodPost, "/user/register", bytes.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var errResp models.Response
		err := json.NewDecoder(rec.Body).Decode(&errResp)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		expectedErrMsg := "email is not valid"
		assert.Equal(t, http.StatusBadRequest, rec.Code, "Expected status code to be 400 Bad Request")
		assert.Contains(t, errResp.Message, expectedErrMsg, "Expected error message to contain 'email is not valid'")
	})

	t.Run("username not provided", func(t *testing.T) {
		user := handlers.RegisterRequest{
			Username:       "", // username not provided
			Email:          "test@example.com",
			Password:       "Password123%",
			PasswordRepeat: "Password123%",
		}

		requestBody, _ := json.Marshal(user)
		req := httptest.NewRequest(http.MethodPost, "/user/register", bytes.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var errResp models.Response
		err := json.NewDecoder(rec.Body).Decode(&errResp)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		expectedErrMsg := "username not provided"
		assert.Equal(t, http.StatusBadRequest, rec.Code, "Expected status code to be 400 Bad Request")
		assert.Contains(t, errResp.Message, expectedErrMsg, "Expected error message to contain 'username not provided'")
	})

	t.Run("password not provided", func(t *testing.T) {
		user := handlers.RegisterRequest{
			Username:       "testuser7",
			Email:          "test@example.com",
			Password:       "", // password not provided
			PasswordRepeat: "",
		}

		requestBody, _ := json.Marshal(user)
		req := httptest.NewRequest(http.MethodPost, "/user/register", bytes.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var errResp models.Response
		err := json.NewDecoder(rec.Body).Decode(&errResp)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		expectedErrMsg := "password not provided"
		assert.Equal(t, http.StatusBadRequest, rec.Code, "Expected status code to be 400 Bad Request")
		assert.Contains(t, errResp.Message, expectedErrMsg, "Expected error message to contain 'password not provided'")
	})

	t.Run("password not secure", func(t *testing.T) {
		insecurePasswordUser := handlers.RegisterRequest{
			Username:       "securetestuser",
			Email:          "securetest@example.com",
			Password:       "123",
			PasswordRepeat: "123",
		}

		requestBody, _ := json.Marshal(insecurePasswordUser)
		req := httptest.NewRequest(http.MethodPost, "/user/register", bytes.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var errResp models.Response
		err := json.NewDecoder(rec.Body).Decode(&errResp)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		expectedErrMsg := "password is not secure"
		assert.Equal(t, http.StatusBadRequest, rec.Code, "Expected status code to be 400 Bad Request")
		assert.Contains(t, errResp.Message, expectedErrMsg, "Expected error message to contain 'password is not secure'")
	})
}
