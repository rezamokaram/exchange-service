package cryptocurrency

type UpdateCryptoRequest struct {
	Id           uint   `json:"id,omitempty" example:"1"`
	Name         string `json:"name,omitempty" example:""`
	Symbol       string `json:"symbol,omitempty" example:""`
	CurrentPrice int    `json:"current_price,omitempty" example:"510"`
}

func (ucr UpdateCryptoRequest) IsValid() bool {
	return ucr.Id != 0
}

func (req UpdateCryptoRequest) UpdateCrypto(c Crypto) Crypto {
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
