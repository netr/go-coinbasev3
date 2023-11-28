package coinbasev3

import (
	"github.com/mitchellh/mapstructure"
	"time"
)

// Event represents a standard event message from the websocket connection.
type Event struct {
	Channel     string        `json:"channel"`
	ClientId    string        `json:"client_id"`
	Timestamp   time.Time     `json:"timestamp"`
	SequenceNum int           `json:"sequence_num"`
	Events      []interface{} `json:"events"`
}

// IsTickerEvent returns true if the event is a ticker event.
// TickerBatch has the same JSON message schema as the ticker channel, except the channel field will have a value of ticker_batch.
func (e Event) IsTickerEvent() bool {
	return e.Channel == string(ChannelTypeTicker) || e.Channel == string(ChannelTypeTickerBatch)
}

// GetTickerEvent converts a generic event to a ticker/ticket_batch event. Returns an error if the event is not a ticker/ticket_batch event.
func (e Event) GetTickerEvent() (TickerEvent, error) {
	var evt TickerEvent
	evt.Event = e

	for _, ev := range e.Events {
		ne, ok := ev.(map[string]interface{})
		if !ok {
			return evt, ErrFailedToUnmarshal
		}

		var event TickerEventType
		err := mapstructure.Decode(ne, &event)
		if err != nil {
			return evt, err
		}

		evt.Events = append(evt.Events, TickerEventType{
			Type:    event.Type,
			Tickers: event.Tickers,
		})
	}

	return evt, nil
}

// IsHeartbeatsEvent returns true if the event is a heartbeat event.
func (e Event) IsHeartbeatsEvent() bool {
	return e.Channel == string(ChannelTypeHeartbeats)
}

// GetHeartbeatsEvent converts a generic event to a heartbeat event. Returns an error if the event is not a heartbeat event.
func (e Event) GetHeartbeatsEvent() (HeartbeatsEvent, error) {
	var evt HeartbeatsEvent
	for _, ev := range e.Events {
		ne, ok := ev.(map[string]interface{})
		if !ok {
			return evt, ErrFailedToUnmarshal
		}

		evt.Events = append(evt.Events, HeartbeatsEventType{
			CurrentTime:      ne["current_time"].(string),
			HeartbeatCounter: ne["heartbeat_counter"].(string),
		})
	}

	return evt, nil
}

// IsCandlesEvent returns true if the event is a candle event.
func (e Event) IsCandlesEvent() bool {
	return e.Channel == string(ChannelTypeCandles)
}

// GetCandlesEvent converts a generic event to a candle event. Returns an error if the event is not a candle event.
func (e Event) GetCandlesEvent() (CandlesEvent, error) {
	var evt CandlesEvent
	for _, ev := range e.Events {
		ne, ok := ev.(map[string]interface{})
		if !ok {
			return evt, ErrFailedToUnmarshal
		}

		var event CandlesEventType
		err := mapstructure.Decode(ne, &event)
		if err != nil {
			return evt, err
		}
		evt.Events = append(evt.Events, CandlesEventType{
			Type:    event.Type,
			Candles: event.Candles,
		})
	}

	return evt, nil
}

// IsMarketTradesEvent returns true if the event is a market trades event.
func (e Event) IsMarketTradesEvent() bool {
	return e.Channel == string(ChannelTypeMarketTrades)
}

// GetMarketTradesEvent converts a generic event to a market trades event. Returns an error if the event is not a market trades event.
func (e Event) GetMarketTradesEvent() (MarketTradesEvent, error) {
	var evt MarketTradesEvent
	for _, ev := range e.Events {
		ne, ok := ev.(map[string]interface{})
		if !ok {
			return evt, ErrFailedToUnmarshal
		}

		var event MarketTradesEventType
		err := mapstructure.Decode(ne, &event)
		if err != nil {
			return evt, err
		}
		evt.Events = append(evt.Events, MarketTradesEventType{
			Type:   event.Type,
			Trades: event.Trades,
		})
	}

	return evt, nil
}

