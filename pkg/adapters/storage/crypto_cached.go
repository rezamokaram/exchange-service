package storage

import (
	"context"

	"github.com/RezaMokaram/ExchangeService/internal/crypto/domain"
	"github.com/RezaMokaram/ExchangeService/internal/crypto/port"
	"github.com/RezaMokaram/ExchangeService/pkg/cache"
)

type cryptoCachedRepo struct {
	repo     port.Repo
	provider cache.Provider
}

func (r *cryptoCachedRepo) Create(ctx context.Context, user domain.Crypto) (domain.CryptoID, error) {
	panic("impl")
}

func (r *cryptoCachedRepo) GetByFilter(ctx context.Context, filter *domain.CryptoFilter) (*domain.Crypto, error) {
	panic("impl")
}
