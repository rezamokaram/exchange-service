package trade

import "time"

type FilterTradesResponse struct {
	CryptoList    []uint        `json:"crypto_list,omitempty"`
	Start         time.Time     `json:"start,omitempty"`
	End           time.Time     `json:"end,omitempty"`
	ProfitOverAll int           `json:"profit_over_all,omitempty"`
	ClosedTrades  []ClosedTrade `json:"closed_trades,omitempty"`
}
