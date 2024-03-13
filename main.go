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
	userModels "qexchange/models/user"
	bankModels "qexchange/models/bank"

	"gorm.io/gorm"
)

//	@Title			ExchangeService
//	@version		1.1
//	@description	Exchange Service

//	@contact.name	Reza Mokaram
//	@contact.url	https://github.com/Quera-Go-Zilla

// @host			localhost:8080
// @BasePath		/
func main() {
	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("db connection failed: %v\n", err.Error())
	}

	err = migrate(db)
	if err != nil {
		log.Fatalf("migrations failed: %v\n", err.Error())
	}

	// start to dump test data into db if it hasn't been done already TODO
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
		&userModels.User{},
		&userModels.Profile{},
		&cryptocurrency.Crypto{},
		&bankModels.PaymentInfo{},
		&models.Transaction{},
		&bankModels.BankingInfo{},
		&models.SupportTicket{},
		&models.TicketMessage{},
		&trade.OpenTrade{},
		&trade.ClosedTrade{},
		&trade.FutureOrder{},
	)
}

func hasTestData(db *gorm.DB) bool {
	var count int64
	db.Model(&userModels.User{}).Count(&count)
	return count > 0
}
