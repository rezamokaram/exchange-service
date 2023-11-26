package database

import (
	"qexchange/models"
	"qexchange/models/cryptocurrency"
	"qexchange/models/trade"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func CreateTestDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Set log level to Silent
	})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Profile{},
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
