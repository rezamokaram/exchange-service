package main

import (
	"log"
	"qexchange/database"
	"qexchange/models"
	"qexchange/models/trade"
	"qexchange/server"

	"gorm.io/gorm"
)

func main() {
	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("db connection failed: %v\n", err.Error())
	}

	err = migrate(db)
	if err != nil {
		log.Fatalf("migrations failed: %v\n", err.Error())
	}

	e := server.NewServer()

	server.RunServer(e, db)
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Profile{},
		&models.Crypto{},
		&models.PaymentInfo{},
		&models.Transaction{},
		&models.BankingInfo{},
		&models.SupportTicket{},
		&trade.OpenTrade{},
		&trade.ClosedTrade{},
	)
}
