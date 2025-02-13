package config

type ExchangeConfig struct {
	DB     DBConfig     `json:"db"`
	Server ServerConfig `json:"app"`
	Redis  RedisConfig  `json:"redis"`
}

func (ExchangeConfig) configSignature(){}