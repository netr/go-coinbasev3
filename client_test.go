package coinbasev3

import "testing"

func TestNewApiClient(t *testing.T) {
	api := NewApiClient("api_key", "secret_key")
	if api == nil {
		t.Errorf("Expected api to be initialized")
	}
	if api.httpClient == nil {
		t.Errorf("Expected client to be initialized")
	}

	// check base urls
	if api.baseUrlV3 != "https://api.coinbase.com/api/v3" {
		t.Errorf("Expected base url to be https://api.coinbase.com/api/v3, got %s", api.baseUrlV3)
	}
	if api.baseUrlV2 != "https://api.coinbase.com/api/v2" {
		t.Errorf("Expected base url to be https://api.coinbase.com/api/v2, got %s", api.baseUrlV2)
	}
	if api.baseExchangeUrl != "https://api.exchange.coinbase.com" {
		t.Errorf("Expected base url to be https://api.exchange.coinbase.com, got %s", api.baseExchangeUrl)
	}
}

func TestNewApiClient_WithCustomClient(t *testing.T) {
	mock := NewMockHttpClient(nil)
	api := NewApiClient("api_key", "secret_key", mock)
	if api == nil {
		t.Errorf("Expected api to be initialized")
	}
	if api.httpClient.GetClient() != nil {
		t.Errorf("Expected client to be nil from the mock")
	}

	// check base urls
	if api.baseUrlV3 != "https://api.coinbase.com/api/v3" {
		t.Errorf("Expected base url to be https://api.coinbase.com/api/v3, got %s", api.baseUrlV3)
	}
	if api.baseUrlV2 != "https://api.coinbase.com/api/v2" {
		t.Errorf("Expected base url to be https://api.coinbase.com/api/v2, got %s", api.baseUrlV2)
	}
	if api.baseExchangeUrl != "https://api.exchange.coinbase.com" {
		t.Errorf("Expected base url to be https://api.exchange.coinbase.com, got %s", api.baseExchangeUrl)
	}
}

func TestApiClient_SetBaseExchangeUrl(t *testing.T) {
	api := NewApiClient("api_key", "secret_key")
	if api == nil {
		t.Errorf("Expected api to be initialized")
	}

	api.SetBaseExchangeUrl("https://testbase.com")
	if api.baseExchangeUrl != "https://testbase.com" {
		t.Errorf("Expected base url to be https://testbase.com, got %s", api.baseExchangeUrl)
	}
}

func TestApiClient_SetBaseV2Url(t *testing.T) {
	api := NewApiClient("api_key", "secret_key")
	if api == nil {
		t.Errorf("Expected api to be initialized")
	}

	api.SetBaseUrlV2("https://testbase.com")
	if api.baseUrlV2 != "https://testbase.com" {
		t.Errorf("Expected base url to be https://testbase.com, got %s", api.baseUrlV2)
	}
}

func TestApiClient_SetBaseV3Url(t *testing.T) {
	api := NewApiClient("api_key", "secret_key")
	if api == nil {
		t.Errorf("Expected api to be initialized")
	}

	api.SetBaseUrlV3("https://testbase.com")
	if api.baseUrlV3 != "https://testbase.com" {
		t.Errorf("Expected base url to be https://testbase.com, got %s", api.baseUrlV3)
	}
}

func TestApiClient_SetSandboxUrls(t *testing.T) {
	api := NewApiClient("api_key", "secret_key")
	if api == nil {
		t.Errorf("Expected api to be initialized")
	}
	api.SetSandboxUrls()

	// check base urls. TODO: Fix this when sandbox is available
	checks := map[string]string{
		"v3":       "https://api-public.sandbox.pro.coinbase.com",
		"v2":       "https://api-public.sandbox.pro.coinbase.com",
		"exchange": "https://api-public.sandbox.exchange.coinbase.com",
	}

	for k, v := range checks {
		switch k {
		case "v3":
			if api.baseUrlV3 != v {
				t.Errorf("Expected base url to be %s, got %s", v, api.baseUrlV3)
			}
		case "v2":
			if api.baseUrlV2 != v {
				t.Errorf("Expected base url to be %s, got %s", v, api.baseUrlV2)
			}
		case "exchange":
			if api.baseExchangeUrl != v {
				t.Errorf("Expected base url to be %s, got %s", v, api.baseExchangeUrl)
			}
		}
	}
}
