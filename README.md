# Coinbase Advanced Trades API (v3)

This is a Go client for the Coinbase Advanced Trades API (v3). Work in progress.

## Installation

```bash
go get github.com/netr/go-coinbasev3
```

## Progress
- [X] Websocket Feed [Advanced Trade WebSocket Docs](https://docs.cloud.coinbase.com/advanced-trade-api/docs/ws-overview)
  - [ ] Create individual helpers for each channel type
    - [ ] `NewWsTickerSubscription(['ETH-USD', 'BTC-USD'], apiKey, secretKey, readCh, etc...)` - Useful for automatically typing the messages that are being received from the read channel.
- [ ] Advanced Trade API (V3) [Advanced Trade REST API Docs](https://docs.cloud.coinbase.com/advanced-trade-api/docs/rest-api-overview)
    - [ ] List Accounts
    - [ ] Get Account
    - [ ] Create Order
    - [ ] Cancel Orders
    - [ ] List Orders
    - [ ] List Fills
    - [ ] Get Order
    - [ ] Get Best Bid/Ask
    - [ ] Get Product Book
    - [ ] List Products
    - [ ] Get Product
    - [ ] Get Product Candles
    - [ ] Get Market Trades
    - [ ] Get Transactions Summary

## Websocket

The websocket client is a wrapper around the gorilla websocket with a few extra features to make it easier to use with the Coinbase Advanced Trade API.
The socket will automatically handle reconnections in the background to ensure you can run your sockets for long periods of time without worrying about the connection dropping. This is the reason why the configuration requires the channel subscriptions prior to initializing the client. In order for the connection to reconnect properly it will need to re-subscribe to the channels that were previously subscribed to.
```go
// create a ticker channel subscription. See: https://docs.cloud.coinbase.com/advanced-trade-api/docs/ws-channels
ticker := coinbasev3.NewWsChannelSub(
  []string{"ETH-USD", "BTC-USD"},
  coinbasev3.ChannelTypeTicker,
)
heartbeat := coinbasev3.NewWsChannelSub(
  []string{"ETH-USD", "ETH-EUR"},
  coinbasev3.ChannelTypeHeartbeats,
)

// create a channel to read messages from the websocket
readCh := make(chan []byte)

// create a minimal websocket config
wsConfig := coinbasev3.NewWsConfig("api_key", "secret_key", readCh, []coinbasev3.WsChannel{ticker, heartbeat})

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

## Run tests

```bash
go test ./... -v
```
