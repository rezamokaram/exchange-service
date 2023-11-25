package models

import "gorm.io/gorm"

// Transaction represents a cryptocurrency transaction
type Transaction struct {
	gorm.Model
	UserID      uint    `gorm:"not null" json:"user_id"`
	CryptoID    uint    `gorm:"not null" json:"crypto_id"`
	Amount      float64 `gorm:"not null" json:"amount"`
	PriceAtTime int     `gorm:"not null" json:"price_at_time"`
}

func (Transaction) TableName() string {
	return "transactions"
}
