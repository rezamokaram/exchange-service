package trade

type ClosedTradeRequest struct {
	OpenTradeID			int     `json:"id"`
	Amount      		float64 `json:"amount"`
}