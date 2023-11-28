package coinbasev3

import (
	"testing"
)

func TestNewWsFeedSubscription(t *testing.T) {
	feed := NewWebsocketChannel(SubTypeSubscribe, ChannelTypeHeartbeats, []string{"BTC-USD"})
	if feed.Type != "subscribe" {
		t.Errorf("Expected subscribe, got %s", feed.Type)
	}

	if len(feed.ProductIds) != 1 {
		t.Errorf("Expected 1 product id, got %d", len(feed.ProductIds))
	}
}
