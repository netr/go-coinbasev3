package coinbasev3

import (
	"github.com/jarcoal/httpmock"
	"net/http"
	"testing"
)

func TestApiClient_GetTransactionSummary(t *testing.T) {
	api := NewApiClient("api_key", "secret_key")

	httpmock.ActivateNonDefault(api.client.GetClient())
	httpmock.RegisterResponder("GET", "https://api.coinbase.com/api/v3/brokerage/transaction_summary", func(request *http.Request) (*http.Response, error) {
		respBody := `{"advanced_trade_only_fees":0,"advanced_trade_only_volume":0,"coinbase_pro_fees":0,"coinbase_pro_volume":0,"fee_tier":{"aop_from":"","aop_to":"","maker_fee_rate":"0.006","pricing_tier":"Advanced 1","taker_fee_rate":"0.008","usd_from":"0","usd_to":"1000"},"goods_and_services_tax":{"rate":"","type":""},"margin_rate":{"value":""},"total_fees":0,"total_volume":0}`
		resp := httpmock.NewStringResponse(http.StatusOK, respBody)
		resp.Header.Set("Content-Type", "application/json; charset=utf-8")
		return resp, nil
	})

	data, err := api.GetTransactionSummary(TransactionSummaryRequest{})
	if err != nil {
		t.Fatalf("Expected no error, got %s", err)
	}

	if data.TotalVolume != 0 {
		t.Fatalf("Expected TotalVolume to be 0, got %d", data.TotalVolume)
	}

	if data.FeeTier.MakerFeeRate != "0.006" {
		t.Fatalf("Expected FeeTier.MakerFeeRate to be 0.006, got %s", data.FeeTier.MakerFeeRate)
	}
}
