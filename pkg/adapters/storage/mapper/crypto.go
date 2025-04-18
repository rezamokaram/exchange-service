package mapper

import (
	"github.com/rezamokaram/exchange-service/internal/crypto/domain"
	"github.com/rezamokaram/exchange-service/pkg/adapters/storage/types"
	"gorm.io/gorm"
)

func CryptoDomain2Storage(cryptoDomain domain.Crypto) *types.Crypto {
	return &types.Crypto{
		Model: gorm.Model{
			ID:        uint(cryptoDomain.ID),
			CreatedAt: cryptoDomain.CreatedAt,
			DeletedAt: gorm.DeletedAt(ToNullTime(cryptoDomain.DeletedAt)),
		},
		Name:         cryptoDomain.Name,
		Symbol:       cryptoDomain.Symbol,
		CurrentPrice: cryptoDomain.CurrentPrice,
		BuyFee:       cryptoDomain.BuyFee,
		SellFee:      cryptoDomain.SellFee,
	}
}

func CryptoStorage2Domain(crypto types.Crypto) *domain.Crypto {
	return &domain.Crypto{
		ID:           domain.CryptoID(crypto.ID),
		CreatedAt:    crypto.CreatedAt,
		DeletedAt:    crypto.DeletedAt.Time,
		Name:         crypto.Name,
		Symbol:       crypto.Symbol,
		CurrentPrice: crypto.CurrentPrice,
		BuyFee:       crypto.BuyFee,
		SellFee:      crypto.SellFee,
	}
}
