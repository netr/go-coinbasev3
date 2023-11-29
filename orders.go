package coinbasev3

import (
	"fmt"
	"strings"
	"time"
)

type ListFillsRequest struct {
	// OrderId string ID of order
	OrderId string `json:"order_id"`
	// ProductId string The ID of the product this order was created for.
	ProductId string `json:"product_id"`
	// StartSequenceTimestamp date-time Start date. Only fills with a trade time at or after this start date are returned.
	StartSequenceTimestamp string `json:"start_sequence_timestamp"`
	// EndSequenceTimestamp date-time End date. Only fills with a trade time before this start date are returned.
	EndSequenceTimestamp string `json:"end_sequence_timestamp"`
	// Limit int64 Maximum number of fills to return in response. Defaults to 100.
	Limit int64 `json:"limit"`
	// Cursor string Cursor used for pagination. When provided, the response returns responses after this cursor.
	Cursor string `json:"cursor"`
}

// GetListFills get a list of fills filtered by optional query parameters (product_id, order_id, etc).
func (c *ApiClient) GetListFills(req ListFillsRequest) (ListFillsData, error) {
	sb := strings.Builder{}
	if req.OrderId != "" {
		sb.WriteString(fmt.Sprintf("&order_id=%s", req.OrderId))
	}
	if req.ProductId != "" {
		sb.WriteString(fmt.Sprintf("&product_id=%s", req.ProductId))
	}
	if req.StartSequenceTimestamp != "" {
		sb.WriteString(fmt.Sprintf("&start_sequence_timestamp=%s", req.StartSequenceTimestamp))
	}
	if req.EndSequenceTimestamp != "" {
		sb.WriteString(fmt.Sprintf("&end_sequence_timestamp=%s", req.EndSequenceTimestamp))
	}
	if req.Limit > 0 {
		sb.WriteString(fmt.Sprintf("&limit=%d", req.Limit))
	}
	if req.Cursor != "" {
		sb.WriteString(fmt.Sprintf("&cursor=%s", req.Cursor))
	}

	// Better way to do this?
	query := ""
	if sb.Len() > 0 {
		query = "?" + sb.String()[1:]
	}

	u := c.makeV3Url(fmt.Sprintf("/brokerage/orders/historical/fills%s", query))
	var data ListFillsData
	if c.get(u, &data) != nil {
		return data, ErrFailedToUnmarshal
	}
	return data, nil
}

type ListFillsData struct {
	Fills  []Fill `json:"fills"`
	Cursor string `json:"cursor"`
}

type Fill struct {
	EntryId            string    `json:"entry_id"`
	TradeId            string    `json:"trade_id"`
	OrderId            string    `json:"order_id"`
	TradeTime          time.Time `json:"trade_time"`
	TradeType          string    `json:"trade_type"`
	Price              string    `json:"price"`
	Size               string    `json:"size"`
	Commission         string    `json:"commission"`
	ProductId          string    `json:"product_id"`
	SequenceTimestamp  time.Time `json:"sequence_timestamp"`
	LiquidityIndicator string    `json:"liquidity_indicator"`
	SizeInQuote        bool      `json:"size_in_quote"`
	UserId             string    `json:"user_id"`
	Side               string    `json:"side"`
}
