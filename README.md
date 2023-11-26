# Coinbase Advanced Trades API (v3)

This is a Go client for the Coinbase Advanced Trades API (v3). Work in progress.

## Installation

```bash
go get github.com/netr/go-coinbasev3
```

## Progress
- [X] Websocket Feed [Advanced Trade WebSocket Docs](https://docs.cloud.coinbase.com/advanced-trade-api/docs/ws-overview)
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

```go
sub := coinbasev3.NewWsFeedSubscription(
    coinbasev3.SubTypeSubscribe,
    []string{"ETH-USD", "BTC-USD"},
    coinbasev3.ChannelTypeTicker,
    "api_key",
    "secret_key",
)

readCh := make(chan []byte)
ws, err := coinbasev3.NewWsClient(coinbase.WsClientConfig{
    Url:              "wss://advanced-trade-ws.coinbase.com",
    ReadChannel:      readCh,
    SubscriptionData: sub.Marshal(),
    OnConnect:        func() {},
    OnDisconnect:     func() {},
    OnReconnect:      func() {},
    UseBackoff:       true,
    Debug:            true,
})
if err != nil {
    panic(err)
}

_, err = ws.Connect()
if err != nil {
    panic(err)
}

for {
    select {
        case msg := <-readCh:
            log.Println(string(msg))
    }
}
```

## Run tests

```bash
go test ./... -v
```
