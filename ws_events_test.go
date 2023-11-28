package coinbasev3

import "testing"

func TestEvent_GetTickerEvent(t *testing.T) {
	evt := Event{
		Channel: "ticker",
		Events: []interface{}{
			struct {
				Type    string   `json:"type"`
				Tickers []Ticker `json:"tickers"`
			}{
				Type: "ticker",
				Tickers: []Ticker{
					{
						ProductId: "BTC-USD",
						Price:     "1000",
					},
				},
			},
		},
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
