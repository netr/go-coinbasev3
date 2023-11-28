package coinbasev3

import "github.com/imroc/req/v3"

type MockHttpClient struct {
	Response *req.Response
	Err      error
	client   *req.Client
}

func (m *MockHttpClient) Get(url string) (*req.Response, error) {
	return m.Response, m.Err
}

func (m *MockHttpClient) GetClient() *req.Client {
	return nil
}

func NewMockHttpClient(resp *req.Response) HttpClient {
	return &MockHttpClient{
		Response: resp,
	}
}