// IsStatusEvent returns true if the event is a status event.
func (e Event) IsStatusEvent() bool {
	return e.Channel == string(ChannelTypeStatus)
}

// GetStatusEvent converts a generic event to a status event. Returns an error if the event is not a status event.
func (e Event) GetStatusEvent() (StatusEvent, error) {
	var evt StatusEvent
	for _, ev := range e.Events {
		ne, ok := ev.(map[string]interface{})
		if !ok {
			return evt, ErrFailedToUnmarshal
		}

		var event StatusEventType
		err := mapstructure.Decode(ne, &event)
		if err != nil {
			return evt, err
		}
		evt.Events = append(evt.Events, StatusEventType{
			Type:     event.Type,
			Products: event.Products,
		})
	}

	return evt, nil
}

// IsLevel2Event returns true if the event is a level2 event.
func (e Event) IsLevel2Event() bool {
	return e.Channel == string(ChannelTypeLevel2)
}

// GetLevel2Event converts a generic event to a level 2 event. Returns an error if the event is not level 2 event.
func (e Event) GetLevel2Event() (Level2Event, error) {
	var evt Level2Event
	for _, ev := range e.Events {
		ne, ok := ev.(map[string]interface{})
		if !ok {
			return evt, ErrFailedToUnmarshal
		}

		var event Level2EventType
		err := mapstructure.Decode(ne, &event)
		if err != nil {
			return evt, err
		}

		evt.Events = append(evt.Events, Level2EventType{
			Type:      event.Type,
			ProductId: event.ProductId,
			Updates:   event.Updates,
		})
	}

	return evt, nil
}

// IsUserEvent returns true if the event is a user's order event.
func (e Event) IsUserEvent() bool {
	return e.Channel == string(ChannelTypeUser)
}

// GetUserEvent converts a generic event to a user's order event. Returns an error if the event is not a user's order event.
func (e Event) GetUserEvent() (UserEvent, error) {
	var evt UserEvent
	for _, ev := range e.Events {
		ne, ok := ev.(map[string]interface{})
		if !ok {
			return evt, ErrFailedToUnmarshal
		}

		var event UserEventType
		err := mapstructure.Decode(ne, &event)
		if err != nil {
			return evt, err
		}
		evt.Events = append(evt.Events, UserEventType{
			Type:   event.Type,
			Orders: event.Orders,
		})
	}

	return evt, nil
}

// TickerEvent represents a ticker event message from the websocket connection.
type TickerEvent struct {
	Event
	Events []TickerEventType `json:"events"`
}

type TickerEventType struct {
	Type    string   `json:"type"`
	Tickers []Ticker `json:"tickers"`
}

// Ticker represents a ticker from the websocket connection.
type Ticker struct {
	Type               string `json:"type" mapstructure:"type"`
	ProductId          string `json:"product_id" mapstructure:"product_id"`
	Price              string `json:"price" mapstructure:"price"`
	Volume24H          string `json:"volume_24_h" mapstructure:"volume_24_h"`
	Low24H             string `json:"low_24_h" mapstructure:"low_24_h"`
	High24H            string `json:"high_24_h" mapstructure:"high_24_h"`
	Low52W             string `json:"low_52_w" mapstructure:"low_52_w"`
	High52W            string `json:"high_52_w" mapstructure:"high_52_w"`
	PricePercentChg24H string `json:"price_percent_chg_24_h" mapstructure:"price_percent_chg_24_h"`
}

type HeartbeatsEvent struct {
	Event
	Events []HeartbeatsEventType `json:"events"`
}

type HeartbeatsEventType struct {
	CurrentTime      string `json:"current_time" mapstructure:"current_time"`
	HeartbeatCounter string `json:"heartbeat_counter" mapstructure:"heartbeat_counter"`
}

type CandlesEvent struct {
	Event
	Events []CandlesEventType `json:"events"`
}

type CandlesEventType struct {
	Type    string   `json:"type"`
	Candles []Candle `json:"candles"`
}

