package coinbasev3

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type ChannelType string

const (
	ChannelTypeHeartbeats   ChannelType = "heartbeats"
	ChannelTypeCandles      ChannelType = "candles"
	ChannelTypeStatus       ChannelType = "status"
	ChannelTypeTicker       ChannelType = "ticker"
	ChannelTypeTickerBatch  ChannelType = "ticker_batch"
	ChannelTypeLevel2       ChannelType = "level2"
	ChannelTypeUser         ChannelType = "user"
	ChannelTypeMarketTrades ChannelType = "market_trades"
)

type SubType string

const (
	SubTypeSubscribe   SubType = "subscribe"
	SubTypeUnsubscribe SubType = "unsubscribe"
)

type WsChannel struct {
	Type       SubType     `json:"type"`
	ProductIds []string    `json:"product_ids"`
	Channel    ChannelType `json:"channel"`
	Signature  string      `json:"signature"`
	ApiKey     string      `json:"api_key"`
	SecretKey  string      `json:"-"`
	Timestamp  string      `json:"timestamp"`
}

func NewWsChannel(subType SubType, productIds []string, channel ChannelType) WsChannel {
	s := WsChannel{
		Type:       subType,
		ProductIds: productIds,
		Channel:    channel,
	}
	return s
}

func NewWsChannelSub(productIds []string, channel ChannelType) WsChannel {
	return NewWsChannel(SubTypeSubscribe, productIds, channel)
}

func NewWsChannelUnsub(productIds []string, channel ChannelType) WsChannel {
	return NewWsChannel(SubTypeUnsubscribe, productIds, channel)
}

func (s *WsChannel) marshal(apiKey, secretKey string) []byte {
	s.ApiKey = apiKey
	s.SecretKey = secretKey

	s.setTimestamp()
	s.setSignature()

	b, err := json.Marshal(s)
	if err != nil {
		return nil
	}

	return b
}

func (s *WsChannel) setSignature() {
	// Concatenating and comma-separating the timestamp, channel name, and product Ids, for example: 1660838876level2ETH-USD,ETH-EUR.
	sig := fmt.Sprintf("%s%s%s", s.Timestamp, s.Channel, strings.Join(s.ProductIds, ","))
	s.Signature = string(sign(sig, s.SecretKey))
}

func (s *WsChannel) setTimestamp() {
	s.Timestamp = fmt.Sprintf("%d", int(time.Now().Unix()))
}

func sign(str, secret string) []byte {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(str))
	sha := h.Sum(nil)
	return []byte(hex.EncodeToString(sha))
}
