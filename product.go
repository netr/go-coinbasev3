package coinbasev3

import (
	"fmt"
	"strings"
)

// GetProduct get information on a single product by product ID.
func (c *ApiClient) GetProduct(productId string) (Product, error) {
	u := c.makeV3Url(fmt.Sprintf("/brokerage/products/%s", productId))

	var data Product
	resp, err := c.httpClient.Get(u)
	if err != nil {
		return data, err
	}

	if !resp.IsSuccessState() {
		return data, ErrFailedToUnmarshal
	}

	err = resp.Unmarshal(&data)
	if err != nil {
		return data, err
	}

	return data, nil
}

type Product struct {
	ProductId                 string                   `json:"product_id"`
	Price                     string                   `json:"price"`
	PricePercentageChange24H  string                   `json:"price_percentage_change_24h"`
	Volume24H                 string                   `json:"volume_24h"`
	VolumePercentageChange24H string                   `json:"volume_percentage_change_24h"`
	BaseIncrement             string                   `json:"base_increment"`
	QuoteIncrement            string                   `json:"quote_increment"`
	QuoteMinSize              string                   `json:"quote_min_size"`
	QuoteMaxSize              string                   `json:"quote_max_size"`
	BaseMinSize               string                   `json:"base_min_size"`
	BaseMaxSize               string                   `json:"base_max_size"`
	BaseName                  string                   `json:"base_name"`
	QuoteName                 string                   `json:"quote_name"`
	Watched                   bool                     `json:"watched"`
	IsDisabled                bool                     `json:"is_disabled"`
	New                       bool                     `json:"new"`
	Status                    string                   `json:"status"`
	CancelOnly                bool                     `json:"cancel_only"`
	LimitOnly                 bool                     `json:"limit_only"`
	PostOnly                  bool                     `json:"post_only"`
	TradingDisabled           bool                     `json:"trading_disabled"`
	AuctionMode               bool                     `json:"auction_mode"`
	ProductType               string                   `json:"product_type"`
	QuoteCurrencyId           string                   `json:"quote_currency_id"`
	BaseCurrencyId            string                   `json:"base_currency_id"`
	FcmTradingSessionDetails  FcmTradingSessionDetails `json:"fcm_trading_session_details"`
	MidMarketPrice            string                   `json:"mid_market_price"`
	Alias                     string                   `json:"alias"`
	AliasTo                   []string                 `json:"alias_to"`
	BaseDisplaySymbol         string                   `json:"base_display_symbol"`
	QuoteDisplaySymbol        string                   `json:"quote_display_symbol"`
	ViewOnly                  bool                     `json:"view_only"`
	PriceIncrement            string                   `json:"price_increment"`
	FutureProductDetails      FutureProductDetails     `json:"future_product_details"`
}

type FcmTradingSessionDetails struct {
	IsSessionOpen string `json:"is_session_open"`
	OpenTime      string `json:"open_time"`
	CloseTime     string `json:"close_time"`
}

type FutureProductDetails struct {
	Venue                  string           `json:"venue"`
	ContractCode           string           `json:"contract_code"`
	ContractExpiry         string           `json:"contract_expiry"`
	ContractSize           string           `json:"contract_size"`
	ContractRootUnit       string           `json:"contract_root_unit"`
	GroupDescription       string           `json:"group_description"`
	ContractExpiryTimezone string           `json:"contract_expiry_timezone"`
	GroupShortDescription  string           `json:"group_short_description"`
	RiskManagedBy          string           `json:"risk_managed_by"`
	ContractExpiryType     string           `json:"contract_expiry_type"`
	PerpetualDetails       PerpetualDetails `json:"perpetual_details"`
	ContractDisplayName    string           `json:"contract_display_name"`
}

type PerpetualDetails struct {
	OpenInterest string `json:"open_interest"`
	FundingRate  string `json:"funding_rate"`
	FundingTime  string `json:"funding_time"`
}

