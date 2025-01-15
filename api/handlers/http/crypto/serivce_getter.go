package crypto

import (
	"context"

	"github.com/RezaMokaram/ExchangeService/api/handlers/http/common"
	"github.com/RezaMokaram/ExchangeService/api/service"
	"github.com/RezaMokaram/ExchangeService/app"
	"github.com/RezaMokaram/ExchangeService/config"
)

// user service transient instance handler
func cryptoServiceGetter(appContainer app.App, cfg config.ServerConfig) common.ServiceGetter[*service.CryptoService] {
	return func(ctx context.Context) *service.CryptoService {
		return service.NewCryptoService(appContainer.CryptoService(ctx),
			cfg.Secret, cfg.AuthExpMinute, cfg.AuthRefreshMinute, appContainer.NotificationService(ctx))
	}
}
