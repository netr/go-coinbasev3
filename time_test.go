package go_coinbasev3

import "testing"

func TestApiClient_GetServerTime(t *testing.T) {
	api := NewApiClient("api_key", "secret_key")
	_, err := api.GetServerTime()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}
