package coinbasev3

import (
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
	ChannelTypeLevel2       ChannelType = "l2_data"
	ChannelTypeUser         ChannelType = "user"
	ChannelTypeMarketTrades ChannelType = "market_trades"
)

type SubType string

const (
	SubTypeSubscribe   SubType = "subscribe"
	SubTypeUnsubscribe SubType = "unsubscribe"
)

type WebsocketChannel struct {
	Type       SubType     `json:"type"`
	ProductIds []string    `json:"product_ids"`
	Channel    ChannelType `json:"channel"`
	Signature  string      `json:"signature"`
	ApiKey     string      `json:"api_key"`
	SecretKey  string      `json:"-"`
	Timestamp  string      `json:"timestamp"`
}

func NewWebsocketChannel(subType SubType, channel ChannelType, productIds []string) WebsocketChannel {
	s := WebsocketChannel{
		Type:       subType,
		ProductIds: productIds,
		Channel:    channel,
	}
	return s
}

func NewChannelSubscribe(channel ChannelType, productIds []string) WebsocketChannel {
	return NewWebsocketChannel(SubTypeSubscribe, channel, productIds)
}

func NewChannelUnsubscribe(channel ChannelType, productIds []string) WebsocketChannel {
	return NewWebsocketChannel(SubTypeUnsubscribe, channel, productIds)
}

func NewTickerChannel(productIds []string) WebsocketChannel {
	return NewChannelSubscribe(ChannelTypeTicker, productIds)
}

func NewTickerBatchChannel(productIds []string) WebsocketChannel {
	return NewChannelSubscribe(ChannelTypeTickerBatch, productIds)
}

func NewCandlesChannel(productIds []string) WebsocketChannel {
	return NewChannelSubscribe(ChannelTypeCandles, productIds)
}

func NewHeartbeatsChannel(productIds []string) WebsocketChannel {
	return NewChannelSubscribe(ChannelTypeHeartbeats, productIds)
}

func NewStatusChannel(productIds []string) WebsocketChannel {
	return NewChannelSubscribe(ChannelTypeStatus, productIds)
}

func NewLevel2Channel(productIds []string) WebsocketChannel {
	return NewChannelSubscribe(ChannelTypeLevel2, productIds)
}

func NewUserChannel(productIds []string) WebsocketChannel {
	return NewChannelSubscribe(ChannelTypeUser, productIds)
}

func (s *WebsocketChannel) marshal(apiKey, secretKey string) []byte {
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

func (s *WebsocketChannel) setSignature() {
	// Concatenating and comma-separating the timestamp, channel name, and product Ids, for example: 1660838876level2ETH-USD,ETH-EUR.
	sig := fmt.Sprintf("%s%s%s", s.Timestamp, s.Channel, strings.Join(s.ProductIds, ","))
	s.Signature = string(SignHmacSha256(sig, s.SecretKey))
}

func (s *WebsocketChannel) setTimestamp() {
	s.Timestamp = fmt.Sprintf("%d", int(time.Now().Unix()))
}
