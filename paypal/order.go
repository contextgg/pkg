package paypal

import "fmt"

// GetOrder retrieves order by ID
// Endpoint: GET /v2/checkout/orders/ID
func (c *client) GetOrder(orderID string) (*Order, error) {
	order := &Order{}

	req, err := c.newRequest("GET", fmt.Sprintf("%s%s%s", c.APIBase, "/v2/checkout/orders/", orderID), nil)
	if err != nil {
		return order, err
	}

	if err = c.sendWithAuth(req, order); err != nil {
		return order, err
	}

	return order, nil
}

// CreateOrder - Use this call to create an order
// Endpoint: POST /v2/checkout/orders
func (c *client) CreateOrder(intent string, purchaseUnits []PurchaseUnitRequest, payer *CreateOrderPayer, appContext *ApplicationContext) (*Order, error) {
	type createOrderRequest struct {
		Intent             string                `json:"intent"`
		Payer              *CreateOrderPayer     `json:"payer,omitempty"`
		PurchaseUnits      []PurchaseUnitRequest `json:"purchase_units"`
		ApplicationContext *ApplicationContext   `json:"application_context,omitempty"`
	}

	order := &Order{}

	req, err := c.newRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v2/checkout/orders"), createOrderRequest{Intent: intent, PurchaseUnits: purchaseUnits, Payer: payer, ApplicationContext: appContext})
	if err != nil {
		return order, err
	}

	if err = c.sendWithAuth(req, order); err != nil {
		return order, err
	}

	return order, nil
}

// UpdateOrder updates the order by ID
// Endpoint: PATCH /v2/checkout/orders/ID
func (c *client) UpdateOrder(orderID string, purchaseUnits []PurchaseUnitRequest) (*Order, error) {
	order := &Order{}

	req, err := c.newRequest("PATCH", fmt.Sprintf("%s%s%s", c.APIBase, "/v2/checkout/orders/", orderID), purchaseUnits)
	if err != nil {
		return order, err
	}

	if err = c.sendWithAuth(req, order); err != nil {
		return order, err
	}

	return order, nil
}

// AuthorizeOrder - https://developer.paypal.com/docs/api/orders/v2/#orders_authorize
// Endpoint: POST /v2/checkout/orders/ID/authorize
func (c *client) AuthorizeOrder(orderID string, authorizeOrderRequest AuthorizeOrderRequest) (*Authorization, error) {
	auth := &Authorization{}

	req, err := c.newRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v2/checkout/orders/"+orderID+"/authorize"), authorizeOrderRequest)
	if err != nil {
		return auth, err
	}

	if err = c.sendWithAuth(req, auth); err != nil {
		return auth, err
	}

	return auth, nil
}

// CaptureOrder - https://developer.paypal.com/docs/api/orders/v2/#orders_capture
// Endpoint: POST /v2/checkout/orders/ID/capture
func (c *client) CaptureOrder(orderID string, captureOrderRequest CaptureOrderRequest) (*CaptureOrderResponse, error) {
	capture := &CaptureOrderResponse{}

	req, err := c.newRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v2/checkout/orders/"+orderID+"/capture"), captureOrderRequest)
	if err != nil {
		return capture, err
	}

	if err = c.sendWithAuth(req, capture); err != nil {
		return capture, err
	}

	return capture, nil
}
