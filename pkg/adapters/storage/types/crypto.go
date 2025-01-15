package types

import "gorm.io/gorm"

type Crypto struct {
	gorm.Model
	Name         string
	Symbol       string
	CurrentPrice uint64
	BuyFee       uint64
	SellFee      uint64
}
