package trade

// ClosedTradeRequest represents the request to close a trade
type ClosedTradeRequest struct {
	OpenTradeID int     `json:"id"`
	Amount      int 	`json:"amount"`
}
