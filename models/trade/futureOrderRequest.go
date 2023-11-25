package trade

// FutureOrderRequest represents the request to set a future order
type FutureOrderRequest struct {
	CryptoID    uint    `json:"crypto_id,omitempty"`
	Amount      float64 `json:"amount,omitempty"`
	TargetPrice int     `json:"target_price,omitempty"`
	StopLoss    int     `json:"stop_loss,omitempty"`
	TakeProfit  int     `json:"take_profit,omitempty"`
}

func (req FutureOrderRequest) ToFutureOrder(
	userID uint,
) FutureOrder {
	return FutureOrder{
		UserID:      userID,
		CryptoID:    req.CryptoID,
		Amount:      req.Amount,
		TargetPrice: req.TargetPrice,
		StopLoss:    req.StopLoss,
		TakeProfit:  req.TakeProfit,
	}
}
