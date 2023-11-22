package trade

import "time"

type FilterTradesRequest struct {
	CryptoList 	[]uint    `json:"crypto_list,omitempty"`
	Start      	time.Time `json:"start,omitempty"`
	End        	time.Time `json:"end,omitempty"`
}
