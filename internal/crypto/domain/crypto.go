package domain

import (
	"time"
)

type (
	CryptoID uint
)

type Crypto struct {
	ID           CryptoID
	CreatedAt    time.Time
	DeletedAt    time.Time
	Name         string
	Symbol       string
	CurrentPrice uint64
	BuyFee       uint64
	SellFee      uint64
}

type CryptoFilter struct {
	ID CryptoID
}
