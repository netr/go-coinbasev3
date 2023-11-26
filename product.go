package coinbasev3

// GetProducts gets a list of available currency pairs for trading.
func (c *ApiClient) GetProducts() ([]Products, error) {
	u := "https://api.exchange.coinbase.com/products"

	var products []Products
	resp, err := c.client.R().
		SetSuccessResult(&products).
		Get(u)
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccessState() {
		return nil, ErrFailedToUnmarshal
	}

	return products, nil
}

type Products struct {
	Id                     string `json:"id"`
	BaseCurrency           string `json:"base_currency"`
	QuoteCurrency          string `json:"quote_currency"`
	QuoteIncrement         string `json:"quote_increment"`
	BaseIncrement          string `json:"base_increment"`
	DisplayName            string `json:"display_name"`
	MinMarketFunds         string `json:"min_market_funds"`
	MarginEnabled          bool   `json:"margin_enabled"`
	PostOnly               bool   `json:"post_only"`
	LimitOnly              bool   `json:"limit_only"`
	CancelOnly             bool   `json:"cancel_only"`
	Status                 string `json:"status"`
	StatusMessage          string `json:"status_message"`
	TradingDisabled        bool   `json:"trading_disabled"`
	FxStablecoin           bool   `json:"fx_stablecoin"`
	MaxSlippagePercentage  string `json:"max_slippage_percentage"`
	AuctionMode            bool   `json:"auction_mode"`
	HighBidLimitPercentage string `json:"high_bid_limit_percentage"`
}
