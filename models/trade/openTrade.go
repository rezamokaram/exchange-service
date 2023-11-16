package trade

import "gorm.io/gorm"

type OpenTrade struct {
	gorm.Model
	UserID      uint    `gorm:"not null" json:"user_id"`
	CryptoID    uint    `gorm:"not null" json:"crypto_id"`
	Amount      float64 `gorm:"not null" json:"amount"`
	BuyFee 		int     `gorm:"not null" json:"price_at_time"`
	StopLoss	int		`gorm:"" json:"stop_loss"`
	TakeProfit  int		`gorm:"" json:"take_profit"`
}

func (OpenTrade) TableName() string {
	return "open_trade"
}

func (openTrade OpenTrade)ToCloseTrade(
	sellFee int,
) ClosedTrade {
	return ClosedTrade{
		UserID: openTrade.UserID,
		CryptoID: openTrade.CryptoID,
		Amount: openTrade.Amount,
		BuyFee: openTrade.BuyFee,
		SellFee: sellFee,
	}
}