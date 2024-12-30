package port

import (
	"context"

	"github.com/RezaMokaram/ExchangeService/internal/common"
	"github.com/RezaMokaram/ExchangeService/internal/notification/domain"
	userDomain "github.com/RezaMokaram/ExchangeService/internal/user/domain"
)

type Service interface {
	Send(ctx context.Context, notif *domain.Notification) error
	CheckUserNotifValue(ctx context.Context, userID userDomain.UserID, val string) (bool, error)
	common.OutboxHandler[domain.NotificationOutbox]
}
