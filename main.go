package main

import (
	"fmt"
	"log"

	"qexchange/config"
	"qexchange/database"
	"qexchange/server"
)

//	@Title			ExchangeService
//	@version		2.1
//	@description	Exchange Service

//	@contact.name	Reza Mokaram
//	@contact.url	https://github.com/RezaMokaram

// @host			0.0.0.0:8080
// @BasePath		/
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("load config failed: %v\n", err.Error())
	}

	db, err := database.NewConnection(&cfg.Postgres)
	if err != nil {
		log.Fatalf("db connection failed: %v\n", err.Error())
	}

	fmt.Println("Database operations done.")

	e := server.NewServer()
	server.RunServer(e, db, &cfg.App)
}
