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

	data, err := api.GetListFills(ListFillsQuery{})
	if err != nil {
		t.Fatalf("Expected no error, got %s", err)
	}

	if len(data.Fills) != 0 {
		t.Fatalf("Expected Fills to be empty, got %d", len(data.Fills))
	}
}

func TestApiClient_GetListFills_WithFillData_Single(t *testing.T) {
	api := NewApiClient("api_key", "secret_key")

	query := ListFillsQuery{
		OrderId:                "0000-000000-000000",
		ProductId:              "4444-44444-444444",
		StartSequenceTimestamp: "2021-05-31T09:59:59Z",
		EndSequenceTimestamp:   "2021-05-31T10:59:59Z",
		Limit:                  100,
		Cursor:                 "789100",
	}

	u := fmt.Sprintf(
		"https://api.coinbase.com/api/v3/brokerage/orders/historical/fills?order_id=%s&product_id=%s&start_sequence_timestamp=%s&end_sequence_timestamp=%s&limit=%d&cursor=%s",
		query.OrderId,
		query.ProductId,
		query.StartSequenceTimestamp,
		query.EndSequenceTimestamp,
		query.Limit,
		query.Cursor,
	)

	httpmock.ActivateNonDefault(api.client.GetClient())
	httpmock.RegisterResponder("GET", u, func(request *http.Request) (*http.Response, error) {
		respBody := `{"fills":{"entry_id":"22222-2222222-22222222","trade_id":"1111-11111-111111","order_id":"0000-000000-000000","trade_time":"2021-05-31T09:59:59Z","trade_type":"FILL","price":"10000.00","size":"0.001","commission":"1.25","product_id":"BTC-USD","sequence_timestamp":"2021-05-31T09:58:59Z","liquidity_indicator":"UNKNOWN_LIQUIDITY_INDICATOR","size_in_quote":false,"user_id":"3333-333333-3333333","side":"UNKNOWN_ORDER_SIDE"},"cursor":"789100"}`
		resp := httpmock.NewStringResponse(http.StatusOK, respBody)
		resp.Header.Set("Content-Type", "application/json; charset=utf-8")
		return resp, nil
	})

	data, err := api.GetListFills(query)
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

func TestApiClient_GetListFills_WithFillData_Array(t *testing.T) {
	api := NewApiClient("api_key", "secret_key")

	query := ListFillsQuery{
		OrderId:                "0000-000000-000000",
		ProductId:              "4444-44444-444444",
		StartSequenceTimestamp: "2021-05-31T09:59:59Z",
		EndSequenceTimestamp:   "2021-05-31T10:59:59Z",
		Limit:                  100,
		Cursor:                 "789100",
	}

	u := fmt.Sprintf(
		"https://api.coinbase.com/api/v3/brokerage/orders/historical/fills?order_id=%s&product_id=%s&start_sequence_timestamp=%s&end_sequence_timestamp=%s&limit=%d&cursor=%s",
		query.OrderId,
		query.ProductId,
		query.StartSequenceTimestamp,
		query.EndSequenceTimestamp,
		query.Limit,
		query.Cursor,
	)

	httpmock.ActivateNonDefault(api.client.GetClient())
	httpmock.RegisterResponder("GET", u, func(request *http.Request) (*http.Response, error) {
		respBody := `{"fills":[{"entry_id":"22222-2222222-22222222","trade_id":"1111-11111-111111","order_id":"0000-000000-000000","trade_time":"2021-05-31T09:59:59Z","trade_type":"FILL","price":"10000.00","size":"0.001","commission":"1.25","product_id":"BTC-USD","sequence_timestamp":"2021-05-31T09:58:59Z","liquidity_indicator":"UNKNOWN_LIQUIDITY_INDICATOR","size_in_quote":false,"user_id":"3333-333333-3333333","side":"UNKNOWN_ORDER_SIDE"}],"cursor":"789100"}`
		resp := httpmock.NewStringResponse(http.StatusOK, respBody)
		resp.Header.Set("Content-Type", "application/json; charset=utf-8")
		return resp, nil
	})

	data, err := api.GetListFills(query)
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

func TestApiClient_GetListOrders_Empty(t *testing.T) {
	api := NewApiClient("api_key", "secret_key")

	httpmock.ActivateNonDefault(api.client.GetClient())
	httpmock.RegisterResponder("GET", "https://api.coinbase.com/api/v3/brokerage/orders/historical/batch", func(request *http.Request) (*http.Response, error) {
		respBody := `{"cursor":"","has_next":false,"orders":null,"sequence":""}`
		resp := httpmock.NewStringResponse(http.StatusOK, respBody)
		resp.Header.Set("Content-Type", "application/json; charset=utf-8")
		return resp, nil
	})

	data, err := api.GetListOrders(ListOrdersQuery{})
	if err != nil {
		t.Fatalf("Expected no error, got %s", err)
	}

	if len(data.Orders) != 0 {
		t.Fatalf("Expected Fills to be empty, got %d", len(data.Orders))
	}
}

func TestApiClient_GetListOrders_WithOrderDataArray(t *testing.T) {
	api := NewApiClient("api_key", "secret_key")

	httpmock.ActivateNonDefault(api.client.GetClient())
	httpmock.RegisterResponder("GET", "https://api.coinbase.com/api/v3/brokerage/orders/historical/batch", func(request *http.Request) (*http.Response, error) {
		respBody := `{"orders":[{"order_id":"0000-000000-000000","product_id":"BTC-USD","user_id":"2222-000000-000000","order_configuration":{"market_market_ioc":{"quote_size":"10.00","base_size":"0.001"},"limit_limit_gtc":{"base_size":"0.001","limit_price":"10000.00","post_only":false},"limit_limit_gtd":{"base_size":"0.001","limit_price":"10000.00","end_time":"2021-05-31T09:59:59Z","post_only":false},"stop_limit_stop_limit_gtc":{"base_size":"0.001","limit_price":"10000.00","stop_price":"20000.00","stop_direction":"UNKNOWN_STOP_DIRECTION"},"stop_limit_stop_limit_gtd":{"base_size":0.001,"limit_price":"10000.00","stop_price":"20000.00","end_time":"2021-05-31T09:59:59Z","stop_direction":"UNKNOWN_STOP_DIRECTION"}},"side":"UNKNOWN_ORDER_SIDE","client_order_id":"11111-000000-000000","status":"OPEN","time_in_force":"UNKNOWN_TIME_IN_FORCE","created_time":"2021-05-31T09:59:59Z","completion_percentage":"50","filled_size":"0.001","average_filled_price":"50","fee":"string","number_of_fills":"2","filled_value":"10000","pending_cancel":true,"size_in_quote":false,"total_fees":"5.00","size_inclusive_of_fees":false,"total_value_after_fees":"string","trigger_status":"UNKNOWN_TRIGGER_STATUS","order_type":"UNKNOWN_ORDER_TYPE","reject_reason":"REJECT_REASON_UNSPECIFIED","settled":"boolean","product_type":"SPOT","reject_message":"string","cancel_message":"string","order_placement_source":"RETAIL_ADVANCED","outstanding_hold_amount":"string","is_liquidation":"boolean","last_fill_time":"string","edit_history":[{"price":"string","size":"string","replace_accept_timestamp":"string"}]}],"sequence":"string","has_next":true,"cursor":"789100"}`
		resp := httpmock.NewStringResponse(http.StatusOK, respBody)
		resp.Header.Set("Content-Type", "application/json; charset=utf-8")
		return resp, nil
	})

	data, err := api.GetListOrders(ListOrdersQuery{})
	if err != nil {
		t.Fatalf("Expected no error, got %s", err)
	}

	if len(data.Orders) != 1 {
		t.Fatalf("Expected Orders to be empty, got %d", len(data.Orders))
	}
}

func TestApiClient_GetListOrders_WithOrderData_Single(t *testing.T) {
	api := NewApiClient("api_key", "secret_key")

	httpmock.ActivateNonDefault(api.client.GetClient())
	httpmock.RegisterResponder("GET", "https://api.coinbase.com/api/v3/brokerage/orders/historical/batch", func(request *http.Request) (*http.Response, error) {
		respBody := `{"orders":{"order_id":"0000-000000-000000","product_id":"BTC-USD","user_id":"2222-000000-000000","order_configuration":{"market_market_ioc":{"quote_size":"10.00","base_size":"0.001"},"limit_limit_gtc":{"base_size":"0.001","limit_price":"10000.00","post_only":false},"limit_limit_gtd":{"base_size":"0.001","limit_price":"10000.00","end_time":"2021-05-31T09:59:59Z","post_only":false},"stop_limit_stop_limit_gtc":{"base_size":"0.001","limit_price":"10000.00","stop_price":"20000.00","stop_direction":"UNKNOWN_STOP_DIRECTION"},"stop_limit_stop_limit_gtd":{"base_size":0.001,"limit_price":"10000.00","stop_price":"20000.00","end_time":"2021-05-31T09:59:59Z","stop_direction":"UNKNOWN_STOP_DIRECTION"}},"side":"UNKNOWN_ORDER_SIDE","client_order_id":"11111-000000-000000","status":"OPEN","time_in_force":"UNKNOWN_TIME_IN_FORCE","created_time":"2021-05-31T09:59:59Z","completion_percentage":"50","filled_size":"0.001","average_filled_price":"50","fee":"string","number_of_fills":"2","filled_value":"10000","pending_cancel":true,"size_in_quote":false,"total_fees":"5.00","size_inclusive_of_fees":false,"total_value_after_fees":"string","trigger_status":"UNKNOWN_TRIGGER_STATUS","order_type":"UNKNOWN_ORDER_TYPE","reject_reason":"REJECT_REASON_UNSPECIFIED","settled":"boolean","product_type":"SPOT","reject_message":"string","cancel_message":"string","order_placement_source":"RETAIL_ADVANCED","outstanding_hold_amount":"string","is_liquidation":"boolean","last_fill_time":"string","edit_history":[{"price":"string","size":"string","replace_accept_timestamp":"string"}]},"sequence":"string","has_next":true,"cursor":"789100"}`
		resp := httpmock.NewStringResponse(http.StatusOK, respBody)
		resp.Header.Set("Content-Type", "application/json; charset=utf-8")
		return resp, nil
	})

	data, err := api.GetListOrders(ListOrdersQuery{})
	if err != nil {
		t.Fatalf("Expected no error, got %s", err)
	}

	if len(data.Orders) != 1 {
		t.Fatalf("Expected Orders to be empty, got %d", len(data.Orders))
	}
}

func TestApiClient_GetOrder(t *testing.T) {
	api := NewApiClient("api_key", "secret_key")

	httpmock.ActivateNonDefault(api.client.GetClient())
	httpmock.RegisterResponder("GET", "https://api.coinbase.com/api/v3/brokerage/orders/historical/0000-000000-000000", func(request *http.Request) (*http.Response, error) {
		respBody := `{"order":{"order_id":"0000-000000-000000","product_id":"BTC-USD","user_id":"2222-000000-000000","order_configuration":{"market_market_ioc":{"quote_size":"10.00","base_size":"0.001"},"limit_limit_gtc":{"base_size":"0.001","limit_price":"10000.00","post_only":false},"limit_limit_gtd":{"base_size":"0.001","limit_price":"10000.00","end_time":"2021-05-31T09:59:59Z","post_only":false},"stop_limit_stop_limit_gtc":{"base_size":"0.001","limit_price":"10000.00","stop_price":"20000.00","stop_direction":"UNKNOWN_STOP_DIRECTION"},"stop_limit_stop_limit_gtd":{"base_size":0.001,"limit_price":"10000.00","stop_price":"20000.00","end_time":"2021-05-31T09:59:59Z","stop_direction":"UNKNOWN_STOP_DIRECTION"}},"side":"UNKNOWN_ORDER_SIDE","client_order_id":"11111-000000-000000","status":"OPEN","time_in_force":"UNKNOWN_TIME_IN_FORCE","created_time":"2021-05-31T09:59:59Z","completion_percentage":"50","filled_size":"0.001","average_filled_price":"50","fee":"string","number_of_fills":"2","filled_value":"10000","pending_cancel":true,"size_in_quote":false,"total_fees":"5.00","size_inclusive_of_fees":false,"total_value_after_fees":"string","trigger_status":"UNKNOWN_TRIGGER_STATUS","order_type":"UNKNOWN_ORDER_TYPE","reject_reason":"REJECT_REASON_UNSPECIFIED","settled":"boolean","product_type":"SPOT","reject_message":"string","cancel_message":"string","order_placement_source":"RETAIL_ADVANCED","outstanding_hold_amount":"string","is_liquidation":"boolean","last_fill_time":"string","edit_history":[{"price":"string","size":"string","replace_accept_timestamp":"string"}]}}`
		resp := httpmock.NewStringResponse(http.StatusOK, respBody)
		resp.Header.Set("Content-Type", "application/json; charset=utf-8")
		return resp, nil
	})

	data, err := api.GetOrder("0000-000000-000000")
	if err != nil {
		t.Fatalf("Expected no error, got %s", err)
	}

	if data.OrderId != "0000-000000-000000" {
		t.Fatalf("Expected OrderId to be 0000-000000-000000, got %s", data.OrderId)
	}
}
