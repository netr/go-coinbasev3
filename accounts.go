package coinbasev3

import (
	"fmt"
	"time"
)

type Account struct {
	Uuid             string                  `json:"uuid"`
	Name             string                  `json:"name"`
	Currency         string                  `json:"currency"`
	AvailableBalance AccountAvailableBalance `json:"available_balance"`
	Default          bool                    `json:"default"`
	Active           bool                    `json:"active"`
	CreatedAt        time.Time               `json:"created_at"`
	UpdatedAt        time.Time               `json:"updated_at"`
	DeletedAt        time.Time               `json:"deleted_at"`
	Type             string                  `json:"type"`
	Ready            bool                    `json:"ready"`
	Hold             AccountHold             `json:"hold"`
}

type AccountAvailableBalance struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type AccountHold struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

// ListAccounts gets a list of authenticated accounts for the current user.
func (c *ApiClient) ListAccounts(limit int, cursor string) (ListAccountsData, error) {
	// A pagination limit with default of 49 and maximum of 250.
	if limit < 49 {
		limit = 49
	}
	if limit > 250 {
		limit = 250
	}

	u := fmt.Sprintf("https://api.coinbase.com/api/v3/brokerage/accounts?limit=%d&cursor=%s", limit, cursor)

	var data ListAccountsData
	resp, err := c.client.R().SetSuccessResult(&data).Get(u)
	if err != nil {
		return data, err
	}

	if !resp.IsSuccessState() {
		return data, ErrFailedToUnmarshal
	}

	return data, nil
}

type ListAccountsData struct {
	Accounts []Account `json:"accounts"`
	HasNext  bool      `json:"has_next"`
	Cursor   string    `json:"cursor"`
	Size     int       `json:"size"`
}

// GetAccount get a list of information about an account, given an account UUID.
func (c *ApiClient) GetAccount(uuid string) (Account, error) {
	u := fmt.Sprintf("https://api.coinbase.com/api/v3/brokerage/accounts/%s", uuid)

	var data GetAccountData
	resp, err := c.client.R().SetSuccessResult(&data).Get(u)
	if err != nil {
		return data.Account, err
	}

	if !resp.IsSuccessState() {
		return data.Account, ErrFailedToUnmarshal
	}

	return data.Account, nil
}

type GetAccountData struct {
	Account Account `json:"account"`
}
