package paypal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"sync"
	"time"
)

type Client interface {
	SetLog(log io.Writer)
	IsSandbox() bool

	GetBalanceAccounts() (*BalanceAccounts, error)

	GetCapture(captureId string) (*Capture, error)

	GetOrder(orderID string) (*Order, error)
	CreateOrder(intent string, purchaseUnits []PurchaseUnitRequest, payer *CreateOrderPayer, appContext *ApplicationContext) (*Order, error)
	CaptureOrder(orderID string, captureOrderRequest CaptureOrderRequest) (*CaptureOrderResponse, error)

	CreateSubscriptionPlan(plan *SubscriptionPlanInput) (*SubscriptionPlanResp, error)

	NewSubscription(sub *SubscriptionInput) (*SubscriptionResp, error)
	GetSubscription(subId string) (*SubscriptionResp, error)
	ActivateSubscription(subId string) error
	CancelSubscription(subId string) error
	SubscriptionListTransactions(subId string, startTime, endTime time.Time) (*SubscriptionTransactionsResp, error)

	CreateSinglePayout(p Payout) (*PayoutResponse, error)

	Transactions(query *TransactionQuery) (*TransactionsResponse, error)
}

// NewClient returns new Client struct
// APIBase is a base API URL, for testing you can use paypal.APIBaseSandBox
func NewClient(clientID string, secret string, APIBase string) (Client, error) {
	if clientID == "" || secret == "" || APIBase == "" {
		return nil, errors.New("ClientID, Secret and APIBase are required to create a Client")
	}

	return &client{
		Client:   &http.Client{},
		ClientID: clientID,
		Secret:   secret,
		APIBase:  APIBase,
	}, nil
}

// Client represents a Paypal REST API Client
type client struct {
	sync.Mutex
	Client         *http.Client
	ClientID       string
	Secret         string
	APIBase        string
	Log            io.Writer // If user set log file name all requests will be logged there
	Token          *TokenResponse
	tokenExpiresAt time.Time
}

// SetLog will set/change the output destination.
// If log file is set paypal will log all requests and responses to this Writer
func (c *client) SetLog(log io.Writer) {
	c.Log = log
}

// IsSandbox checks the current APIBase against the sandbox url
func (c *client) IsSandbox() bool {
	return c.APIBase == APIBaseSandBox
}

// Send makes a request to the API, the response body will be
// unmarshaled into v, or if v is an io.Writer, the response will
// be written to it without decoding
func (c *client) send(req *http.Request, v interface{}) error {
	var (
		err  error
		resp *http.Response
		data []byte
	)

	// Set default headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en_US")
	req.Header.Set("Prefer", "return=representation")

	// Default values for headers
	if req.Header.Get("Content-type") == "" {
		req.Header.Set("Content-type", "application/json")
	}

	resp, err = c.Client.Do(req)
	c.log(req, resp)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		errResp := &ErrorResponse{Response: resp}
		data, err = ioutil.ReadAll(resp.Body)

		if err == nil && len(data) > 0 {
			json.Unmarshal(data, errResp)
		}

		return errResp
	}
	if v == nil {
		return nil
	}

	if w, ok := v.(io.Writer); ok {
		io.Copy(w, resp.Body)
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(v)
}

// GetAccessToken returns struct of TokenResponse
// No need to call SetAccessToken to apply new access token for current Client
// Endpoint: POST /v1/oauth2/token
func (c *client) getAccessToken() (*TokenResponse, error) {
	buf := bytes.NewBuffer([]byte("grant_type=client_credentials"))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/oauth2/token"), buf)
	if err != nil {
		return &TokenResponse{}, err
	}

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")

	response := &TokenResponse{}
	err = c.sendWithBasicAuth(req, response)

	// Set Token fur current Client
	if response.Token != "" {
		c.Token = response
		c.tokenExpiresAt = time.Now().Add(time.Duration(response.ExpiresIn) * time.Second)
	}

	return response, err
}

// ensureToken will fetch the token if needed or expired
func (c *client) ensureToken() error {
	c.Lock()
	defer c.Unlock()

	if c.Token == nil {
		// c.Token will be updated in GetAccessToken call
		if _, err := c.getAccessToken(); err != nil {
			return err
		}
	}

	if !c.tokenExpiresAt.IsZero() && c.tokenExpiresAt.Sub(time.Now()) < RequestNewTokenBeforeExpiresIn {
		// c.Token will be updated in GetAccessToken call
		if _, err := c.getAccessToken(); err != nil {
			return err
		}
	}

	return nil
}

// sendWithAuth makes a request to the API and apply OAuth2 header automatically.
// If the access token soon to be expired or already expired, it will try to get a new one before
// making the main request
// client.Token will be updated when changed
func (c *client) sendWithAuth(req *http.Request, v interface{}) error {
	if err := c.ensureToken(); err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.Token.Token)
	return c.send(req, v)
}

// sendWithBasicAuth makes a request to the API using clientID:secret basic auth
func (c *client) sendWithBasicAuth(req *http.Request, v interface{}) error {
	req.SetBasicAuth(c.ClientID, c.Secret)

	return c.send(req, v)
}

// newRequest constructs a request
// Convert payload to a JSON
func (c *client) newRequest(method, url string, payload interface{}) (*http.Request, error) {
	var buf io.Reader
	if payload != nil {
		b, err := json.Marshal(&payload)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(b)
	}
	return http.NewRequest(method, url, buf)
}

// log will dump request and response to the log file
func (c *client) log(r *http.Request, resp *http.Response) {
	if c.Log != nil {
		var (
			reqDump  string
			respDump []byte
		)

		if r != nil {
			reqDump = fmt.Sprintf("%s %s. Data: %s", r.Method, r.URL.String(), r.Form.Encode())
		}
		if resp != nil {
			respDump, _ = httputil.DumpResponse(resp, true)
		}

		c.Log.Write([]byte(fmt.Sprintf("Request: %s\nResponse: %s\n", reqDump, string(respDump))))
	}
}
