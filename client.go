package coinbasev3

import (
	"fmt"
	"github.com/imroc/req/v3"
	"time"
)

type ApiClient struct {
	apiKey    string
	secretKey string
	client    *req.Client
}

var (
	ErrFailedToUnmarshal = fmt.Errorf("failed to unmarshal response")
)

func NewApiClient(apiKey, secretKey string) *ApiClient {
	return &ApiClient{
		apiKey:    apiKey,
		secretKey: secretKey,
		client: req.C().
			SetTimeout(time.Second*10).
			SetUserAgent("GoCoinbaseV3/1.0.0").
			SetCommonRetryCount(3).
			SetCommonRetryBackoffInterval(1*time.Second, 5*time.Second),
	}
}
