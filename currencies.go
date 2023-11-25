package go_coinbasev3

// GetFiatCurrencies lists known fiat currencies. Currency codes conform to the ISO 4217 standard where possible
func (c *ApiClient) GetFiatCurrencies() (*FiatCurrencies, error) {
	u := "https://api.coinbase.com/v2/currencies"

	fiats := new(FiatCurrencies)
	resp, err := c.client.R().SetSuccessResult(&fiats).Get(u)
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccessState() {
		return nil, ErrFailedToUnmarshal
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
func (c *ApiClient) GetCurrencies() (*Currencies, error) {
	u := "https://api.coinbase.com/v2/currencies/crypto"

	curr := new(Currencies)
	resp, err := c.client.R().SetSuccessResult(&curr).Get(u)
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccessState() {
		return nil, ErrFailedToUnmarshal
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
