package paypalnvp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	sandboxEndpoint    = "https://api-3t.sandbox.paypal.com/nvp"
	productionEndpoint = "https://api-3t.paypal.com/nvp"
	version            = "204"
)

type Client interface {
	GetBalance(returnAll bool) ([]Amount, error)
}

type client struct {
	username    string
	password    string
	signature   string
	usesSandbox bool
	client      *http.Client
}

func (c *client) performRequest(values url.Values) (*PayPalResponse, error) {
	values.Add("USER", c.username)
	values.Add("PWD", c.password)
	values.Add("SIGNATURE", c.signature)
	values.Add("VERSION", version)

	endpoint := productionEndpoint
	if c.usesSandbox {
		endpoint = sandboxEndpoint
	}

	formResponse, err := c.client.PostForm(endpoint, values)
	defer formResponse.Body.Close()
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(formResponse.Body)
	if err != nil {
		return nil, err
	}

	responseValues, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, err
	}

	response := &PayPalResponse{
		Ack:           responseValues.Get("ACK"),
		CorrelationId: responseValues.Get("CORRELATIONID"),
		Timestamp:     responseValues.Get("TIMESTAMP"),
		Version:       responseValues.Get("VERSION"),
		Build:         responseValues.Get("BUILD"),
		Values:        responseValues,
	}

	errorCode := responseValues.Get("L_ERRORCODE0")
	if len(errorCode) != 0 || strings.ToLower(response.Ack) == "failure" || strings.ToLower(response.Ack) == "failurewithwarning" {
		return nil, &PayPalError{
			Ack:          response.Ack,
			ErrorCode:    errorCode,
			ShortMessage: responseValues.Get("L_SHORTMESSAGE0"),
			LongMessage:  responseValues.Get("L_LONGMESSAGE0"),
			SeverityCode: responseValues.Get("L_SEVERITYCODE0"),
		}
	}

	return response, nil
}

func (c *client) GetBalance(returnAll bool) ([]Amount, error) {
	r := "0"
	if returnAll {
		r = "1"
	}

	values := url.Values{}
	values.Set("METHOD", "GetBalance")
	values.Add("RETURNALLCURRENCIES", r)

	resp, err := c.performRequest(values)
	if err != nil {
		return nil, err
	}

	var i int
	var amounts []Amount

	for {
		code := fmt.Sprintf("L_CURRENCYCODE%d", i)
		amt := fmt.Sprintf("L_AMT%d", i)

		if !resp.Values.Has(code) || !resp.Values.Has(amt) {
			return amounts, nil
		}

		amounts = append(amounts, Amount{
			Currency: resp.Values.Get(code),
			Value:    resp.Values.Get(amt),
		})
		i = i + 1
	}
}

func NewClient(username, password, signature string, usesSandbox bool) Client {
	return &client{username, password, signature, usesSandbox, new(http.Client)}
}
