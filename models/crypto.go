package models

import "gorm.io/gorm"

type Crypto struct {
	gorm.Model
	Name         string `gorm:"type:varchar(100);not null"`
	Symbol       string `gorm:"type:varchar(5);not null"`
	CurrentPrice int    `gorm:"not null"`
	BuyFee       int    `gorm:"not null"`
	SellFee      int    `gorm:"not null"`
}

func (Crypto) TableName() string {
	return "crypto"
}
