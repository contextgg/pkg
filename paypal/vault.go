package paypal

import (
	"fmt"
)

// StoreCreditCard func
// Endpoint: POST /v1/vault/credit-cards
func (c *client) StoreCreditCard(cc CreditCard) (*CreditCard, error) {
	req, err := c.newRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/vault/credit-cards"), cc)
	if err != nil {
		return nil, err
	}

	response := &CreditCard{}

	if err = c.sendWithAuth(req, response); err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteCreditCard func
// Endpoint: DELETE /v1/vault/credit-cards/credit_card_id
func (c *client) DeleteCreditCard(id string) error {
	req, err := c.newRequest("DELETE", fmt.Sprintf("%s/v1/vault/credit-cards/%s", c.APIBase, id), nil)
	if err != nil {
		return err
	}

	if err = c.sendWithAuth(req, nil); err != nil {
		return err
	}

	return nil
}

// GetCreditCard func
// Endpoint: GET /v1/vault/credit-cards/credit_card_id
func (c *client) GetCreditCard(id string) (*CreditCard, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("%s/v1/vault/credit-cards/%s", c.APIBase, id), nil)
	if err != nil {
		return nil, err
	}

	response := &CreditCard{}

	if err = c.sendWithAuth(req, response); err != nil {
		return nil, err
	}

	return response, nil
}

// GetCreditCards func
// Endpoint: GET /v1/vault/credit-cards
func (c *client) GetCreditCards(ccf *CreditCardsFilter) (*CreditCards, error) {
	page := 1
	if ccf != nil && ccf.Page > 0 {
		page = ccf.Page
	}
	pageSize := 10
	if ccf != nil && ccf.PageSize > 0 {
		pageSize = ccf.PageSize
	}

	req, err := c.newRequest("GET", fmt.Sprintf("%s/v1/vault/credit-cards?page=%d&page_size=%d", c.APIBase, page, pageSize), nil)
	if err != nil {
		return nil, err
	}

	response := &CreditCards{}

	if err = c.sendWithAuth(req, response); err != nil {
		return nil, err
	}

	return response, nil
}

// PatchCreditCard func
// Endpoint: PATCH /v1/vault/credit-cards/credit_card_id
func (c *client) PatchCreditCard(id string, ccf []CreditCardField) (*CreditCard, error) {
	req, err := c.newRequest("PATCH", fmt.Sprintf("%s/v1/vault/credit-cards/%s", c.APIBase, id), ccf)
	if err != nil {
		return nil, err
	}

	response := &CreditCard{}

	if err = c.sendWithAuth(req, response); err != nil {
		return nil, err
	}

	return response, nil
}
