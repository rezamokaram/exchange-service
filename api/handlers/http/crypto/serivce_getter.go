package crypto

import (
	"context"

	"github.com/rezamokaram/exchange-service/api/handlers/http/common"
	"github.com/rezamokaram/exchange-service/api/service"
	"github.com/rezamokaram/exchange-service/app"
	"github.com/rezamokaram/exchange-service/config"
)

// user service transient instance handler
func cryptoServiceGetter(appContainer app.App, cfg config.ServerConfig) common.ServiceGetter[*service.CryptoService] {
	return func(ctx context.Context) *service.CryptoService {
		return service.NewCryptoService(appContainer.CryptoService(ctx),
			cfg.Secret, cfg.AuthExpMinute, cfg.AuthRefreshMinute, appContainer.NotificationService(ctx))
	}
}
