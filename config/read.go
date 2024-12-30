package config

import (
	"encoding/json"
	"os"
)

func ReadConfig(configPath string) (AConfig, error) {
	var c AConfig
	all, err := os.ReadFile(configPath)
	if err != nil {
		return c, err
	}

	return c, json.Unmarshal(all, &c)
}

func MustReadConfig(configPath string) AConfig {
	c, err := ReadConfig(configPath)
	if err != nil {
		panic(err)
	}
	return c
}