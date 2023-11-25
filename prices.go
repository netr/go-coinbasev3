package go_coinbasev3

// GetBuyPrice get the total price to buy a currency.
func (c *ApiClient) GetBuyPrice(pair string) (*CurrencyPairPrice, error) {
	return c.getPairPrice(pair, "buy")
}

// GetSellPrice get the total price to sell a currency.
func (c *ApiClient) GetSellPrice(pair string) (*CurrencyPairPrice, error) {
	return c.getPairPrice(pair, "sell")
}

// GetSpotPrice get the current market price of a currency.
func (c *ApiClient) GetSpotPrice(pair string) (*CurrencyPairPrice, error) {
	return c.getPairPrice(pair, "spot")
}

// getPairPrice get the price of a currency pair.
func (c *ApiClient) getPairPrice(pair string, side string) (*CurrencyPairPrice, error) {
	u := "https://api.coinbase.com/v2/prices/{currency_pair}/{side}"

	price := new(CurrencyPairPrice)
	resp, err := c.client.R().
		SetPathParam("currency_pair", pair).
		SetPathParam("side", side).
		SetSuccessResult(&price).Get(u)
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccessState() {
		return nil, ErrFailedToUnmarshal
	}

	return price, nil
}

// CurrencyPairPrice represents the price of a currency pair.
type CurrencyPairPrice struct {
	Data struct {
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	} `json:"data"`
}
