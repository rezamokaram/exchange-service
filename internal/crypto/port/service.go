package port

import (
	"context"

	"github.com/RezaMokaram/ExchangeService/internal/crypto/domain"
)

type Service interface {
	CreateCrypto(ctx context.Context, crypto domain.Crypto) (domain.CryptoID, error)
	GetCryptoByFilter(ctx context.Context, filter *domain.CryptoFilter) (*domain.Crypto, error)
	// TODO
	// UpdateCrypto(ctx context.Context, crypto domain.TODO) (int, error)
	// GetListCrypto(ctx context.Context) ([]domain.TODO, int, error)
}
