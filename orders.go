package coinbasev3

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type OrderType string

const (
	OrderTypeMarket    OrderType = "MARKET"
	OrderTypeLimit     OrderType = "LIMIT"
	OrderTypeStop      OrderType = "STOP"
	OrderTypeStopLimit OrderType = "STOP_LIMIT"
	OrderTypeUnknown   OrderType = "UNKNOWN_ORDER_TYPE"
)

type OrderSide string

const (
	OrderSideBuy     OrderSide = "BUY"
	OrderSideSell    OrderSide = "SELL"
	OrderSideUnknown OrderSide = "UNKNOWN_ORDER_SIDE"
)

type OrderPlacementSource string

const (
	OrderPlacementSourceRetailAdvanced OrderPlacementSource = "RETAIL_ADVANCED"
	OrderPlacementSourceRetailSimple   OrderPlacementSource = "RETAIL_SIMPLE"
)

// ListFillsQuery represents the request parameters for the GetListFills function.
type ListFillsQuery struct {
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

// BuildQueryString creates a query string from the request parameters. If no parameters are set, an empty string is returned.
func (q ListFillsQuery) BuildQueryString() string {
	sb := strings.Builder{}
	if q.OrderId != "" {
		sb.WriteString(fmt.Sprintf("&order_id=%s", q.OrderId))
	}
	if q.ProductId != "" {
		sb.WriteString(fmt.Sprintf("&product_id=%s", q.ProductId))
	}
	if q.StartSequenceTimestamp != "" {
		sb.WriteString(fmt.Sprintf("&start_sequence_timestamp=%s", q.StartSequenceTimestamp))
	}
	if q.EndSequenceTimestamp != "" {
		sb.WriteString(fmt.Sprintf("&end_sequence_timestamp=%s", q.EndSequenceTimestamp))
	}
	if q.Limit > 0 {
		sb.WriteString(fmt.Sprintf("&limit=%d", q.Limit))
	}
	if q.Cursor != "" {
		sb.WriteString(fmt.Sprintf("&cursor=%s", q.Cursor))
	}

	// Better way to do this?
	if sb.Len() > 0 {
		return fmt.Sprintf("?%s", sb.String()[1:])
	}
	return ""
}

// GetListFills get a list of fills filtered by optional query parameters (product_id, order_id, etc).
func (c *ApiClient) GetListFills(q ListFillsQuery) (ListFillsData, error) {
	u := c.makeV3Url(fmt.Sprintf("/brokerage/orders/historical/fills%s", q.BuildQueryString()))
	var data ListFillsData
	if c.get(u, &data) != nil {
		return data, ErrFailedToUnmarshal
	}
	return data, nil
}

type ListFillsData struct {
	Fills  Fills  `json:"fills"`
	Cursor string `json:"cursor"`
}

type Fills []Fill

// UnmarshalJSON implements the json.Unmarshaler interface. Required because Coinbase returns an array of fills or a single fill object.
func (f *Fills) UnmarshalJSON(data []byte) error {
	// First, try unmarshaling into a slice
	var fillsSlice []Fill
	if err := json.Unmarshal(data, &fillsSlice); err == nil {
		*f = fillsSlice
		return nil
	}

	// If slice fails, try unmarshaling as a single object
	var singleFill Fill
	if err := json.Unmarshal(data, &singleFill); err == nil {
		*f = Fills{singleFill}
		return nil
	}

	// If both attempts fail, return an error
	return errors.New("fills should be an array or a single object")
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

// ListOrdersQuery represents the request parameters for the GetListOrders function.
type ListOrdersQuery struct {
	// ProductId Optional string of the product ID. Defaults to null, or fetch for all products.
	ProductId string `json:"product_id,omitempty"`
	// OrderStatus A list of order statuses.
	OrderStatus []string `json:"order_status"`
	// Limit A pagination limit with no default set. If has_next is true, additional orders are available to be fetched with pagination.
	Limit int32 `json:"limit,omitempty"`
	// StartDate Start date to fetch orders from, inclusive.
	StartDate string `json:"start_date"`
	// EndDate An optional end date for the query window, exclusive.
	EndDate string `json:"end_date,omitempty"`
	// OrderType Type of orders to return. Default is to return all order types.
	OrderType OrderType `json:"order_type,omitempty"`
	// OrderSide Only orders matching this side are returned. Default is to return all sides.
	OrderSide OrderSide `json:"order_side,omitempty"`
	// Cursor used for pagination. When provided, the response returns responses after this cursor.
	Cursor string `json:"cursor,omitempty"`
	// ProductType Only orders matching this product type are returned. Default is to return all product types.
	ProductType ProductType `json:"product_type,omitempty"`
	// OrderPlacementSource Only orders matching this placement source are returned. Default is to return RETAIL_ADVANCED placement source.
	OrderPlacementSource OrderPlacementSource `json:"order_placement_source,omitempty"`
	// ContractExpiryType Only orders matching this contract expiry type are returned. Filter is only applied if ProductType is set to FUTURE in the request.
	ContractExpiryType ContractExpiryType `json:"contract_expiry_type,omitempty"`
}

// BuildQueryString creates a query string from the request parameters. If no parameters are set, an empty string is returned.
func (q *ListOrdersQuery) BuildQueryString() string {
	var sb strings.Builder

	if q.ProductId != "" {
		sb.WriteString(fmt.Sprintf("&product_id=%s", q.ProductId))
	}
	if len(q.OrderStatus) > 0 {
		for _, status := range q.OrderStatus {
			sb.WriteString(fmt.Sprintf("&order_status=%s", status))
		}
	}
	if q.Limit != 0 {
		sb.WriteString(fmt.Sprintf("&limit=%d", q.Limit))
	}
	if q.StartDate != "" {
		sb.WriteString(fmt.Sprintf("&start_date=%s", q.StartDate))
	}
	if q.EndDate != "" {
		sb.WriteString(fmt.Sprintf("&end_date=%s", q.EndDate))
	}
	if q.OrderType != "" {
		sb.WriteString(fmt.Sprintf("&order_type=%s", q.OrderType))
	}
	if q.OrderSide != "" {
		sb.WriteString(fmt.Sprintf("&order_side=%s", q.OrderSide))
	}
	if q.Cursor != "" {
		sb.WriteString(fmt.Sprintf("&cursor=%s", q.Cursor))
	}
	if q.ProductType != "" {
		sb.WriteString(fmt.Sprintf("&product_type=%s", q.ProductType))
	}
	if q.OrderPlacementSource != "" {
		sb.WriteString(fmt.Sprintf("&order_placement_source=%s", q.OrderPlacementSource))
	}
	if q.ContractExpiryType != "" {
		sb.WriteString(fmt.Sprintf("&contract_expiry_type=%s", q.ContractExpiryType))
	}

	// Remove the first '&' for a clean query string
	if sb.Len() > 0 {
		return sb.String()[1:]
	}
	return ""
}

// GetListOrders get a list of orders filtered by optional query parameters (product_id, order_status, etc). Note: You cannot pair open orders with other order types. Example: order_status=OPEN,CANCELLED will return an error.
func (c *ApiClient) GetListOrders(q ListOrdersQuery) (ListOrdersData, error) {
	u := c.makeV3Url(fmt.Sprintf("/brokerage/orders/historical/batch%s", q.BuildQueryString()))

	var data ListOrdersData
	if c.get(u, &data) != nil {
		return data, ErrFailedToUnmarshal
	}
	return data, nil
}

type ListOrdersData struct {
	Orders   Orders `json:"orders"`
	Sequence string `json:"sequence"`
	HasNext  bool   `json:"has_next"`
	Cursor   string `json:"cursor"`
}

type Orders []Order

// UnmarshalJSON implements the json.Unmarshaler interface. Required because Coinbase returns an array of orders or a single order object.
func (o *Orders) UnmarshalJSON(data []byte) error {
	// First, try unmarshaling into a slice
	var ordersSlice []Order
	if err := json.Unmarshal(data, &ordersSlice); err == nil {
		*o = ordersSlice
		return nil
	}

	// If slice fails, try unmarshaling as a single object
	var singleOrder Order
	if err := json.Unmarshal(data, &singleOrder); err == nil {
		*o = Orders{singleOrder}
		return nil
	}

	// If both attempts fail, return an error
	return errors.New("orders should be an array or a single object")
}

type Order struct {
	OrderId               string             `json:"order_id"`
	ProductId             string             `json:"product_id"`
	UserId                string             `json:"user_id"`
	OrderConfiguration    OrderConfiguration `json:"order_configuration"`
	Side                  string             `json:"side"`
	ClientOrderId         string             `json:"client_order_id"`
	Status                string             `json:"status"`
	TimeInForce           string             `json:"time_in_force"`
	CreatedTime           time.Time          `json:"created_time"`
	CompletionPercentage  string             `json:"completion_percentage"`
	FilledSize            string             `json:"filled_size"`
	AverageFilledPrice    string             `json:"average_filled_price"`
	Fee                   string             `json:"fee"`
	NumberOfFills         string             `json:"number_of_fills"`
	FilledValue           string             `json:"filled_value"`
	PendingCancel         bool               `json:"pending_cancel"`
	SizeInQuote           bool               `json:"size_in_quote"`
	TotalFees             string             `json:"total_fees"`
	SizeInclusiveOfFees   bool               `json:"size_inclusive_of_fees"`
	TotalValueAfterFees   string             `json:"total_value_after_fees"`
	TriggerStatus         string             `json:"trigger_status"`
	OrderType             string             `json:"order_type"`
	RejectReason          string             `json:"reject_reason"`
	Settled               string             `json:"settled"`
	ProductType           string             `json:"product_type"`
	RejectMessage         string             `json:"reject_message"`
	CancelMessage         string             `json:"cancel_message"`
	OrderPlacementSource  string             `json:"order_placement_source"`
	OutstandingHoldAmount string             `json:"outstanding_hold_amount"`
	IsLiquidation         string             `json:"is_liquidation"`
	LastFillTime          string             `json:"last_fill_time"`
	EditHistory           []EditHistory      `json:"edit_history"`
}

type EditHistory struct {
	Price                  string `json:"price"`
	Size                   string `json:"size"`
	ReplaceAcceptTimestamp string `json:"replace_accept_timestamp"`
}

type OrderConfiguration struct {
	MarketMarketIoc       MarketMarketIoc       `json:"market_market_ioc"`
	LimitLimitGtc         LimitLimitGtc         `json:"limit_limit_gtc"`
	LimitLimitGtd         LimitLimitGtd         `json:"limit_limit_gtd"`
	StopLimitStopLimitGtc StopLimitStopLimitGtc `json:"stop_limit_stop_limit_gtc"`
	StopLimitStopLimitGtd StopLimitStopLimitGtd `json:"stop_limit_stop_limit_gtd"`
}

type MarketMarketIoc struct {
	QuoteSize string `json:"quote_size"`
	BaseSize  string `json:"base_size"`
}

type LimitLimitGtc struct {
	BaseSize   string `json:"base_size"`
	LimitPrice string `json:"limit_price"`
	PostOnly   bool   `json:"post_only"`
}

type LimitLimitGtd struct {
	BaseSize   string    `json:"base_size"`
	LimitPrice string    `json:"limit_price"`
	EndTime    time.Time `json:"end_time"`
	PostOnly   bool      `json:"post_only"`
}

type StopLimitStopLimitGtc struct {
	BaseSize      string `json:"base_size"`
	LimitPrice    string `json:"limit_price"`
	StopPrice     string `json:"stop_price"`
	StopDirection string `json:"stop_direction"`
}

type StopLimitStopLimitGtd struct {
	BaseSize      float64   `json:"base_size"`
	LimitPrice    string    `json:"limit_price"`
	StopPrice     string    `json:"stop_price"`
	EndTime       time.Time `json:"end_time"`
	StopDirection string    `json:"stop_direction"`
}

// GetOrder get a single order by order ID.
func (c *ApiClient) GetOrder(orderId string) (Order, error) {
	u := c.makeV3Url(fmt.Sprintf("/brokerage/orders/historical/%s", orderId))

	var data GetOrderData
	if c.get(u, &data) != nil {
		return data.Order, ErrFailedToUnmarshal
	}
	return data.Order, nil
}

type GetOrderData struct {
	Order Order `json:"order"`
}

type CreateOrderRequest struct {
	// ClientOrderId string Client set unique uuid for this order
	ClientOrderID string `json:"client_order_id"`
	// ProductId string The product this order was created for e.g. 'BTC-USD'
	ProductID string `json:"product_id"`
	// OrderType string Possible values: [UNKNOWN_ORDER_SIDE, BUY, SELL]
	Side               OrderSide          `json:"side"`
	OrderConfiguration OrderConfiguration `json:"order_configuration"`
}

func (req CreateOrderRequest) ToJson() ([]byte, error) {
	return json.Marshal(req)
}

// CreateOrder create an order with a specified product_id (asset-pair), side (buy/sell), etc.
func (c *ApiClient) CreateOrder(req CreateOrderRequest) (CreateOrderData, error) {
	var data CreateOrderData

	u := c.makeV3Url("/brokerage/orders")

	body, err := req.ToJson()
	if err != nil {
		return data, err
	}

	if c.post(u, body, &data) != nil {
		return data, ErrFailedToUnmarshal
	}
	return data, nil
}

type CreateOrderData struct {
	Success            bool                       `json:"success"`
	FailureReason      string                     `json:"failure_reason"`
	OrderId            string                     `json:"order_id"`
	SuccessResponse    CreateOrderSuccessResponse `json:"success_response"`
	ErrorResponse      CreatOrderErrorResponse    `json:"error_response"`
	OrderConfiguration OrderConfiguration         `json:"order_configuration"`
}

type CreatOrderErrorResponse struct {
	Error                 string `json:"error"`
	Message               string `json:"message"`
	ErrorDetails          string `json:"error_details"`
	PreviewFailureReason  string `json:"preview_failure_reason"`
	NewOrderFailureReason string `json:"new_order_failure_reason"`
}

type CreateOrderSuccessResponse struct {
	OrderId       string `json:"order_id"`
	ProductId     string `json:"product_id"`
	Side          string `json:"side"`
	ClientOrderId string `json:"client_order_id"`
}

// CancelOrders initiate cancel requests for one or more orders.
func (c *ApiClient) CancelOrders(orderIds []string) (CancelOrdersData, error) {
	var data CancelOrdersData

	u := c.makeV3Url("/brokerage/orders/batch_cancel")

	ords := struct {
		OrderIds []string `json:"order_ids"`
	}{OrderIds: orderIds}

	body, err := json.Marshal(ords)
	if err != nil {
		return data, err
	}

	if c.post(u, body, &data) != nil {
		return data, ErrFailedToUnmarshal
	}
	return data, nil
}

type CancelOrdersData struct {
	Results CancelOrderResults `json:"results"`
}

type CancelOrderResult struct {
	Success       bool   `json:"success"`
	FailureReason string `json:"failure_reason"`
	OrderId       string `json:"order_id"`
}

type CancelOrderResults []CancelOrderResult

// UnmarshalJSON implements the json.Unmarshaler interface. Required because Coinbase returns an array of orders or a single order object.
func (o *CancelOrderResults) UnmarshalJSON(data []byte) error {
	// First, try unmarshaling into a slice
	var ordersSlice []CancelOrderResult
	if err := json.Unmarshal(data, &ordersSlice); err == nil {
		*o = ordersSlice
		return nil
	}

	// If slice fails, try unmarshaling as a single object
	var singleOrder CancelOrderResult
	if err := json.Unmarshal(data, &singleOrder); err == nil {
		*o = CancelOrderResults{singleOrder}
		return nil
	}

	// If both attempts fail, return an error
	return errors.New("orders should be an array or a single object")
}
