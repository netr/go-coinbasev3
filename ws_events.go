package coinbasev3

import "time"

// Event represents a standard event message from the websocket connection.
type Event struct {
	Channel     string        `json:"channel"`
	ClientId    string        `json:"client_id"`
	Timestamp   time.Time     `json:"timestamp"`
	SequenceNum int           `json:"sequence_num"`
	Events      []interface{} `json:"events"`
}

// TickerEvent represents a ticker event message from the websocket connection.
type TickerEvent struct {
	Event
	Events []struct {
		Type    string   `json:"type"`
		Tickers []Ticker `json:"tickers"`
	} `json:"events"`
}

// Ticker represents a ticker from the websocket connection.
type Ticker struct {
	Type               string `json:"type"`
	ProductId          string `json:"product_id"`
	Price              string `json:"price"`
	Volume24H          string `json:"volume_24_h"`
	Low24H             string `json:"low_24_h"`
	High24H            string `json:"high_24_h"`
	Low52W             string `json:"low_52_w"`
	High52W            string `json:"high_52_w"`
	PricePercentChg24H string `json:"price_percent_chg_24_h"`
}

type HeartbeatsEvent struct {
	Event
	Events []struct {
		CurrentTime      string `json:"current_time"`
		HeartbeatCounter string `json:"heartbeat_counter"`
	} `json:"events"`
}

type CandlesEvent struct {
	Event
	Events []struct {
		Type    string   `json:"type"`
		Candles []Candle `json:"candles"`
	} `json:"events"`
}

type Candle struct {
	Start     string `json:"start"`
	High      string `json:"high"`
	Low       string `json:"low"`
	Open      string `json:"open"`
	Close     string `json:"close"`
	Volume    string `json:"volume"`
	ProductId string `json:"product_id"`
}

type MarketTradesEvent struct {
	Event
	Events []struct {
		Type   string        `json:"type"`
		Trades []MarketTrade `json:"trades"`
	} `json:"events"`
}

type MarketTrade struct {
	TradeId   string    `json:"trade_id"`
	ProductId string    `json:"product_id"`
	Price     string    `json:"price"`
	Size      string    `json:"size"`
	Side      string    `json:"side"`
	Time      time.Time `json:"time"`
}

type StatusEvent struct {
	Event
	Events []struct {
		Type     string          `json:"type"`
		Products []ProductStatus `json:"products"`
	} `json:"events"`
}

type ProductStatus struct {
	ProductType    string `json:"product_type"`
	Id             string `json:"id"`
	BaseCurrency   string `json:"base_currency"`
	QuoteCurrency  string `json:"quote_currency"`
	BaseIncrement  string `json:"base_increment"`
	QuoteIncrement string `json:"quote_increment"`
	DisplayName    string `json:"display_name"`
	Status         string `json:"status"`
	StatusMessage  string `json:"status_message"`
	MinMarketFunds string `json:"min_market_funds"`
}

type Level2Event struct {
	Event
	Events []struct {
		Type      string         `json:"type"`
		ProductId string         `json:"product_id"`
		Updates   []Level2Update `json:"updates"`
	} `json:"events"`
}

type Level2Update struct {
	Side        string    `json:"side"`
	EventTime   time.Time `json:"event_time"`
	PriceLevel  string    `json:"price_level"`
	NewQuantity string    `json:"new_quantity"`
}

type UserEvent struct {
	Event
	Events []struct {
		Type   string      `json:"type"`
		Orders []UserOrder `json:"orders"`
	} `json:"events"`
}

type UserOrder struct {
	OrderId            string    `json:"order_id"`
	ClientOrderId      string    `json:"client_order_id"`
	CumulativeQuantity string    `json:"cumulative_quantity"`
	LeavesQuantity     string    `json:"leaves_quantity"`
	AvgPrice           string    `json:"avg_price"`
	TotalFees          string    `json:"total_fees"`
	Status             string    `json:"status"`
	ProductId          string    `json:"product_id"`
	CreationTime       time.Time `json:"creation_time"`
	OrderSide          string    `json:"order_side"`
	OrderType          string    `json:"order_type"`
}
