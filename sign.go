package coinbasev3

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func SignHmacSha256(str, secret string) []byte {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(str))
	sha := h.Sum(nil)
	return []byte(hex.EncodeToString(sha))
}
