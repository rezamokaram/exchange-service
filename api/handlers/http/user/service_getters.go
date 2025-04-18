package user

import (
	"context"

	"github.com/rezamokaram/exchange-service/api/handlers/http/common"
	"github.com/rezamokaram/exchange-service/api/service"
	"github.com/rezamokaram/exchange-service/app"
	"github.com/rezamokaram/exchange-service/config"
)

// user service transient instance handler
func userServiceGetter(appContainer app.App, cfg config.ServerConfig) common.ServiceGetter[*service.UserService] {
	return func(ctx context.Context) *service.UserService {
		return service.NewUserService(appContainer.UserService(ctx),
			cfg.Secret, cfg.AuthExpMinute, cfg.AuthRefreshMinute, appContainer.NotificationService(ctx))
	}
}
