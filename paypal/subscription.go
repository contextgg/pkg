package paypal

import (
	"fmt"
	"time"
)

// CreateSubscriptionPlan creates a new plan for users to subscribe to
func (c *client) CreateSubscriptionPlan(plan *SubscriptionPlanInput) (*SubscriptionPlanResp, error) {
	req, err := c.newRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/plans"), plan)
	response := &SubscriptionPlanResp{}
	if err != nil {
		return response, err
	}
	err = c.sendWithAuth(req, response)
	return response, err
}

// DeactivateSubscriptionPlan deactivates an existing subscription plan
func (c *client) DeactivateSubscriptionPlan(planID string) (*SubscriptionPlanResp, error) {
	req, err := c.newRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/plans/"+planID+"/deactivate"), nil)
	response := &SubscriptionPlanResp{}
	if err != nil {
		return response, err
	}
	err = c.sendWithAuth(req, response)
	return response, err
}

// UpdateSubscriptionPlanPricing deactivates an existing subscription plan
func (c *client) UpdateSubscriptionPlanPricing(planID string, scheme *UpdatePricingSchemeInput) (*SubscriptionPlanResp, error) {
	req, err := c.newRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/plans/"+planID+"/update-pricing-schemes"), scheme)
	response := &SubscriptionPlanResp{}
	if err != nil {
		return response, err
	}
	err = c.sendWithAuth(req, response)
	return response, err
}

// ActivateSubscriptionPlan activates the specified subscription plan so it can be subscribed to
func (c *client) ActivateSubscriptionPlan(planID string) (*SubscriptionPlanResp, error) {
	req, err := c.newRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/plans/"+planID+"/activate"), nil)
	response := &SubscriptionPlanResp{}
	if err != nil {
		return response, err
	}
	err = c.sendWithAuth(req, response)
	return response, err
}

// GetSubscriptionPlan gets a specific subscription plan
func (c *client) GetSubscriptionPlan(planID string) (*SubscriptionPlanResp, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/plans/"+planID), nil)
	response := &SubscriptionPlanResp{}
	if err != nil {
		return response, err
	}
	err = c.sendWithAuth(req, response)
	return response, err
}

// ListSubscriptionPlans list all subscription plans
func (c *client) ListSubscriptionPlans() ([]ListSubscriptionPlan, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/plans"), nil)
	res := &ListSubscriptionPlansResp{}
	if err != nil {
		return res.Plans, err
	}
	err = c.sendWithAuth(req, res)
	return res.Plans, err
}

// NewSubscription list all subscription plans
func (c *client) NewSubscription(sub *SubscriptionInput) (*SubscriptionResp, error) {
	req, err := c.newRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/subscriptions"), sub)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	response := &SubscriptionResp{}
	if err != nil {
		return response, err
	}
	err = c.sendWithAuth(req, response)
	return response, err
}

// GetSubscription list all subscription plans
func (c *client) GetSubscription(subId string) (*SubscriptionResp, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/subscriptions/"+subId), nil)
	response := &SubscriptionResp{}
	if err != nil {
		return response, err
	}
	err = c.sendWithAuth(req, response)
	return response, err
}

// ActivateSubscription activates the specified subscription
func (c *client) ActivateSubscription(subId string) error {
	req, err := c.newRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/subscriptions/"+subId+"/activate"), nil)
	if err != nil {
		return err
	}
	err = c.sendWithAuth(req, nil)
	return err
}

// CancelSubscription cancel the specified subscription
func (c *client) CancelSubscription(subId string) error {
	req, err := c.newRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/subscriptions/"+subId+"/cancel"), nil)
	if err != nil {
		return err
	}
	err = c.sendWithAuth(req, nil)
	return err
}

func (c *client) SubscriptionListTransactions(subId string, startTime, endTime time.Time) (*SubscriptionTransactionsResp, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/subscriptions/"+subId+"/transactions"), nil)
	q := req.URL.Query()
	q.Add("start_time", startTime.UTC().Format(time.RFC3339))
	q.Add("end_time", endTime.UTC().Format(time.RFC3339))
	req.URL.RawQuery = q.Encode()
	response := &SubscriptionTransactionsResp{}
	if err != nil {
		return response, err
	}
	err = c.sendWithAuth(req, response)
	return response, err
}
