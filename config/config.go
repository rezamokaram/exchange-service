package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App      APP      `env-required:"true" json:"app"`
		Postgres POSTGRES `env-required:"true" json:"postgres"`
	}

	APP struct {
		Name    string `env-required:"true" json:"name" env:"INSPECTION_APP_NAME"`
		Version string `env-required:"true" json:"version" env:"INSPECTION_APP_VERSION"`
		Host    string `env-required:"true" json:"host" env:"INSPECTION_APP_HOST"`
		Port    string `env-required:"true" json:"port" env:"INSPECTION_APP_PORT"`
	}

	POSTGRES struct {
		DB       string `env-required:"true" json:"db" env:"POSTGRES_DB"`
		User     string `env-required:"true" json:"user" env:"POSTGRES_USER"`
		Password string `env-required:"true" json:"password" env:"POSTGRES_PASSWORD"`
		Host     string `env-required:"true" json:"host" env:"POSTGRES_HOST"`
		Port     string `env-required:"true" json:"port" env:"POSTGRES_PORT"`
		SSLMode  string `env-required:"true" json:"ssl_mode" env:"POSTGRES_SSLMODE"`
		Timezone string `env-required:"true" json:"timezone" env:"POSTGRES_TIMEZONE"`
	}
)

func LoadConfig(path string) (*Config, error) {
	config := new(Config)
	err := cleanenv.ReadConfig(path, config)
	if err != nil {
		return nil, err
	}

	err = cleanenv.ReadEnv(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
