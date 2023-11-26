package coinbasev3

import "testing"

func TestApiClient_GetProducts(t *testing.T) {
	api := NewApiClient("api_key", "secret_key")
	_, err := api.GetProducts()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}
