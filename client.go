package coinbasev3

import (
	"fmt"
	"github.com/imroc/req/v3"
	"net/url"
	"time"
)

type HttpClient interface {
	Get(url string) (*req.Response, error)
	GetClient() *req.Client
}

type ApiClient struct {
	apiKey     string
	secretKey  string
	client     *req.Client
	httpClient HttpClient
}

func (c *ApiClient) GetClient() *req.Client {
	return c.client
}

type ReqClient struct {
	client *req.Client
}

func (c *ReqClient) GetClient() *req.Client {
	return c.client
}

var (
	ErrFailedToUnmarshal = fmt.Errorf("failed to unmarshal response")
)

// NewApiClient creates a new Coinbase API client. The API key and secret key are used to sign requests. The default timeout is 10 seconds. The default retry count is 3. The default retry backoff interval is 1 second to 5 seconds.
func NewApiClient(apiKey, secretKey string, clients ...HttpClient) *ApiClient {
	if clients != nil && len(clients) > 0 {
		return &ApiClient{
			apiKey:     apiKey,
			secretKey:  secretKey,
			client:     newClient(apiKey, secretKey),
			httpClient: clients[0],
		}
	}

	client := newClient(apiKey, secretKey)

	return &ApiClient{
		apiKey:     apiKey,
		secretKey:  secretKey,
		client:     client,
		httpClient: &ReqClient{client: client},
	}
}

func (c *ApiClient) get(url string, out interface{}) error {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return err
	}

	if !resp.IsSuccessState() {
		return ErrFailedToUnmarshal
	}

	err = resp.Unmarshal(&out)
	if err != nil {
		return err
	}
	return nil
}

func newClient(apiKey, secretKey string) *req.Client {
	client := req.C().
		SetTimeout(time.Second * 10).
		SetUserAgent("GoCoinbaseV3/1.0.0")

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

	return client
}

func (c *ReqClient) Get(url string) (*req.Response, error) {
	resp, err := c.client.R().Get(url)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
