package paypal

import (
	"fmt"
	"strings"
	"time"
)

const timeFormat = "2006-01-02T15:04:05-0700"

type Time time.Time

func (t Time) MarshalJSON() ([]byte, error) {
	str := time.Time(t).UTC().Format(timeFormat)
	return []byte(str), nil
}

func (t *Time) UnmarshalJSON(s []byte) error {
	str := strings.Trim(string(s), "\"")

	at, err := time.Parse(timeFormat, str)
	if err != nil {
		return err
	}

	*(*time.Time)(t) = at
	return nil
}

func (t Time) String() string { return time.Time(t).String() }

type TransactionQuery struct {
	TransactionId       *string
	TransactionType     *string
	TransactionStatus   *string
	TransactionAmount   *string
	TransactionCurrency *string

	StartDate time.Time
	EndDate   time.Time
	Fields    *string

	PaymentInstrumentType       *string
	BalanceAffectingRecordsOnly *string
	PageSize                    *string
	Page                        *string
}

type TransactionsResponse struct {
	TransactionDetails    []*TransactionDetail `json:"transaction_details,omitempty"`
	AccountNumber         string               `json:"account_number,omitempty"`
	StartDate             *Time                `json:"start_date,omitempty"`
	EndDate               *Time                `json:"end_date,omitempty"`
	LastRefreshedDateTime *Time                `json:"last_refreshed_datetime,omitempty"`
	Page                  int                  `json:"page,omitempty"`
	TotalItems            int                  `json:"total_items,omitempty"`
	TotalPages            int                  `json:"total_pages,omitempty"`
	Links                 []Link               `json:"links,omitempty"`
}

type TransactionDetail struct {
	TransactionInfo *TransactionInfo      `json:"transaction_info,omitempty"`
	PayerInfo       *TransactionPayerInfo `json:"payer_info,omitempty"`
	ShippingInfo    *ShippingInfo         `json:"shipping_info,omitempty"`
	CartInfo        *CartInfo             `json:"cart_info,omitempty"`
	StoreInfo       *StoreInfo            `json:"store_info,omitempty"`
	AuctionInfo     *AuctionInfo          `json:"auction_info,omitempty"`
	IncentiveInfo   *IncentiveInfo        `json:"incentive_info,omitempty"`
}

type TransactionInfo struct {
	PayPalAccountId           string `json:"paypal_account_id,omitempty"`
	TransactionId             string `json:"transaction_id,omitempty"`
	PayPalReferenceId         string `json:"paypal_reference_id,omitempty"`
	PayPalReferenceIdType     string `json:"paypal_reference_id_type,omitempty"`
	TransactionEventCode      string `json:"transaction_event_code,omitempty"`
	TransactionInitiationDate *Time  `json:"transaction_initiation_date,omitempty"`
	TransactionUpdatedDate    *Time  `json:"transaction_updated_date,omitempty"`
	TransactionAmount         *Money `json:"transaction_amount,omitempty"`
	FeeAmount                 *Money `json:"fee_amount,omitempty"`
	DiscountAmount            *Money `json:"discount_amount,omitempty"`
	InsuranceAmount           *Money `json:"insurance_amount,omitempty"`
	SalesTaxAmount            *Money `json:"sales_tax_amount,omitempty"`
	ShippingAmount            *Money `json:"shipping_amount,omitempty"`
	ShippingDiscountAmount    *Money `json:"shipping_discount_amount,omitempty"`
	ShippingTaxAmount         *Money `json:"shipping_tax_amount,omitempty"`
	OtherAmount               *Money `json:"other_amount,omitempty"`
	TipAmount                 *Money `json:"tip_amount,omitempty"`
	TransactionStatus         string `json:"transaction_status,omitempty"`
	TransactionSubject        string `json:"transaction_subject,omitempty"`
	TransactionNote           string `json:"transaction_note,omitempty"`
	PaymentTrackingId         string `json:"payment_tracking_id,omitempty"`
	BankReferenceId           string `json:"bank_reference_id,omitempty"`
	EndingBalance             *Money `json:"ending_balance,omitempty"`
	AvailableBalance          *Money `json:"available_balance,omitempty"`
	InvoiceId                 string `json:"invoice_id,omitempty"`
	CustomField               string `json:"custom_field,omitempty"`
	ProtectionEligibility     string `json:"protection_eligibility,omitempty"`
	CreditTerm                string `json:"credit_term,omitempty"`
	CreditTransactionalFee    *Money `json:"credit_transactional_fee,omitempty"`
	CreditPromotionalFee      *Money `json:"credit_promotional_fee,omitempty"`
	AnnualPercentageRate      string `json:"annual_percentage_rate,omitempty"`
	PaymentMethodType         string `json:"payment_method_type,omitempty"`
}
type TransactionPayerInfo struct {
	AccountId     string     `json:"account_id,omitempty"`
	EmailAddress  string     `json:"email_address,omitempty"`
	PhoneNumber   *Phone     `json:"phone_number,omitempty"`
	AddressStatus string     `json:"address_status,omitempty"`
	PayerStatus   string     `json:"payer_status,omitempty"`
	PayerName     *PayerName `json:"payer_name,omitempty"`
	CountryCode   string     `json:"country_code,omitempty"`
	Address       *Address   `json:"address,omitempty"`
}
type ShippingInfo struct {
	Name                     string   `json:"name,omitempty"`
	Method                   string   `json:"method,omitempty"`
	Address                  *Address `json:"address,omitempty"`
	SecondaryShippingAddress *Address `json:"secondary_shipping_address,omitempty"`
}
type CartInfo struct {
	ItemDetails     []ItemDetail `json:"item_details,omitempty"`
	TaxInclusive    bool         `json:"tax_inclusive,omitempty"`
	PayPalInvoiceId string       `json:"paypal_invoice_id,omitempty"`
}
type StoreInfo struct {
	StoreId    string `json:"store_id,omitempty"`
	TerminalId string `json:"terminal_id,omitempty"`
}
type AuctionInfo struct {
	AuctionSite        string `json:"auction_site,omitempty"`
	AuctionItemSite    string `json:"auction_item_site,omitempty"`
	AuctionBuyerId     string `json:"auction_buyer_id,omitempty"`
	AuctionClosingDate *Time  `json:"auction_closing_date,omitempty"`
}
type IncentiveInfo struct {
	IncentiveDetails []IncentiveDetail `json:"incentive_details,omitempty"`
}
type IncentiveDetail struct {
	IncentiveType        string `json:"incentive_type,omitempty"`
	IncentiveCode        string `json:"incentive_code,omitempty"`
	IncentiveAmount      *Money `json:"incentive_amount,omitempty"`
	IncentiveProgramCode string `json:"incentive_program_code,omitempty"`
}

