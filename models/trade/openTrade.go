package trade

import "gorm.io/gorm"

type OpenTrade struct {
	gorm.Model
	UserID      uint    `gorm:"not null"`
	CryptoID    uint    `gorm:"not null"`
	Amount      int 	`gorm:"not null"`
	BuyFee 		int     `gorm:"not null"`
	StopLoss	int		`gorm:""`
	TakeProfit  int		`gorm:""`
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
		Profit: profitCalculation(openTrade.Amount, openTrade.BuyFee, sellFee),
	}
}

func profitCalculation(amount int, buy , sell int) int {
	return (sell - buy) * amount
}