package paypal

import "fmt"

// GetSale returns a sale by ID
// Use this call to get details about a sale transaction.
// Note: This call returns only the sales that were created via the REST API.
// Endpoint: GET /v2/payments/captures/ID
func (c *client) GetCapture(captureId string) (*Capture, error) {
	capture := &Capture{}

	req, err := c.newRequest("GET", fmt.Sprintf("%s%s", c.APIBase, "/v2/payments/captures/"+captureId), nil)
	if err != nil {
		return capture, err
	}

	if err = c.sendWithAuth(req, capture); err != nil {
		return capture, err
	}

	return capture, nil
}

// RefundSale refunds a completed payment.
// Use this call to refund a completed payment. Provide the sale_id in the URI and an empty JSON payload for a full refund. For partial refunds, you can include an amount.
// Endpoint: POST /v2/payments/sale/ID/refund
func (c *client) RefundCapture(captureId string, a *Amount) (*Refund, error) {
	type refundRequest struct {
		Amount *Amount `json:"amount"`
	}

	refund := &Refund{}

	req, err := c.newRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v2/payments/captures/"+captureId+"/refund"), &refundRequest{Amount: a})
	if err != nil {
		return refund, err
	}

	if err = c.sendWithAuth(req, refund); err != nil {
		return refund, err
	}

	return refund, nil
}

// GetRefund by ID
// Use it to look up details of a specific refund on direct and captured payments.
// Endpoint: GET /v2/payments/refund/ID
func (c *client) GetRefund(refundID string) (*Refund, error) {
	refund := &Refund{}

	req, err := c.newRequest("GET", fmt.Sprintf("%s%s", c.APIBase, "/v2/payments/refund/"+refundID), nil)
	if err != nil {
		return refund, err
	}

	if err = c.sendWithAuth(req, refund); err != nil {
		return refund, err
	}

	return refund, nil
}
