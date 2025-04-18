package main

import (
	"flag"
	"log"

	"github.com/rezamokaram/exchange-service/config"

	"github.com/rezamokaram/exchange-service/api/handlers/http"
	"github.com/rezamokaram/exchange-service/app"
)

func main() {
	var path string
	flag.StringVar(&path, "config_path", "./cmd/exchange/config.yaml", "path to clean env config file")
	flag.Parse()

	cfg := config.MustReadConfig[config.ExchangeConfig](path)
	appContainer := app.NewMustApp(cfg)

	log.Fatal(http.Run(appContainer, cfg.Server))
}
