package cryptocurrency

type MakeCryptoRequest struct {
	Name         string `json:"name,omitempty" example:"SomeCoin"`
	Symbol       string `json:"symbol,omitempty" example:"SMC"`
	CurrentPrice int    `json:"current_price,omitempty" example:"500"`
}

func (mcr MakeCryptoRequest) IsValid() bool {
	if mcr.Name == "" || mcr.Symbol == "" || mcr.CurrentPrice < 100 {
		return false
	}
	return true
}

func (mcr MakeCryptoRequest) ToCrypto() Crypto {
	return Crypto{
		Name:         mcr.Name,
		Symbol:       mcr.Symbol,
		CurrentPrice: mcr.CurrentPrice,
		SellFee:      CalculateSellFee(mcr.CurrentPrice),
		BuyFee:       CalculateBuyFee(mcr.CurrentPrice),
	}
}

func CalculateSellFee(price int) int {
	ans := price - ((price / 100) + 10)
	if ans < 0 {
		ans = 0
	}
	return ans
}

func CalculateBuyFee(price int) int {
	return price + ((price / 100) + 10)
}
