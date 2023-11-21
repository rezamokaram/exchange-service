package main

import (
	"fmt"
	"log"
	"os"
	"qexchange/database"
	"qexchange/models"
	"qexchange/models/cryptocurrency"
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

	// start to dump test data into db if it hasn't been done already
	if !hasTestData(db) {
		// Read SQL file
		sqlFile, err := os.ReadFile("./main-data.sql")
		if err != nil {
			log.Fatalf("reading sql dump file failed: %v\n", err.Error())
		}

		// Execute SQL
		result := db.Exec(string(sqlFile))
		if result.Error != nil {
			log.Fatalf("executing sql dump file failed: %v\n", result.Error)
		}

		fmt.Println("Fake Data Inserted.")
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

func hasTestData(db *gorm.DB) bool {
	var count int64
	db.Model(&models.User{}).Count(&count)
	return count > 0
}
