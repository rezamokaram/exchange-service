package trade

type ClosedTradeRequest struct {
	OpenTradeID			int     `json:"id"`
	CryptoID    		uint    `json:"crypto_id"`
	Amount      		float64 `json:"amount"`
}