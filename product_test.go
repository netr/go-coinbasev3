package coinbasev3

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"net/http"
	"testing"
)

func TestApiClient_GetBestBidAsk(t *testing.T) {
	api := NewApiClient("api_key", "secret_key")
	productId := "BTC-USD"

	httpmock.ActivateNonDefault(api.client.GetClient())
	httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.coinbase.com/api/v3/brokerage/best_bid_ask?product_ids=%s", productId), func(request *http.Request) (*http.Response, error) {
		respBody := `{"pricebooks":[{"asks":[{"price":"2043.89","size":"0.57322806"}],"bids":[{"price":"2043.86","size":"2.24704213"}],"product_id":"ETH-USD","time":"2023-11-28T16:32:37.087555Z"},{"asks":[{"price":"14.408","size":"138.26"}],"bids":[{"price":"14.404","size":"3.8"}],"product_id":"LINK-USD","time":"2023-11-28T16:32:36.851091Z"},{"asks":[{"price":"37685.29","size":"0.06980337"}],"bids":[{"price":"37683.15","size":"0.03990523"}],"product_id":"BTC-USD","time":"2023-11-28T16:32:36.917395Z"}]}`
		resp := httpmock.NewStringResponse(http.StatusOK, respBody)
		resp.Header.Set("Content-Type", "application/json; charset=utf-8")
		return resp, nil
	})

	ask, err := api.GetBestBidAsk([]string{productId})
	if err != nil {
		t.Fatalf("Expected no error, got %s", err)
	}

	mapCheck := map[int]string{
		0: "ETH-USD",
		1: "LINK-USD",
		2: "BTC-USD",
	}
	for i, v := range ask.PriceBooks {
		if v.ProductId != mapCheck[i] {
			t.Errorf("Expected %s, got %s", mapCheck[i], v.ProductId)
		}
	}

	askCheck := map[int]string{
		0: "2043.89",
		1: "14.408",
		2: "37685.29",
	}
	for i, v := range ask.PriceBooks {
		if v.Asks[0].Price != askCheck[i] {
			t.Errorf("Expected %s, got %s", askCheck[i], v.Asks[0].Price)
		}
	}

	bidCheck := map[int]string{
		0: "2043.86",
		1: "14.404",
		2: "37683.15",
	}
	for i, v := range ask.PriceBooks {
		if v.Bids[0].Price != bidCheck[i] {
			t.Errorf("Expected %s, got %s", bidCheck[i], v.Bids[0].Price)
		}
	}
}

func TestApiClient_GetProductBook(t *testing.T) {
	api := NewApiClient("api_key", "secret_key")

	productId := "ETH-USD"
	var limit int32 = 4

	httpmock.ActivateNonDefault(api.client.GetClient())
	httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.coinbase.com/api/v3/brokerage/product_book?product_id=%s&limit=%d", productId, limit), func(request *http.Request) (*http.Response, error) {
		respBody := `{"pricebook":{"asks":[{"price":"2056.29","size":"1.67746848"},{"price":"2056.34","size":"0.350161"},{"price":"2056.35","size":"3.16704916"},{"price":"2056.36","size":"0.00529558"}],"bids":[{"price":"2056.13","size":"0.121583"},{"price":"2056.1","size":"0.14590056"},{"price":"2056.09","size":"1.32474108"},{"price":"2056.02","size":"0.5"}],"product_id":"ETH-USD","time":"2023-11-28T16:56:43.770106Z"}}`
		resp := httpmock.NewStringResponse(http.StatusOK, respBody)
		resp.Header.Set("Content-Type", "application/json; charset=utf-8")
		return resp, nil
	})

	data, err := api.GetProductBook(productId, limit)
	if err != nil {
		t.Fatalf("Expected no error, got %s", err)
	}

	checks := map[string]string{
		"2056.29": "1.67746848",
		"2056.34": "0.350161",
		"2056.35": "3.16704916",
		"2056.36": "0.00529558",
		"2056.13": "0.121583",
		"2056.1":  "0.14590056",
		"2056.09": "1.32474108",
		"2056.02": "0.5",
	}

	for _, v := range data.PriceBook.Asks {
		if checks[v.Price] != v.Size {
			t.Errorf("Expected %s, got %s", checks[v.Price], v.Size)
		}
	}

	if data.PriceBook.ProductId != "ETH-USD" {
		t.Errorf("Expected ETH-USD, got %s", data.PriceBook.ProductId)
	}

	if data.PriceBook.Time != "2023-11-28T16:56:43.770106Z" {
		t.Errorf("Expected 2023-11-28T16:56:43.770106Z, got %s", data.PriceBook.Time)
	}
}

