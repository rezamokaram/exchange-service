package app

import (
	"context"

	"github.com/rezamokaram/exchange-service/config"
	cryptoPort "github.com/rezamokaram/exchange-service/internal/crypto/port"
	notifPort "github.com/rezamokaram/exchange-service/internal/notification/port"
	userPort "github.com/rezamokaram/exchange-service/internal/user/port"

	"gorm.io/gorm"
)

type App interface {
	// OrderService(ctx context.Context) orderPort.Service
	UserService(ctx context.Context) userPort.Service
	CryptoService(ctx context.Context) cryptoPort.Service
	NotificationService(ctx context.Context) notifPort.Service
	DB() *gorm.DB
	Config() config.ExchangeConfig
}
