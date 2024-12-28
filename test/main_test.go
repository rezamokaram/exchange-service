package test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	gormDB "github.com/RezaMokaram/ExchangeService/pkg/gorm_database"
	"github.com/RezaMokaram/ExchangeService/models"
	cryptoModels "github.com/RezaMokaram/ExchangeService/models/crypto"
	userModels "github.com/RezaMokaram/ExchangeService/models/user"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	mockValidUser = userModels.LoginRequest{
		Username: "user1",
		Password: "password",
	}

	mockAdminUser = userModels.LoginRequest{
		Username: "admin",
		Password: "password",
	}

	testDB *gorm.DB
)

func TestMain(m *testing.M) {
	var err error
	testDB, err = gormDB.CreateTestDatabase()
	if err != nil {
		log.Fatalf("Failed to create test database: %v", err)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	_ = gormDB.CloseTestDatabase(testDB)
	err = os.Remove("./test.db") // Replace with the actual path to test.db
	if err != nil {
		log.Fatalf("Failed to delete test database file: %v", err)
	}

	os.Exit(code)
}

func LoginAndGetToken(e *echo.Echo, t *testing.T, user userModels.LoginRequest) string {
	requestBody, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/user/login", bytes.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "Expected status code to be 200 OK")

	var tokenResponse userModels.LoginResponse
	err := json.NewDecoder(rec.Body).Decode(&tokenResponse)
	if err != nil {
		t.Fatalf("Failed to decode token response: %v", err)
	}

	return tokenResponse.Token
}

func ClearDatabaseTables(db *gorm.DB) error {
	tables := []interface{}{
		&models.SupportTicket{},
		&models.TicketMessage{},
		&cryptoModels.Crypto{},
	}

	for _, model := range tables {
		if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(model).Error; err != nil {
			return err
		}
	}

	return nil
}