func TestApiClient_GetMarketTrades(t *testing.T) {
	api := NewApiClient("api_key", "secret_key")

	productId := "ETH-USD"
	var limit int32 = 4

	httpmock.ActivateNonDefault(api.client.GetClient())
	httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.coinbase.com/api/v3/brokerage/products/%s/ticker?limit=%d", productId, limit), func(request *http.Request) (*http.Response, error) {
		respBody := `{"best_ask":"2058.33","best_bid":"2058.3","trades":[{"ask":"","bid":"","price":"2058.47","product_id":"ETH-USD","side":"SELL","size":"1.08006532","time":"2023-11-28T17:02:18.95103Z","trade_id":"480617715"},{"ask":"","bid":"","price":"2058.46","product_id":"ETH-USD","side":"SELL","size":"1.07776503","time":"2023-11-28T17:02:18.95103Z","trade_id":"480617714"},{"ask":"","bid":"","price":"2058.46","product_id":"ETH-USD","side":"SELL","size":"1.38405046","time":"2023-11-28T17:02:18.95103Z","trade_id":"480617713"},{"ask":"","bid":"","price":"2058.42","product_id":"ETH-USD","side":"SELL","size":"0.30674412","time":"2023-11-28T17:02:18.95103Z","trade_id":"480617712"}]}`
		resp := httpmock.NewStringResponse(http.StatusOK, respBody)
		resp.Header.Set("Content-Type", "application/json; charset=utf-8")
		return resp, nil
	})

	data, err := api.GetMarketTrades(productId, limit)
	if err != nil {
		t.Fatalf("Expected no error, got %s", err)
	}

	if data.BestAsk != "2058.33" {
		t.Errorf("Expected 2058.33, got %s", data.BestAsk)
	}
}

func TestApiClient_GetProductCandles(t *testing.T) {
	api := NewApiClient("api_key", "secret_key")

	productId := "ETH-USD"
	start := "1609459200"
	end := "1609545600"
	granularity := GranularityOneHour

	httpmock.ActivateNonDefault(api.client.GetClient())
	httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.coinbase.com/api/v3/brokerage/products/%s/candles?start=%s&end=%s&granularity=%s", productId, start, end, granularity), func(request *http.Request) (*http.Response, error) {
		respBody := `{"candles":[{"close":"721.8","high":"732","low":"715.22","open":"730.97","start":"1609545600","volume":"18729.1044335"},{"close":"730.99","high":"732.3","low":"729.4","open":"730.84","start":"1609542000","volume":"3961.27779606"},{"close":"730.97","high":"733.72","low":"727.42","open":"730.17","start":"1609538400","volume":"5213.35033145"},{"close":"730.18","high":"732.66","low":"728.8","open":"731.88","start":"1609534800","volume":"4663.13978788"},{"close":"731.87","high":"732.86","low":"724.65","open":"727.36","start":"1609531200","volume":"5550.06252723"},{"close":"727.45","high":"730.3","low":"722.69","open":"725.97","start":"1609527600","volume":"7873.28907073"},{"close":"725.99","high":"734.01","low":"717.1","open":"730.05","start":"1609524000","volume":"15402.59025014"},{"close":"730.26","high":"738.56","low":"729","open":"737.64","start":"1609520400","volume":"10272.37245532"},{"close":"737.65","high":"739.6","low":"735.05","open":"737.27","start":"1609516800","volume":"4945.18990967"},{"close":"737.27","high":"741.96","low":"735.3","open":"739.81","start":"1609513200","volume":"6095.96878078"},{"close":"739.86","high":"744.78","low":"738.51","open":"742.33","start":"1609509600","volume":"4612.50555526"},{"close":"742.25","high":"744.57","low":"734.7","open":"734.79","start":"1609506000","volume":"5391.8602088"},{"close":"734.79","high":"745.49","low":"734.06","open":"740.66","start":"1609502400","volume":"9812.80737865"},{"close":"740.71","high":"743.25","low":"738","open":"738.67","start":"1609498800","volume":"2942.76519355"},{"close":"738.26","high":"740.68","low":"734.97","open":"735.2","start":"1609495200","volume":"2548.99863614"},{"close":"735.2","high":"736.4","low":"730.62","open":"731.58","start":"1609491600","volume":"4385.0163215"},{"close":"731.57","high":"740.13","low":"726.71","open":"738.54","start":"1609488000","volume":"31868.7354194"},{"close":"738.44","high":"741.81","low":"735.27","open":"741.07","start":"1609484400","volume":"8189.28536611"},{"close":"741.06","high":"744.42","low":"738.38","open":"742.2","start":"1609480800","volume":"6550.59249879"},{"close":"742.02","high":"744.51","low":"740.65","open":"743.58","start":"1609477200","volume":"8546.58834584"},{"close":"743.57","high":"748.37","low":"740.49","open":"746.17","start":"1609473600","volume":"10497.25613281"},{"close":"746.17","high":"748.5","low":"743.47","open":"745.54","start":"1609470000","volume":"9219.05884386"},{"close":"745.54","high":"750","low":"743.75","open":"749.73","start":"1609466400","volume":"8446.51683824"},{"close":"749.74","high":"750","low":"735.2","open":"735.75","start":"1609462800","volume":"14000.72675479"},{"close":"735.69","high":"740.69","low":"731","open":"737.89","start":"1609459200","volume":"7070.48769111"}]}`
		resp := httpmock.NewStringResponse(http.StatusOK, respBody)
		resp.Header.Set("Content-Type", "application/json; charset=utf-8")
		return resp, nil
	})

	data, err := api.GetProductCandles(productId, start, end, granularity)
	if err != nil {
		t.Fatalf("Expected no error, got %s", err)
	}

	// one for each hour + 1 in the day
	if len(data) != 25 {
		t.Errorf("Expected 25 hours, got %d", len(data))
	}

}

