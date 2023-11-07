package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	UserID      uint    `gorm:"not null"`
	CryptoID    uint    `gorm:"not null"`
	PhoneNumber string  `gorm:"type:varchar(100)"`
	Type        int     `gorm:"not null"`
	Amount      float64 `gorm:"not null"`
	PriceAtTime int     `gorm:"not null"`
}

func (Transaction) TableName() string {
	return "transactions"
}
