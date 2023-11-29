package coinbasev3

import (
	"fmt"
	"strings"
)

type ProductType string

const (
	ProductTypeSpot   ProductType = "SPOT"
	ProductTypeFuture ProductType = "FUTURE"
)

type ContractExpiryType string

const (
	ContractExpiryTypeExpiring ContractExpiryType = "EXPIRING"
	ContractExpiryTypeUnknown  ContractExpiryType = "UNKNOWN_CONTRACT"
)

type TransactionSummaryRequest struct {
	// StartDate date-time RFC3339 format
	StartDate string
	// EndDate date-time RFC3339 format
	EndDate string
	// UserNativeCurrency string USD (default), EUR, GBP, etc. -- Only orders matching this native currency are returned.
	UserNativeCurrency string
	// ProductType string SPOT, FUTURE
	ProductType ProductType
	// ContractExpiryType string EXPIRING, UNKNOWN_CONTRACT (default) -- Only orders matching this contract expiry type are returned. Only filters response if ProductType is set to FUTURE.
	ContractExpiryType ContractExpiryType
}

// GetTransactionSummary get a summary of transactions with fee tiers, total volume, and fees.
func (c *ApiClient) GetTransactionSummary(req TransactionSummaryRequest) (TransactionSummaryData, error) {
	sb := strings.Builder{}
	if req.StartDate != "" {
		sb.WriteString(fmt.Sprintf("&start_date=%s", req.StartDate))
	}
	if req.EndDate != "" {
		sb.WriteString(fmt.Sprintf("&end_date=%s", req.EndDate))
	}
	if req.UserNativeCurrency != "" {
		sb.WriteString(fmt.Sprintf("&user_native_currency=%s", req.UserNativeCurrency))
	}
	if req.ProductType != "" {
		sb.WriteString(fmt.Sprintf("&product_type=%s", req.ProductType))
	}
	if req.ContractExpiryType != "" {
		sb.WriteString(fmt.Sprintf("&contract_expiry_type=%s", req.ContractExpiryType))
	}

	// Better way to do this?
	query := ""
	if sb.Len() > 0 {
		query = "?" + sb.String()[1:]
	}
	u := c.makeV3Url(fmt.Sprintf("/brokerage/transaction_summary%s", query))
	var data TransactionSummaryData
	if res, err := c.get(u, &data); err != nil {
		return data, fmt.Errorf("%s", newErrorResponse(res))
	}
	return data, nil
}

type TransactionSummaryData struct {
	TotalVolume             int                 `json:"total_volume"`
	TotalFees               int                 `json:"total_fees"`
	FeeTier                 FeeTier             `json:"fee_tier"`
	MarginRate              MarginRate          `json:"margin_rate"`
	GoodsAndServicesTax     GoodsAndServicesTax `json:"goods_and_services_tax"`
	AdvancedTradeOnlyVolume int                 `json:"advanced_trade_only_volume"`
	AdvancedTradeOnlyFees   int                 `json:"advanced_trade_only_fees"`
	CoinbaseProVolume       int                 `json:"coinbase_pro_volume"`
	CoinbaseProFees         int                 `json:"coinbase_pro_fees"`
}

type FeeTier struct {
	PricingTier  string `json:"pricing_tier"`
	UsdFrom      string `json:"usd_from"`
	UsdTo        string `json:"usd_to"`
	TakerFeeRate string `json:"taker_fee_rate"`
	MakerFeeRate string `json:"maker_fee_rate"`
	AopFrom      string `json:"aop_from"`
	AopTo        string `json:"aop_to"`
}

type MarginRate struct {
	Value string `json:"value"`
}

type GoodsAndServicesTax struct {
	Rate string `json:"rate"`
	Type string `json:"type"`
}
