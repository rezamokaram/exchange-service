package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/RezaMokaram/ExchangeService/api/server"
	"github.com/RezaMokaram/ExchangeService/config"
	gorm "github.com/RezaMokaram/ExchangeService/pkg/postgres"

	"github.com/RezaMokaram/ExchangeService/api/handlers/http"
	"github.com/RezaMokaram/ExchangeService/app"
)

func main() {
	var path string
	flag.StringVar(&path, "cleanenv", "./config/config.json", "path to clean env config file")
	flag.Parse()

	cfg, err := config.LoadConfig(path)
	if err != nil {
		log.Fatalf("load config failed: %v\n", err.Error())
	}

	db, err := gorm.NewGormDatabase(&cfg.Postgres)
	if err != nil {
		log.Fatalf("db connection failed: %v\n", err.Error())
	}

	fmt.Println("Database operations done.")

	go func() {
		c := config.MustReadConfig(path)
		// c := config.AConfig{}

		appContainer := app.NewMustApp(c)

		log.Fatal(http.Run(appContainer, c.Server))
	}()

	e := server.NewServer()
	server.RunServer(e, db, &cfg.App)
}
