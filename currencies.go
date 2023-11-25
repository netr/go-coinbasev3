package go_coinbasev3

// GetFiatCurrencies lists known fiat currencies. Currency codes conform to the ISO 4217 standard where possible
func (c *ApiClient) GetFiatCurrencies() (FiatCurrencies, error) {
	u := "https://api.coinbase.com/v2/currencies"

	var fiats FiatCurrencies
	resp, err := c.client.R().SetSuccessResult(&fiats).Get(u)
	if err != nil {
		return fiats, err
	}

	if !resp.IsSuccessState() {
		return fiats, ErrFailedToUnmarshal
	}

	return fiats, nil
}

// FiatCurrencies represents a list of fiat currencies.
type FiatCurrencies struct {
	Data []struct {
		Id      string `json:"id"`
		Name    string `json:"name"`
		MinSize string `json:"min_size"`
	} `json:"data"`
}

// GetCurrencies lists known cryptocurrencies.
func (c *ApiClient) GetCurrencies() (Currencies, error) {
	u := "https://api.coinbase.com/v2/currencies/crypto"

	var curr Currencies
	resp, err := c.client.R().SetSuccessResult(&curr).Get(u)
	if err != nil {
		return curr, err
	}

	if !resp.IsSuccessState() {
		return curr, ErrFailedToUnmarshal
	}

	return curr, nil
}

// Currencies represents a list of cryptocurrencies.
type Currencies struct {
	Data []struct {
		AssetId             string `json:"asset_id"`
		Code                string `json:"code"`
		Name                string `json:"name"`
		Color               string `json:"color"`
		SortIndex           int    `json:"sort_index"`
		Exponent            int    `json:"exponent"`
		Type                string `json:"type"`
		AddressRegex        string `json:"address_regex"`
		DestinationTagName  string `json:"destination_tag_name,omitempty"`
		DestinationTagRegex string `json:"destination_tag_regex,omitempty"`
	} `json:"data"`
}

// GetExchangeRates get current exchange rates. Default base currency is USD, but it can be defined as any supported currency
func (c *ApiClient) GetExchangeRates(currency string) (ExchangeRates, error) {
	u := "https://api.coinbase.com/v2/exchange-rates"

	if currency == "" {
		currency = "USD"
	}

	var rates ExchangeRates
	resp, err := c.client.R().
		SetQueryParam("currency", currency).
		SetSuccessResult(&rates).
		Get(u)
	if err != nil {
		return rates, err
	}

	if !resp.IsSuccessState() {
		return rates, ErrFailedToUnmarshal
	}

	return rates, nil
}

// ExchangeRates represents a list of exchange rates for a given currency.
type ExchangeRates struct {
	Data struct {
		Currency string            `json:"currency"`
		Rates    map[string]string `json:"rates"`
	} `json:"data"`
}
