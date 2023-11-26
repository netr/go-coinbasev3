# Coinbase Advanced Trades API (v3)

## Run tests

```bash
go test ./... -v
```

## Websocket

```go
sub := coinbase.NewWsFeedSubscription(
    coinbase.SubTypeSubscribe,
    []string{"ETH-USD", "BTC-USD"},
    coinbase.ChannelTypeTicker,
    "api_key",
    "secret_key",
)

readCh := make(chan []byte)
ws, err := coinbase.NewWsClient(coinbase.WsClientConfig{
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