package paypal

import (
	"fmt"
	"net/http"
)

type BalanceAccounts struct {
	TotalAvailable string `json:"total_available"`
}

// GetBalanceAccounts - Get the balance accounts for the current user
func (c *client) GetBalanceAccounts() (*BalanceAccounts, error) {
	out := &BalanceAccounts{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.APIBase, "/v2/wallet/balance-accounts"), nil)
	if err != nil {
		return out, err
	}

	if err := c.sendWithAuth(req, out); err != nil {
		return out, err
	}
	return out, nil
}