type Candle struct {
	Start     string `json:"start" mapstructure:"start"`
	High      string `json:"high" mapstructure:"high"`
	Low       string `json:"low" mapstructure:"low"`
	Open      string `json:"open" mapstructure:"open"`
	Close     string `json:"close" mapstructure:"close"`
	Volume    string `json:"volume" mapstructure:"volume"`
	ProductId string `json:"product_id" mapstructure:"product_id"`
}

type MarketTradesEvent struct {
	Event
	Events []MarketTradesEventType `json:"events"`
}

type MarketTradesEventType struct {
	Type   string        `json:"type"`
	Trades []MarketTrade `json:"trades"`
}

type MarketTrade struct {
	TradeId   string    `json:"trade_id" mapstructure:"trade_id"`
	ProductId string    `json:"product_id" mapstructure:"product_id"`
	Price     string    `json:"price" mapstructure:"price"`
	Size      string    `json:"size" mapstructure:"size"`
	Side      string    `json:"side" mapstructure:"side"`
	Time      time.Time `json:"time" mapstructure:"time"`
	Bid       string    `json:"bid"  mapstructure:"bid"`
	Ask       string    `json:"ask"  mapstructure:"ask"`
}

type StatusEvent struct {
	Event
	Events []StatusEventType `json:"events"`
}

type StatusEventType struct {
	Type     string          `json:"type"`
	Products []ProductStatus `json:"products"`
}

type ProductStatus struct {
	ProductType    string `json:"product_type" mapstructure:"product_type"`
	Id             string `json:"id" mapstructure:"id"`
	BaseCurrency   string `json:"base_currency" mapstructure:"base_currency"`
	QuoteCurrency  string `json:"quote_currency" mapstructure:"quote_currency"`
	BaseIncrement  string `json:"base_increment" mapstructure:"base_increment"`
	QuoteIncrement string `json:"quote_increment" mapstructure:"quote_increment"`
	DisplayName    string `json:"display_name" mapstructure:"display_name"`
	Status         string `json:"status" mapstructure:"status"`
	StatusMessage  string `json:"status_message" mapstructure:"status_message"`
	MinMarketFunds string `json:"min_market_funds" mapstructure:"min_market_funds"`
}

type Level2Event struct {
	Event
	Events []Level2EventType `json:"events"`
}

type Level2EventType struct {
	Type      string         `json:"type" mapstructure:"type"`
	ProductId string         `json:"product_id" mapstructure:"product_id"`
	Updates   []Level2Update `json:"updates" mapstructure:"updates"`
}

type Level2Update struct {
	Side        string    `json:"side" mapstructure:"side"`
	EventTime   time.Time `json:"event_time" mapstructure:"event_time"`
	PriceLevel  string    `json:"price_level" mapstructure:"price_level"`
	NewQuantity string    `json:"new_quantity" mapstructure:"new_quantity"`
}

type UserEvent struct {
	Event
	Events []UserEventType `json:"events"`
}

type UserEventType struct {
	Type   string      `json:"type"`
	Orders []UserOrder `json:"orders"`
}

type UserOrder struct {
	OrderId            string    `json:"order_id" mapstructure:"order_id"`
	ClientOrderId      string    `json:"client_order_id" mapstructure:"client_order_id"`
	CumulativeQuantity string    `json:"cumulative_quantity" mapstructure:"cumulative_quantity"`
	LeavesQuantity     string    `json:"leaves_quantity" mapstructure:"leaves_quantity"`
	AvgPrice           string    `json:"avg_price" mapstructure:"avg_price"`
	TotalFees          string    `json:"total_fees" mapstructure:"total_fees"`
	Status             string    `json:"status" mapstructure:"status"`
	ProductId          string    `json:"product_id" mapstructure:"product_id"`
	CreationTime       time.Time `json:"creation_time" mapstructure:"creation_time"`
	OrderSide          string    `json:"order_side" mapstructure:"order_side"`
	OrderType          string    `json:"order_type" mapstructure:"order_type"`
}
