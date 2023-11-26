package coinbasev3

import (
	"testing"
)

func TestApiClient_GetFiatCurrencies(t *testing.T) {
	api := NewApiClient("api_key", "secret_key")
	_, err := api.GetFiatCurrencies()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func TestApiClient_GetCurrencies(t *testing.T) {
	api := NewApiClient("api_key", "secret_key")
	_, err := api.GetCurrencies()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func TestApiClient_GetExchangeRates(t *testing.T) {
	api := NewApiClient("api_key", "secret_key")
	_, err := api.GetExchangeRates("BTC")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}
