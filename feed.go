package go_coinbasev3

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type WsFeedSubscription struct {
	Type       string   `json:"type"`
	ProductIds []string `json:"product_ids"`
	Channel    string   `json:"channel"`
	Signature  string   `json:"signature"`
	ApiKey     string   `json:"api_key"`
	SecretKey  string   `json:"-"`
	Timestamp  int      `json:"timestamp"`
}

func NewWsFeedSubscription(subType string, productIds []string, channel, apiKey, secretKey string) *WsFeedSubscription {
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