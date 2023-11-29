# Coinbase Advanced Trades API (v3) ![Tests Status](https://github.com/netr/go-coinbasev3/actions/workflows/ci.yml/badge.svg) [![GoDoc](https://godoc.org/github.com/netr/go-coinbasev3?status.svg)](https://godoc.org/github.com/netr/go-coinbasev3)
This is a Go client for the Coinbase Advanced Trades API (v3). Work in progress. 

**Note:** The advanced trading API does not have a sandbox available at this time. All tests have been mocked to ensure the client works as expected. Once the sandbox is available for the advanced trading API the tests will be updated to use the sandbox.

## Installation

```bash
go get github.com/netr/go-coinbasev3
```

## Progress
Resource: [Rest API Pro Mapping](https://docs.cloud.coinbase.com/advanced-trade-api/docs/rest-api-pro-mapping)
- [X] Websocket Feed [Advanced Trade WebSocket Docs](https://docs.cloud.coinbase.com/advanced-trade-api/docs/ws-overview)
- [X] Public API (No authentication needed)
- [X] Advanced Trade API (V3) [Advanced Trade REST API Docs](https://docs.cloud.coinbase.com/advanced-trade-api/docs/rest-api-overview)
    - [x] List Accounts
    - [x] Get Account
    - [X] Create Order
    - [X] Cancel Orders
    - [X] Edit Order
    - [X] Edit Order Preview
    - [X] List Orders
    - [X] List Fills
    - [X] Get Order
    - [X] Get Best Bid/Ask
    - [X] Get Product Book
    - [X] List Products
    - [X] Get Product
    - [X] Get Product Candles
    - [X] Get Market Trades
    - [X] Get Transactions Summary
- [ ] Sign In with Coinbase API v2
  - [ ] Show an Account
  - [ ] List Transactions
  - [ ] Show Address
  - [ ] Create Address
  - [ ] Get Currencies
  - [ ] Deposit funds
  - [ ] List Payment Methods
  - [ ] List Transactions
  - [ ] Show a Transaction
  - [ ] Send Money
  - [ ] Withdraw Funds

## HTTP Client Usage

```go
// the client will automatically sign requests with the api_key and secret_key using req's OnBeforeRequest callback
client := coinbasev3.NewApiClient("api_key", "secret_key")

// product is a struct defined in the coinbasev3 package
product, err := client.GetProduct(productId)
if err != nil {
    panic("Failed to get product")
}
```

### Error Handling

The client will return an error if the response status code is not 200, the retries run out of attempts, or the []byte response can't properly marshal the result into json.

There is an error wrapper `coinbasev3.ResponseError` that will contain the error message and a marshaled coinbase error struct.

For most cases, the normal error message should contain enough information to determine what went wrong. However, if you need to access the coinbase error struct you can use the `errors.As` function to check if the error is a `coinbasev3.ResponseError` and then access the `CoinbaseError` field.

```go
res, err := client.EditOrder(coinbasev3.EditOrderRequest{
    OrderId: "d0c5340b-6d6c-49d9-b567-48c4bfca13d2",
    Size:    "1.00",
    Price:   "0.01",	
})
if err != nil {
    var myErr coinbasev3.ResponseError
    if errors.As(err, &myErr) {
        panic(err.(coinbasev3.ResponseError).CoinbaseError)
    } else {
        panic("Failed to edit order")
    }
}
```

### Changing base URL

The advanced trading API does not currently have a sandbox available. The sandbox is only available for the Coinbase API v2. The base URL can be changed to the sandbox URL for the Coinbase API v2 endpoints. Will need to revisit this once the sandbox is available for the advanced trading API. 

```go
client := coinbasev3.NewApiClient("api_key", "secret_key")
client.SetBaseUrlV3("https://api-public.sandbox.pro.coinbase.com")
client.SetBaseUrlV2("https://api-public.sandbox.coinbase.com")
client.SetBaseExchangeUrl("https://api.exchange.coinbase.com")
```

## Websocket

The websocket client is a wrapper around the gorilla websocket with a few extra features to make it easier to use with the Coinbase Advanced Trade API.
The socket will automatically handle reconnections in the background to ensure you can run your sockets for long periods of time without worrying about the connection dropping. This is the reason why the configuration requires the channel subscriptions prior to initializing the client. In order for the connection to reconnect properly it will need to re-subscribe to the channels that were previously subscribed to.

```go
// create a list of products to subscribe to 
// follows best practices mentioned in docs: https://docs.cloud.coinbase.com/advanced-trade-api/docs/ws-best-practices)
products := []string{"ETH-USD", "BTC-USD"}
// use one read channel for all products
readCh := make(chan []byte)
for _, product := range products {
  // create a ticker + heartbeat channel subscription. See: https://docs.cloud.coinbase.com/advanced-trade-api/docs/ws-channels
  wsChannels := []coinbasev3.WebsocketChannel{
    coinbasev3.NewTickerChannel([]string{product}),
    coinbasev3.NewHeartbeatsChannel([]string{product}),
  }
  
  // create a minimal websocket config
  wsConfig := coinbasev3.NewWsConfig("api_key", "secret_key", readCh, wsChannels)
  
  // create a websocket client with the config
  ws, err := coinbasev3.NewWsClient(wsConfig)
  if err != nil {
  panic("Failed to create Coinbase websocket client")
  }
  
  // open the websocket connection
  _, err = ws.Connect()
  if err != nil {
  panic("Failed to connect to Coinbase websocket")
  }
}

// read messages from the websocket
for {
    select {
        case msg := <-readCh:
            log.Println(string(msg))
    }
}

// will close the connection and stop the reconnection loop
ws.Shutdown()
```

### Reading messages from the websocket
The decision to use a read channel (instead of a callback) is to allow a more flexible dx, while also allowing the underlying socket to re-connect without any external interruptions or maintenance. The read channel will remain open for the entirety of the scope of the websocket client. If the websocket client is shutdown the read channel will be closed as well.

#### Note on the read channel data type
The read channel uses []byte instead of a default struct because the websocket client does not know what type of message it is receiving. The websocket client will not attempt to parse the message in any way. It is up to the developer to parse the message into the appropriate struct.

#### Parsing into pre-defined structs
For ease of use there are some helper functions to parse the messages into the appropriate struct. These helpers are not required to use the websocket client.

To utilize the helper functions you can parse the message into the default event struct, which is `coinbasev3.Event`.

Once you have the event struct you can use the `evt.IsTickerEvent()`, `evt.IsHeartbeatEvent()`, etc... to determine what type of event it is. Then you can use the `evt.GetTickerEvent()`, `evt.GetHeartbeatEvent()`, etc... to convert the default event into the appropriate struct.

```go
var evt coinbasev3.Event // default event struct (the events array is defined as an interface{})
err := json.Unmarshal(msg, &evt)
if err != nil {
  panic(err)
}

// check if the event is a ticker event
if evt.IsTickerEvent() {
  // convert the event into a ticker event struct
  tick, err := evt.GetTickerEvent()
  if err != nil {
      panic(err)
  }
}
```
## Run tests

```bash
go test ./... -v
```
