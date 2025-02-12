package main

import (
	"flag"
	"log"

	"github.com/RezaMokaram/ExchangeService/config"

	"github.com/RezaMokaram/ExchangeService/api/handlers/http"
	"github.com/RezaMokaram/ExchangeService/app"
)

func main() {
	var path string
	flag.StringVar(&path, "cleanenv", "./config/config.json", "path to clean env config file")
	flag.Parse()

	c := config.MustReadConfig(path)
	// c := config.AConfig{}

	appContainer := app.NewMustApp(c)

	log.Fatal(http.Run(appContainer, c.Server))
}
