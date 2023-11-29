package coinbasev3

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/imroc/req/v3"
	"net/url"
	"strings"
	"time"
)

var (
	ErrFailedToUnmarshal = fmt.Errorf("failed to unmarshal response")
)

type HttpClient interface {
	Get(url string) (*req.Response, error)
	Post(url string, data []byte) (*req.Response, error)
	GetClient() *req.Client
}

type ApiClient struct {
	apiKey          string
	secretKey       string
	client          *req.Client
	httpClient      HttpClient
	baseUrlV3       string
	baseUrlV2       string
	baseExchangeUrl string
}

// NewApiClient creates a new Coinbase API client. The API key and secret key are used to sign requests. The default timeout is 10 seconds. The default retry count is 3. The default retry backoff interval is 1 second to 5 seconds.
func NewApiClient(apiKey, secretKey string, clients ...HttpClient) *ApiClient {
	if clients != nil && len(clients) > 0 {
		ac := &ApiClient{
			apiKey:     apiKey,
			secretKey:  secretKey,
			client:     newClient(apiKey, secretKey),
			httpClient: clients[0],
		}
		ac.setBaseUrls()
		return ac
	}

	client := newClient(apiKey, secretKey)

	ac := &ApiClient{
		apiKey:     apiKey,
		secretKey:  secretKey,
		client:     client,
		httpClient: &ReqClient{client: client},
	}
	ac.setBaseUrls()
	return ac
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

func (c *ApiClient) get(url string, out interface{}) ([]byte, error) {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return resp.Bytes(), err
	}

	if !resp.IsSuccessState() {
		return resp.Bytes(), ErrFailedToUnmarshal
	}

	err = resp.Unmarshal(&out)
	if err != nil {
		return resp.Bytes(), err
	}

	return resp.Bytes(), nil
}

func (c *ApiClient) post(url string, data []byte, out interface{}) ([]byte, error) {
	resp, err := c.httpClient.Post(url, data)
	if err != nil {
		return resp.Bytes(), err
	}

	if !resp.IsSuccessState() {
		return resp.Bytes(), ErrFailedToUnmarshal
	}

	err = resp.Unmarshal(&out)
	if err != nil {
		return resp.Bytes(), err
	}
	return resp.Bytes(), nil
}

func (c *ApiClient) setBaseUrls() {
	c.baseUrlV3 = "https://api.coinbase.com/api/v3"
	c.baseUrlV2 = "https://api.coinbase.com/api/v2"
	c.baseExchangeUrl = "https://api.exchange.coinbase.com"
}

// SetSandboxUrls sets the base URLs to the sandbox environment. Note: The sandbox for Advanced Trading is not yet available. This method will be revisited when the sandbox is available.
func (c *ApiClient) SetSandboxUrls() {
	c.baseUrlV3 = "https://api-public.sandbox.pro.coinbase.com"
	c.baseUrlV2 = "https://api-public.sandbox.pro.coinbase.com"
	c.baseExchangeUrl = "https://api-public.sandbox.exchange.coinbase.com"
}

// SetBaseUrlV3 sets the base URL for the Coinbase Advanced Trading API.
func (c *ApiClient) SetBaseUrlV3(url string) {
	c.baseUrlV3 = url
}

func (c *ApiClient) makeV3Url(path string) string {
	if strings.HasPrefix(path, "/") {
		path = strings.TrimPrefix(path, "/")
	}
	return fmt.Sprintf("%s/%s", c.baseUrlV3, path)
}

// SetBaseUrlV2 sets the base URL for the Sign In With Coinbase APIs.
func (c *ApiClient) SetBaseUrlV2(url string) {
	c.baseUrlV2 = url
}

func (c *ApiClient) makeV2Url(path string) string {
	if strings.HasPrefix(path, "/") {
		path = strings.TrimPrefix(path, "/")
	}
	return fmt.Sprintf("%s/%s", c.baseUrlV2, path)
}

// SetBaseExchangeUrl sets the base URL for the Coinbase Exchange API.
func (c *ApiClient) SetBaseExchangeUrl(url string) {
	c.baseExchangeUrl = url
}

func (c *ApiClient) makeExchangeUrl(path string) string {
	if strings.HasPrefix(path, "/") {
		path = strings.TrimPrefix(path, "/")
	}
	return fmt.Sprintf("%s/%s", c.baseExchangeUrl, path)
}

// ReqClient is a wrapper around the req.Client to satisfy the HttpClient interface.
type ReqClient struct {
	client *req.Client
}

// GetClient returns the underlying req.Client.
func (c *ReqClient) GetClient() *req.Client {
	return c.client
}

// Get makes a GET request to the given URL.
func (c *ReqClient) Get(url string) (*req.Response, error) {
	resp, err := c.client.R().Get(url)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Post makes a POST request to the given URL.
func (c *ReqClient) Post(url string, data []byte) (*req.Response, error) {
	resp, err := c.client.R().SetBody(data).Post(url)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type ErrorResponse struct {
	Error        string       `json:"error"`
	Code         string       `json:"code"`
	Message      string       `json:"message"`
	ErrorDetails string       `json:"error_details"`
	Details      ErrorDetails `json:"details"`
}

func newErrorResponse(res []byte) ErrorResponse {
	var errRes ErrorResponse
	err := json.Unmarshal(res, &errRes)
	if err != nil {
		return ErrorResponse{
			Error:   err.Error(),
			Code:    "unknown",
			Message: err.Error(),
		}
	}
	return errRes
}

type ErrorDetail struct {
	TypeUrl string `json:"type_url"`
	Value   string `json:"value"`
}

type ErrorDetails []ErrorDetail

// UnmarshalJSON implements the json.Unmarshaler interface. Required because Coinbase returns an array of error details or a single error detail object.
func (ed *ErrorDetails) UnmarshalJSON(data []byte) error {
	var details []ErrorDetail
	if err := json.Unmarshal(data, &details); err == nil {
		*ed = details
		return nil
	}
	var detail ErrorDetail
	if err := json.Unmarshal(data, &detail); err == nil {
		*ed = ErrorDetails{detail}
		return nil
	}
	return errors.New("error details should be an array or a single object")
}
