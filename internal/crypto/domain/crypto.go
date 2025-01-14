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
	CurrentPrice int
	BuyFee       int
	SellFee      int
}

type CryptoFilter struct {
	ID CryptoID
}
