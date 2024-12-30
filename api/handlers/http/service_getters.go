package http

import (
	"context"

	"github.com/RezaMokaram/ExchangeService/api/service"
	"github.com/RezaMokaram/ExchangeService/app"
	"github.com/RezaMokaram/ExchangeService/config"
)

// user service transient instance handler
func userServiceGetter(appContainer app.App, cfg config.ServerConfig) ServiceGetter[*service.UserService] {
	return func(ctx context.Context) *service.UserService {
		return service.NewUserService(appContainer.UserService(ctx),
			cfg.Secret, cfg.AuthExpMinute, cfg.AuthRefreshMinute, appContainer.NotificationService(ctx))
	}
}