// GetProducts gets a list of available currency pairs for trading.
func (c *ApiClient) GetProducts() ([]Products, error) {
	u := c.makeExchangeUrl("/products")

	var data []Products
	resp, err := c.httpClient.Get(u)
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccessState() {
		return nil, ErrFailedToUnmarshal
	}

	err = resp.Unmarshal(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
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

type Granularity string

const (
	GranularityUnknown    Granularity = "UNKNOWN_GRANULARITY"
	GranularityOneMin     Granularity = "ONE_MINUTE"
	GranularityFiveMin    Granularity = "FIVE_MINUTE"
	GranularityFifteenMin Granularity = "FIFTEEN_MINUTE"
	GranularityThirtyMin  Granularity = "THIRTY_MINUTE"
	GranularityOneHour    Granularity = "ONE_HOUR"
	GranularityTwoHour    Granularity = "TWO_HOUR"
	GranularitySixHour    Granularity = "SIX_HOUR"
	GranularityOneDay     Granularity = "ONE_DAY"
)

// GetProductCandles get rates for a single product by product ID, grouped in buckets.
func (c *ApiClient) GetProductCandles(productId, start, end string, granularity Granularity) ([]ProductCandles, error) {
	u := c.makeV3Url(fmt.Sprintf("/brokerage/products/%s/candles?start=%s&end=%s&granularity=%s", productId, start, end, granularity))

	var data ProductCandlesData
	resp, err := c.httpClient.Get(u)
	if err != nil {
		return data.Candles, err
	}

	if !resp.IsSuccessState() {
		return data.Candles, ErrFailedToUnmarshal
	}

	err = resp.Unmarshal(&data)
	if err != nil {
		return data.Candles, err
	}

	return data.Candles, nil
}

type ProductCandlesData struct {
	Candles []ProductCandles `json:"candles"`
}

type ProductCandles struct {
	Start  string `json:"start"`
	Low    string `json:"low"`
	High   string `json:"high"`
	Open   string `json:"open"`
	Close  string `json:"close"`
	Volume string `json:"volume"`
}

// GetMarketTrades get snapshot information, by product ID, about the last trades (ticks), best bid/ask, and 24h volume.
func (c *ApiClient) GetMarketTrades(productId string, limit int32) (MarketTradesData, error) {
	u := c.makeV3Url(fmt.Sprintf("/brokerage/products/%s/ticker?limit=%d", productId, limit))

	var data MarketTradesData
	resp, err := c.httpClient.Get(u)
	if err != nil {
		return data, err
	}

	if !resp.IsSuccessState() {
		return data, ErrFailedToUnmarshal
	}

	err = resp.Unmarshal(&data)
	if err != nil {
		return MarketTradesData{}, err
	}

	return data, nil
}

type MarketTradesData struct {
	Trades  []MarketTrade `json:"trades"`
	BestBid string        `json:"best_bid"`
	BestAsk string        `json:"best_ask"`
}

// GetProductBook get a list of bids/asks for a single product. The amount of detail shown can be customized with the limit parameter.
func (c *ApiClient) GetProductBook(productId string, limit int32) (ProductBookData, error) {
	u := c.makeV3Url(fmt.Sprintf("/brokerage/product_book?product_id=%s&limit=%d", productId, limit))

	var data ProductBookData
	if res, err := c.get(u, &data); err != nil {
		return data, fmt.Errorf("%s", newErrorResponse(res))
	}
	return data, nil
}

type ProductBookData struct {
	PriceBook PriceBook `json:"pricebook"`
}

type PriceBook struct {
	ProductId string           `json:"product_id"`
	Bids      []PriceBookOrder `json:"bids"`
	Asks      []PriceBookOrder `json:"asks"`
	Time      string           `json:"time"`
}

type PriceBookOrder struct {
	Price string `json:"price"`
	Size  string `json:"size"`
}

// GetBestBidAsk get the best bid/ask for all products. A subset of all products can be returned instead by using the product_ids input.
func (c *ApiClient) GetBestBidAsk(productIds []string) (BestBidAskData, error) {
	query := strings.Join(productIds, "&product_ids=")
	if query != "" {
		query = "product_ids=" + query
	} else {
		return BestBidAskData{}, fmt.Errorf("no product ids provided")
	}

	u := c.makeV3Url(fmt.Sprintf("/brokerage/best_bid_ask?%s", query))
	var data BestBidAskData
	if res, err := c.get(u, &data); err != nil {
		return data, fmt.Errorf("%s", newErrorResponse(res))
	}
	return data, nil
}

type BestBidAskData struct {
	PriceBooks []PriceBook `json:"pricebooks"`
}
