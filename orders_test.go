package coinbasev3

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"net/http"
	"testing"
)

func TestApiClient_GetListFills_Empty(t *testing.T) {
	api := NewApiClient("api_key", "secret_key")

	httpmock.ActivateNonDefault(api.client.GetClient())
	httpmock.RegisterResponder("GET", "https://api.coinbase.com/api/v3/brokerage/orders/historical/fills", func(request *http.Request) (*http.Response, error) {
		respBody := `{"cursor":"","fills":[]}`
		resp := httpmock.NewStringResponse(http.StatusOK, respBody)
		resp.Header.Set("Content-Type", "application/json; charset=utf-8")
		return resp, nil
	})

	data, err := api.GetListFills(ListFillsRequest{})
	if err != nil {
		t.Fatalf("Expected no error, got %s", err)
	}

	if len(data.Fills) != 0 {
		t.Fatalf("Expected Fills to be empty, got %d", len(data.Fills))
	}
}

func TestApiClient_GetListFills_WithFillData(t *testing.T) {
	api := NewApiClient("api_key", "secret_key")

	req := ListFillsRequest{
		OrderId:                "0000-000000-000000",
		ProductId:              "4444-44444-444444",
		StartSequenceTimestamp: "2021-05-31T09:59:59Z",
		EndSequenceTimestamp:   "2021-05-31T10:59:59Z",
		Limit:                  100,
		Cursor:                 "789100",
	}

	u := fmt.Sprintf(
		"https://api.coinbase.com/api/v3/brokerage/orders/historical/fills?order_id=%s&product_id=%s&start_sequence_timestamp=%s&end_sequence_timestamp=%s&limit=%d&cursor=%s",
		req.OrderId,
		req.ProductId,
		req.StartSequenceTimestamp,
		req.EndSequenceTimestamp,
		req.Limit,
		req.Cursor,
	)

	httpmock.ActivateNonDefault(api.client.GetClient())
	httpmock.RegisterResponder("GET", u, func(request *http.Request) (*http.Response, error) {
		respBody := `{"fills":[{"entry_id":"22222-2222222-22222222","trade_id":"1111-11111-111111","order_id":"0000-000000-000000","trade_time":"2021-05-31T09:59:59Z","trade_type":"FILL","price":"10000.00","size":"0.001","commission":"1.25","product_id":"BTC-USD","sequence_timestamp":"2021-05-31T09:58:59Z","liquidity_indicator":"UNKNOWN_LIQUIDITY_INDICATOR","size_in_quote":false,"user_id":"3333-333333-3333333","side":"UNKNOWN_ORDER_SIDE"}],"cursor":"789100"}`
		resp := httpmock.NewStringResponse(http.StatusOK, respBody)
		resp.Header.Set("Content-Type", "application/json; charset=utf-8")
		return resp, nil
	})

	data, err := api.GetListFills(req)
	if err != nil {
		t.Fatalf("Expected no error, got %s", err)
	}

	if len(data.Fills) != 1 {
		t.Fatalf("Expected Fills to be empty, got %d", len(data.Fills))
	}

	if data.Fills[0].EntryId != "22222-2222222-22222222" {
		t.Fatalf("Expected Fills[0].EntryId to be 22222-2222222-22222222, got %s", data.Fills[0].EntryId)
	}
}
