package database

import (
	"log"
	"os"

	"qexchange/models"
	"qexchange/models/cryptocurrency"
	"qexchange/models/trade"
	userModels "qexchange/models/user"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func CreateTestDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file:test.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Set log level to Silent
	})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&userModels.User{},
		&userModels.Profile{},
		&cryptocurrency.Crypto{},
		&models.PaymentInfo{},
		&models.Transaction{},
		&models.BankingInfo{},
		&models.SupportTicket{},
		&models.TicketMessage{},
		&trade.OpenTrade{},
		&trade.ClosedTrade{},
		&trade.FutureOrder{},
	)
	if err != nil {
		return nil, err
	}

	sqlFile, err := os.ReadFile("../database/main-data.sql")
	if err != nil {
		log.Fatalf("reading sql dump file failed: %v\n", err.Error())
	}

	// Execute SQL
	result := db.Exec(string(sqlFile))
	if result.Error != nil {
		log.Fatalf("executing sql dump file failed: %v\n", result.Error)
	}

	return db, nil
}

func CloseTestDatabase(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	err = sqlDB.Close()
	if err != nil {
		return err
	}

	return nil
}
