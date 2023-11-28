package trade

// OpenTradeRequest represents the request to open a trade
type OpenTradeRequest struct {
	CryptoID   uint    `json:"crypto_id"`
	Amount     float64 `json:"amount"`
	StopLoss   int     `json:"stop_loss"`
	TakeProfit int     `json:"take_profit"`
}

func (req OpenTradeRequest) ToOpenTrade(
	userId uint,
	buyFee int,
) OpenTrade {
	return OpenTrade{
		UserID:     userId,
		CryptoID:   req.CryptoID,
		Amount:     req.Amount,
		BuyFee:     buyFee,
		StopLoss:   req.StopLoss,
		TakeProfit: req.TakeProfit,
	}
}
