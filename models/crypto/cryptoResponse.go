package cryptocurrency

// CryptoResponse represents the response structure for cryptocurrency data
type CryptoResponse struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
	BuyFee  int    `json:"buy_fee"`
	SellFee int    `json:"sell_fee"`
}

func NewCryptoResponse(crypto Crypto) CryptoResponse {
	return CryptoResponse{
		Id:      int(crypto.ID),
		Name:    crypto.Name,
		Symbol:  crypto.Symbol,
		BuyFee:  crypto.BuyFee,
		SellFee: crypto.SellFee,
	}
}
