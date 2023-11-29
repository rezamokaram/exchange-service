package trade

import "time"

// FilterTradesRequest represents the request to filter trades
type FilterTradesRequest struct {
	CryptoList []uint    `json:"crypto_list,omitempty"`
	Start      time.Time `json:"start,omitempty"`
	End        time.Time `json:"end,omitempty"`
}
