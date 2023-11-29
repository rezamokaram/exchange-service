package trade

import (
	"time"

	"gorm.io/gorm"
)

// closed trade should change a lot
// ClosedTrade represents a closed trade
type ClosedTrade struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID   uint `gorm:"not null" json:"user_id"`
	CryptoID uint `gorm:"not null" json:"crypto_id"`
	Amount   int  `gorm:"not null" json:"amount"`
	BuyFee   int  `gorm:"not null" json:"buy_fee"`
	SellFee  int  `gorm:"not null" json:"sell_fee"`
	Profit   int  `gorm:"not null" json:"profit"`
}

func (ClosedTrade) TableName() string {
	return "closed_trade"
}
