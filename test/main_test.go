package test

import (
	"log"
	"os"
	"qexchange/database"
	"testing"

	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = database.CreateTestDatabase()
	if err != nil {
		log.Fatalf("Failed to create test database: %v", err)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	_ = database.CloseTestDatabase(testDB)
	err = os.Remove("./test.db") // Replace with the actual path to test.db
	if err != nil {
		log.Fatalf("Failed to delete test database file: %v", err)
	}

	os.Exit(code)
}
