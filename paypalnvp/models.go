package paypalnvp

import "net/url"

type PayPalResponse struct {
	Ack           string
	CorrelationId string
	Timestamp     string
	Version       string
	Build         string
	Values        url.Values
}

type PayPalError struct {
	Ack          string
	ErrorCode    string
	ShortMessage string
	LongMessage  string
	SeverityCode string
}

func (e *PayPalError) Error() string {
	var message string
	if len(e.ErrorCode) != 0 && len(e.ShortMessage) != 0 {
		message = "PayPal Error " + e.ErrorCode + ": " + e.ShortMessage
	} else if len(e.Ack) != 0 {
		message = e.Ack
	} else {
		message = "PayPal is undergoing maintenance.\nPlease try again later."
	}

	return message
}

type Amount struct {
	Currency string
	Value    string
}
