package go_coinbasev3

import (
	"testing"
)

func TestNewWsFeedSubscription(t *testing.T) {
	feed := NewWsFeedSubscription("subscribe", []string{"BTC-USD"}, "level2", "1234567890", "secretkey")
	if feed.Type != "subscribe" {
		t.Errorf("Expected subscribe, got %s", feed.Type)
	}

	if len(feed.ProductIds) != 1 {
		t.Errorf("Expected 1 product id, got %d", len(feed.ProductIds))
	}

	if feed.Timestamp == 0 {
		t.Errorf("Expected timestamp to be set")
	}

	if feed.Signature == "" {
		t.Errorf("Expected signature to be set")
	}
}
