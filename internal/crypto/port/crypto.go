package port

import (
	"context"

	"github.com/RezaMokaram/ExchangeService/internal/crypto/domain"
)

type Repo interface {
	Create(ctx context.Context, user domain.Crypto) (domain.CryptoID, error)
	// Update(ctx context.Context, user domain.Crypto) (domain.CryptoID, error)
	GetByFilter(ctx context.Context, filter *domain.CryptoFilter) (*domain.Crypto, error)
}
