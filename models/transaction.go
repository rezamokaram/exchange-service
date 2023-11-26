package models

import "gorm.io/gorm"

// Transaction represents a cryptocurrency transaction
type Transaction struct {
	gorm.Model
	UserID      uint    `gorm:"not null" json:"user_id" example:"2"`
	CryptoID    uint    `gorm:"not null" json:"crypto_id" example:"1"`
	Amount      float64 `gorm:"not null" json:"amount" example:"515"`
	PriceAtTime int     `gorm:"not null" json:"price_at_time" example:"500"`
}

func (Transaction) TableName() string {
	return "transactions"
}
