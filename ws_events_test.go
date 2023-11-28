package coinbasev3

import (
	"testing"
	"time"
)

func TestEvent_GetTickerEvent(t *testing.T) {
	evt := Event{
		Channel: "ticker",
		Events: []interface{}{
			map[string]interface{}{
				"type": "update",
				"tickers": []Ticker{
					{
						ProductId: "BTC-USD",
						Price:     "1000",
					},
				},
			},
		},
	}

	if !evt.IsTickerEvent() {
		t.Errorf("Expected ticker event, got %s", evt.Channel)
	}

	tickerEvt, err := evt.GetTickerEvent()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(tickerEvt.Events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(tickerEvt.Events))
	}

	if tickerEvt.Events[0].Tickers[0].ProductId != "BTC-USD" {
		t.Errorf("Expected BTC-USD, got %s", tickerEvt.Events[0].Tickers[0].ProductId)
	}
}

func TestEvent_GetCandlesEvent(t *testing.T) {
	evt := Event{
		Channel: "candles",
		Events: []interface{}{
			map[string]interface{}{
				"type": "snapshot",
				"candles": []Candle{
					{
						Start:     "1688998200",
						High:      "1867.72",
						Low:       "1865.63",
						Open:      "1867.38",
						Close:     "1866.81",
						Volume:    "0.20269406",
						ProductId: "ETH-USD",
					},
				},
			},
		},
	}

	if !evt.IsCandlesEvent() {
		t.Errorf("Expected candles event, got %s", evt.Channel)
	}

	ne, err := evt.GetCandlesEvent()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(ne.Events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(ne.Events))
	}

	if ne.Events[0].Candles[0].ProductId != "ETH-USD" {
		t.Errorf("Expected ETH-USD, got %s", ne.Events[0].Candles[0].ProductId)
	}
}

func TestEvent_GetMarketTradesEvent(t *testing.T) {
	evt := Event{
		Channel: "market_trades",
		Events: []interface{}{
			map[string]interface{}{
				"type": "snapshot",
				"trades": []MarketTrade{
					{
						TradeId:   "123",
						ProductId: "BTC-USD",
						Price:     "1000",
						Size:      "1",
						Side:      "buy",
						Time:      time.Now(),
					},
				},
			},
		},
	}

	if !evt.IsMarketTradesEvent() {
		t.Errorf("Expected market trades event, got %s", evt.Channel)
	}

	ne, err := evt.GetMarketTradesEvent()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(ne.Events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(ne.Events))
	}

	if ne.Events[0].Trades[0].ProductId != "BTC-USD" {
		t.Errorf("Expected BTC-USD, got %s", ne.Events[0].Trades[0].ProductId)
	}
}

func TestEvent_GetHeartbeatsEvent(t *testing.T) {
	evt := Event{
		Channel: "heartbeats",
		Events: []interface{}{
			map[string]interface{}{
				"current_time":      "2023-06-23 20:31:56.121961769 +0000 UTC m=+91717.525857105",
				"heartbeat_counter": "3049",
			},
		},
	}

	if !evt.IsHeartbeatsEvent() {
		t.Errorf("Expected heartbeat event, got %s", evt.Channel)
	}

	ne, err := evt.GetHeartbeatsEvent()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(ne.Events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(ne.Events))
	}

	if ne.Events[0].HeartbeatCounter != "3049" {
		t.Errorf("Expected 3049, got %s", ne.Events[0].HeartbeatCounter)
	}
}

func TestEvent_GetStatusEvent(t *testing.T) {
	evt := Event{
		Channel: "status",
		Events: []interface{}{
			map[string]interface{}{
				"type": "snapshot",
				"products": []ProductStatus{
					{
						ProductType:    "SPOT",
						Id:             "BTC-USD",
						BaseCurrency:   "BTC",
						QuoteCurrency:  "USD",
						BaseIncrement:  "0.00000001",
						QuoteIncrement: "0.01",
						DisplayName:    "BTC/USD",
						Status:         "online",
						StatusMessage:  "",
						MinMarketFunds: "1",
					},
				},
			},
		},
	}

	if !evt.IsStatusEvent() {
		t.Errorf("Expected status event, got %s", evt.Channel)
	}

	ne, err := evt.GetStatusEvent()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(ne.Events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(ne.Events))
	}

	if ne.Events[0].Products[0].BaseIncrement != "0.00000001" {
		t.Errorf("Expected 0.00000001, got %s", ne.Events[0].Products[0].BaseIncrement)
	}
}

func TestEvent_GetUserEvent(t *testing.T) {
	evt := Event{
		Channel: "user",
		Events: []interface{}{
			map[string]interface{}{
				"type": "snapshot",
				"orders": []UserOrder{
					{
						OrderId:            "XXX",
						ClientOrderId:      "YYY",
						CumulativeQuantity: "0",
						LeavesQuantity:     "0.000994",
						AvgPrice:           "0",
						TotalFees:          "0",
						Status:             "OPEN",
						ProductId:          "BTC-USD",
						CreationTime:       time.Now(),
						OrderSide:          "BUY",
						OrderType:          "Limit",
					},
				},
			},
		},
	}

	if !evt.IsUserEvent() {
		t.Errorf("Expected user event, got %s", evt.Channel)
	}

	ne, err := evt.GetUserEvent()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(ne.Events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(ne.Events))
	}

	if ne.Events[0].Orders[0].ProductId != "BTC-USD" {
		t.Errorf("Expected BTC-USD, got %s", ne.Events[0].Orders[0].ProductId)
	}
}

func TestEvent_GetLevel2Event(t *testing.T) {
	evt := Event{
		Channel: "l2_data",
		Events: []interface{}{
			map[string]interface{}{
				"type":       "snapshot",
				"product_id": "BTC-USD",
				"updates": []Level2Update{
					{
						Side:        "bid",
						EventTime:   time.Now(),
						PriceLevel:  "21921.73",
						NewQuantity: "0.06317902",
					},
				},
			},
		},
	}

	if !evt.IsLevel2Event() {
		t.Errorf("Expected level2 event, got %s", evt.Channel)
	}

	ne, err := evt.GetLevel2Event()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(ne.Events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(ne.Events))
	}

	if ne.Events[0].ProductId != "BTC-USD" {
		t.Errorf("Expected BTC-USD, got %s", ne.Events[0].ProductId)
	}

	if ne.Events[0].Updates[0].PriceLevel != "21921.73" {
		t.Errorf("Expected 21921.73, got %s", ne.Events[0].Updates[0].PriceLevel)
	}
}