type ItemDetail struct {
	ItemCode            string                `json:"item_code,omitempty"`
	ItemName            string                `json:"item_name,omitempty"`
	ItemDescription     string                `json:"item_description,omitempty"`
	ItemOptions         string                `json:"item_options,omitempty"`
	ItemQuantity        string                `json:"item_quantity,omitempty"`
	ItemUnitPrice       *Money                `json:"item_unit_price,omitempty"`
	ItemAmount          *Money                `json:"item_amount,omitempty"`
	DiscountAmount      *Money                `json:"discount_amount,omitempty"`
	AdjustmentAmount    *Money                `json:"adjustment_amount,omitempty"`
	GiftWrapAmount      *Money                `json:"gift_wrap_amount,omitempty"`
	TaxPercentage       string                `json:"tax_percentage,omitempty"`
	TaxAmounts          []ItemDetailTaxAmount `json:"tax_amounts,omitempty"`
	BasicShippingAmount *Money                `json:"basic_shipping_amount,omitempty"`
	ExtraShippingAmount *Money                `json:"extra_shipping_amount,omitempty"`
	HandlingAmount      *Money                `json:"handling_amount,omitempty"`
	InsuranceAmount     *Money                `json:"insurance_amount,omitempty"`
	TotalItemAmount     *Money                `json:"total_item_amount,omitempty"`
	InvoiceNumber       string                `json:"invoice_number,omitempty"`
	CheckoutOptions     []CheckoutOption      `json:"checkout_options,omitempty"`
}

type CheckoutOption struct {
	CheckoutOptionName  string `json:"checkout_option_name,omitempty"`
	CheckoutOptionValue string `json:"checkout_option_value,omitempty"`
}

type ItemDetailTaxAmount struct {
	TaxAmount *Money `json:"tax_amount,omitempty"`
}

type Phone struct {
	CountryCode     string `json:"country_code,omitempty"`
	NationalNumber  string `json:"national_number,omitempty"`
	ExtensionNumber string `json:"extension_number,omitempty"`
}

type PayerName struct {
	CountryCode       string `json:"country_code,omitempty"`
	Prefix            string `json:"prefix,omitempty"`
	GivenName         string `json:"given_name,omitempty"`
	Surname           string `json:"surname,omitempty"`
	MiddleName        string `json:"middle_name,omitempty"`
	Suffix            string `json:"suffix,omitempty"`
	AlternateFullName string `json:"alternate_full_name,omitempty"`
	FullName          string `json:"full_name,omitempty"`
}

// Transactions searches the API
// Endpoint: POST /v1/reporting/transactions
func (c *client) Transactions(q *TransactionQuery) (*TransactionsResponse, error) {
	response := &TransactionsResponse{}
	req, err := c.newRequest("GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/reporting/transactions"), nil)
	if err != nil {
		return response, err
	}

	query := req.URL.Query()
	query.Add("start_date", q.StartDate.UTC().Format(timeFormat))
	query.Add("end_date", q.EndDate.UTC().Format(timeFormat))

	if q.TransactionId != nil {
		query.Add("transaction_id", *q.TransactionId)
	}
	if q.TransactionType != nil {
		query.Add("transaction_type", *q.TransactionType)
	}
	if q.TransactionStatus != nil {
		query.Add("transaction_status", *q.TransactionStatus)
	}
	if q.TransactionAmount != nil {
		query.Add("transaction_amount", *q.TransactionAmount)
	}
	if q.TransactionCurrency != nil {
		query.Add("transaction_currency", *q.TransactionCurrency)
	}
	if q.PaymentInstrumentType != nil {
		query.Add("payment_instrument_type", *q.PaymentInstrumentType)
	}
	if q.BalanceAffectingRecordsOnly != nil {
		query.Add("balance_affecting_records_only", *q.BalanceAffectingRecordsOnly)
	}
	if q.PageSize != nil {
		query.Add("page_size", *q.PageSize)
	}
	if q.Page != nil {
		query.Add("page", *q.Page)
	}
	if q.Fields != nil {
		query.Add("fields", *q.Fields)
	} else {
		query.Add("fields", "all")
	}
	req.URL.RawQuery = query.Encode()

	if err := c.sendWithAuth(req, response); err != nil {
		return response, err
	}
	return response, nil
}
