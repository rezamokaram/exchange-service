package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config interface {
	configSignature()
}

func MustReadConfig[T Config](path string) T {
	var cfg T
	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		panic(err)
	}

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}