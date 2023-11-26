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
	ChannelTypeHeartbeats   ChannelType = "heartbeat"
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

type WsFeedSubscription struct {
	Type       SubType     `json:"type"`
	ProductIds []string    `json:"product_ids"`
	Channel    ChannelType `json:"channel"`
	Signature  string      `json:"signature"`
	ApiKey     string      `json:"api_key"`
	SecretKey  string      `json:"-"`
	Timestamp  int         `json:"timestamp"`
}

func NewWsFeedSubscription(subType SubType, productIds []string, channel ChannelType, apiKey, secretKey string) *WsFeedSubscription {
	s := &WsFeedSubscription{
		Type:       subType,
		ProductIds: productIds,
		Channel:    channel,
		ApiKey:     apiKey,
		SecretKey:  secretKey,
	}
	s.setTimestamp()
	s.setSignature()

	return s
}

func (s *WsFeedSubscription) Marshal() []byte {
	b, err := json.Marshal(s)
	if err != nil {
		return nil
	}

	return b
}

func (s *WsFeedSubscription) setSignature() {
	// Concatenating and comma-separating the timestamp, channel name, and product Ids, for example: 1660838876level2ETH-USD,ETH-EUR.
	sig := fmt.Sprintf("%d%s%s", s.Timestamp, s.Channel, strings.Join(s.ProductIds, ","))
	s.Signature = string(sign(sig, s.SecretKey))
}

func (s *WsFeedSubscription) setTimestamp() {
	s.Timestamp = int(time.Now().Unix())
	s.setSignature()
}

func sign(str, secret string) []byte {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(str))
	sha := h.Sum(nil)
	return []byte(hex.EncodeToString(sha))
}
