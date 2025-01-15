package app

import (
	"context"

	"github.com/RezaMokaram/ExchangeService/config"
	cryptoPort "github.com/RezaMokaram/ExchangeService/internal/crypto/port"
	notifPort "github.com/RezaMokaram/ExchangeService/internal/notification/port"
	userPort "github.com/RezaMokaram/ExchangeService/internal/user/port"

	"gorm.io/gorm"
)

type App interface {
	// OrderService(ctx context.Context) orderPort.Service
	UserService(ctx context.Context) userPort.Service
	CryptoService(ctx context.Context) cryptoPort.Service
	NotificationService(ctx context.Context) notifPort.Service
	DB() *gorm.DB
	Config() config.AConfig
}
