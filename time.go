package go_coinbasev3

import "time"

// GetServerTime get the API server time.
func (c *ApiClient) GetServerTime() (ServerTime, error) {
	u := "https://api.coinbase.com/v2/time"

	var servTime ServerTime
	resp, err := c.client.R().SetSuccessResult(&servTime).Get(u)
	if err != nil {
		return servTime, err
	}

	if !resp.IsSuccessState() {
		return servTime, ErrFailedToUnmarshal
	}

	return servTime, nil
}

type ServerTime struct {
	Data struct {
		Iso   time.Time `json:"iso"`
		Epoch int       `json:"epoch"`
	} `json:"data"`
}
