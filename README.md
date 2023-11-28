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

```go
tick := coinbasev3.NewWsChannelSub(
  []string{"ETH-USD", "BTC-USD"},
  coinbasev3.ChannelTypeTicker,
)
hb := coinbasev3.NewWsChannelSub(
  []string{"ETH-USD", "ETH-EUR"},
  coinbasev3.ChannelTypeHeartbeats,
)

readCh := make(chan []byte)
ws, err := coinbasev3.NewWsClient(coinbasev3.WsClientConfig{
  ApiKey:      config.Coinbase.ApiKey,
  SecretKey:   config.Coinbase.SecretKey,
  ReadChannel: readCh,
  WsChannels:  []coinbasev3.WsChannel{tick, hb},
  OnConnect: func() {
    log.Println("Connected to Coinbase")
  },
  OnDisconnect: func() {
    log.Println("Disconnected from Coinbase")
  },
  OnReconnect: func() {
    log.Println("Reconnected to Coinbase")
  },
  UseBackoff: true,
  Debug:      true,
})
if err != nil {
  panic("Failed to create Coinbase websocket client")
}

_, err = ws.Connect()
if err != nil {
    panic("Failed to connect to Coinbase websocket")
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
