package coinbasev3

import (
	"fmt"
	"github.com/imroc/req/v3"
	"net/url"
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

// NewApiClient creates a new Coinbase API client. The API key and secret key are used to sign requests. The default timeout is 10 seconds. The default retry count is 3. The default retry backoff interval is 1 second to 5 seconds.
func NewApiClient(apiKey, secretKey string) *ApiClient {
	client := req.C().
		SetTimeout(time.Second*10).
		SetUserAgent("GoCoinbaseV3/1.0.0").
		SetCommonRetryCount(3).
		SetCommonRetryBackoffInterval(1*time.Second, 5*time.Second)

	// TODO: figure out how to do this where we can use PathParam, QueryParam, etc.
	client.OnBeforeRequest(func(client *req.Client, req *req.Request) error {
		// create a secret key from: `timestamp + method + requestPath + body`
		path := ""
		if req.RawURL != "" {
			u, err := url.Parse(req.RawURL)
			if err != nil {
				return err
			}
			path = u.Path
		} else {
			return fmt.Errorf("no path found")
		}

		sig := fmt.Sprintf("%d%s%s%s", time.Now().Unix(), req.Method, path, req.Body)
		signedSig := string(SignHmacSha256(sig, secretKey))

		client.Headers.Set("CB-ACCESS-KEY", apiKey)
		client.Headers.Set("CB-ACCESS-SIGN", signedSig)
		client.Headers.Set("CB-ACCESS-TIMESTAMP", fmt.Sprintf("%d", time.Now().Unix()))
		return nil
	})

	return &ApiClient{
		apiKey:    apiKey,
		secretKey: secretKey,
		client:    client,
	}

}
