package coinbasev3

import "testing"

func TestApiClient_GetBuyPrice(t *testing.T) {
	api := NewApiClient("api_key", "secret_key")
	_, err := api.GetBuyPrice("BTC-USD")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func TestApiClient_GetSellPrice(t *testing.T) {
	api := NewApiClient("api_key", "secret_key")
	_, err := api.GetSellPrice("BTC-USD")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func TestApiClient_GetSpotPrice(t *testing.T) {
	api := NewApiClient("api_key", "secret_key")
	_, err := api.GetSpotPrice("BTC-USD")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}
