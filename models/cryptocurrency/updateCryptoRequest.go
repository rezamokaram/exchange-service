package cryptocurrency

type UpdateCryptoRequest struct {
	Id           uint   `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Symbol       string `json:"symbol,omitempty"`
	CurrentPrice int    `json:"current_price,omitempty"`
}

func (ucr UpdateCryptoRequest)IsValid() bool {
	return ucr.Id != 0
}

func (req UpdateCryptoRequest)UpdateCrypto(c Crypto) Crypto {
	if req.CurrentPrice >= 100 {
		c.CurrentPrice = req.CurrentPrice
		c.BuyFee = CalculateBuyFee(c.CurrentPrice)
		c.SellFee = CalculateSellFee(c.CurrentPrice)
	}
	if req.Name != "" {
		c.Name = req.Name
	}
	if req.Symbol != "" {
		c.Symbol = req.Symbol
	}

	return c
}