package storage

import (
	"context"

	"github.com/rezamokaram/exchange-service/internal/crypto/domain"
	"github.com/rezamokaram/exchange-service/internal/crypto/port"
	"github.com/rezamokaram/exchange-service/pkg/cache"
	"gorm.io/gorm"
)

type cryptoRepo struct {
	db *gorm.DB
}

func NewCryptoRepo(db *gorm.DB, cached bool, provider cache.Provider) port.Repo {
	repo := &cryptoRepo{db}
	if !cached {
		return repo
	}

	return &cryptoCachedRepo{
		repo:     repo,
		provider: provider,
	}
}

func (r *cryptoRepo) Create(ctx context.Context, user domain.Crypto) (domain.CryptoID, error) {
	panic("impl")
}

func (r *cryptoRepo) GetByFilter(ctx context.Context, filter *domain.CryptoFilter) (*domain.Crypto, error) {
	panic("impl")
}
