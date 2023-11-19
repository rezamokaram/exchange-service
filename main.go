package main

import (
	"fmt"
	"log"
	"os"
	"qexchange/database"
	"qexchange/models"
	"qexchange/models/trade"
	"qexchange/models/cryptocurrency"
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

	// start to dump test data into db
	// Read SQL file
	sqlFile, err := os.ReadFile("./main-data.sql")
	if err != nil {
		log.Fatalf("reading sql dump file failed: %v\n", err.Error())
	}

	sqlStatement := string(sqlFile)
	// Execute SQL
	result := db.Exec(sqlStatement)
	if result.Error != nil {
		log.Fatalf("executing sql dump file failed: %v\n", result.Error)
	}

	fmt.Println("Database operations done.")

	e := server.NewServer()

	server.RunServer(e, db)
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Profile{},
		&cryptocurrency.Crypto{},
		&models.PaymentInfo{},
		&models.Transaction{},
		&models.BankingInfo{},
		&models.SupportTicket{},
		&trade.OpenTrade{},
		&trade.ClosedTrade{},
		&trade.FutureOrder{},
	)
}
