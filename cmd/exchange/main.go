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
	flag.StringVar(&path, "config_path", "./cmd/exchange/config.yaml", "path to clean env config file")
	flag.Parse()

	cfg := config.MustReadConfig[config.ExchangeConfig](path)
	appContainer := app.NewMustApp(cfg)

	log.Fatal(http.Run(appContainer, cfg.Server))
}
