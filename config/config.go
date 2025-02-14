package config

import (
	"encoding/json"
	"fmt"
	"log"
)

type ExchangeConfig struct {
	DB     DBConfig     `json:"db"`
	Server ServerConfig `json:"app" yaml:"app"`
	Redis  RedisConfig  `json:"redis"`
}

func (ExchangeConfig) configSignature() {}

func (cfg ExchangeConfig) Print() {
	jsonData, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal struct to JSON: %v", err)
	}

	fmt.Printf("loaded config: %v", string(jsonData))
}
