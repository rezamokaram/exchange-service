package app

import (
	"context"

	"github.com/RezaMokaram/ExchangeService/config"
	notifPort "github.com/RezaMokaram/ExchangeService/internal/notification/port"
	// orderPort "github.com/RezaMokaram/ExchangeService/internal/order/port"
	userPort "github.com/RezaMokaram/ExchangeService/internal/user/port"

	"gorm.io/gorm"
)

type App interface {
	// OrderService(ctx context.Context) orderPort.Service
	UserService(ctx context.Context) userPort.Service
	NotificationService(ctx context.Context) notifPort.Service
	DB() *gorm.DB
	Config() config.AConfig
}