func TestApiClient_GetProducts(t *testing.T) {
	api := NewApiClient("api_key", "secret_key")

	httpmock.ActivateNonDefault(api.client.GetClient())
	httpmock.RegisterResponder("GET", "https://api.exchange.coinbase.com/products", func(request *http.Request) (*http.Response, error) {
		respBody := `[{"auction_mode":false,"base_currency":"1INCH","base_increment":"0.01","cancel_only":false,"display_name":"1INCH/USD","fx_stablecoin":false,"high_bid_limit_percentage":"","id":"1INCH-USD","limit_only":false,"margin_enabled":false,"max_slippage_percentage":"0.03000000","min_market_funds":"1","post_only":false,"quote_currency":"USD","quote_increment":"0.001","status":"online","status_message":"","trading_disabled":false},{"auction_mode":false,"base_currency":"ARPA","base_increment":"0.1","cancel_only":false,"display_name":"ARPA/USDT","fx_stablecoin":false,"high_bid_limit_percentage":"","id":"ARPA-USDT","limit_only":false,"margin_enabled":false,"max_slippage_percentage":"0.03000000","min_market_funds":"1","post_only":false,"quote_currency":"USDT","quote_increment":"0.0001","status":"delisted","status_message":"","trading_disabled":true},{"auction_mode":false,"base_currency":"XTZ","base_increment":"0.01","cancel_only":false,"display_name":"XTZ/GBP","fx_stablecoin":false,"high_bid_limit_percentage":"","id":"XTZ-GBP","limit_only":false,"margin_enabled":false,"max_slippage_percentage":"0.03000000","min_market_funds":"0.72","post_only":false,"quote_currency":"GBP","quote_increment":"0.001","status":"online","status_message":"","trading_disabled":false}]`
		resp := httpmock.NewStringResponse(http.StatusOK, respBody)
		resp.Header.Set("Content-Type", "application/json; charset=utf-8")
		return resp, nil
	})

	data, err := api.GetProducts()
	if err != nil {
		t.Fatalf("Expected no error, got %s", err)
	}

	if len(data) != 3 {
		t.Errorf("Expected 3 products, got %d", len(data))
	}

	if data[0].Id != "1INCH-USD" {
		t.Errorf("Expected 1INCH-USD, got %s", data[0].Id)
	}
}

func TestApiClient_GetProduct(t *testing.T) {
	api := NewApiClient("api_key", "secret_key")

	productId := "ETH-USD"

	httpmock.ActivateNonDefault(api.client.GetClient())
	httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.coinbase.com/api/v3/brokerage/products/%s", productId), func(request *http.Request) (*http.Response, error) {
		respBody := `{"alias":"","alias_to":["ETH-USDC"],"auction_mode":false,"base_currency_id":"ETH","base_display_symbol":"ETH","base_increment":"0.00000001","base_max_size":"38000","base_min_size":"0.00022","base_name":"Ethereum","cancel_only":false,"fcm_trading_session_details":{"close_time":"","is_session_open":"","open_time":""},"future_product_details":{"contract_code":"","contract_display_name":"","contract_expiry":"","contract_expiry_timezone":"","contract_expiry_type":"","contract_root_unit":"","contract_size":"","group_description":"","group_short_description":"","perpetual_details":{"funding_rate":"","funding_time":"","open_interest":""},"risk_managed_by":"","venue":""},"is_disabled":false,"limit_only":false,"mid_market_price":"","new":false,"post_only":false,"price":"2055.34","price_increment":"0.01","price_percentage_change_24h":"1.55293466606717","product_id":"ETH-USD","product_type":"SPOT","quote_currency_id":"USD","quote_display_symbol":"USD","quote_increment":"0.01","quote_max_size":"50000000","quote_min_size":"1","quote_name":"US Dollar","status":"online","trading_disabled":false,"view_only":false,"volume_24h":"85000.44805841","volume_percentage_change_24h":"-22.83194559854808","watched":false}`
		resp := httpmock.NewStringResponse(http.StatusOK, respBody)
		resp.Header.Set("Content-Type", "application/json; charset=utf-8")
		return resp, nil
	})

	data, err := api.GetProduct(productId)
	if err != nil {
		t.Fatalf("Expected no error, got %s", err)
	}

	if data.ProductId != "ETH-USD" {
		t.Errorf("Expected ETH-USD, got %s", data.ProductId)
	}

	if data.BaseCurrencyId != "ETH" {
		t.Errorf("Expected ETH, got %s", data.BaseCurrencyId)
	}
}
