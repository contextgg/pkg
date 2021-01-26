package paypal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	// APIBaseSandBox points to the sandbox (for testing) version of the API
	APIBaseSandBox = "https://api.sandbox.paypal.com"

	// APIBaseLive points to the live version of the API
	APIBaseLive = "https://api.paypal.com"

	// RequestNewTokenBeforeExpiresIn is used by sendWithAuth and try to get new Token when it's about to expire
	RequestNewTokenBeforeExpiresIn = time.Duration(60) * time.Second
)

// Possible values for `no_shipping` in InputFields
//
// https://developer.paypal.com/docs/api/payment-experience/#definition-input_fields
const (
	NoShippingDisplay      uint = 0
	NoShippingHide         uint = 1
	NoShippingBuyerAccount uint = 2
)

// Possible values for `address_override` in InputFields
//
// https://developer.paypal.com/docs/api/payment-experience/#definition-input_fields
const (
	AddrOverrideFromFile uint = 0
	AddrOverrideFromCall uint = 1
)

// Possible values for `landing_page_type` in FlowConfig
//
// https://developer.paypal.com/docs/api/payment-experience/#definition-flow_config
const (
	LandingPageTypeBilling string = "Billing"
	LandingPageTypeLogin   string = "Login"
)

// Possible value for `allowed_payment_method` in PaymentOptions
//
// https://developer.paypal.com/docs/api/payments/#definition-payment_options
const (
	AllowedPaymentUnrestricted         string = "UNRESTRICTED"
	AllowedPaymentInstantFundingSource string = "INSTANT_FUNDING_SOURCE"
	AllowedPaymentImmediatePay         string = "IMMEDIATE_PAY"
)

// Possible value for `intent` in CreateOrder
//
// https://developer.paypal.com/docs/api/orders/v2/#orders_create
const (
	OrderIntentCapture   string = "CAPTURE"
	OrderIntentAuthorize string = "AUTHORIZE"
)

// Possible values for `category` in Item
//
// https://developer.paypal.com/docs/api/orders/v2/#definition-item
const (
	ItemCategoryDigitalGood  string = "DIGITAL_GOODS"
	ItemCategoryPhysicalGood string = "PHYSICAL_GOODS"
)

// Possible values for `shipping_preference` in ApplicationContext
//
// https://developer.paypal.com/docs/api/orders/v2/#definition-application_context
const (
	ShippingPreferenceGetFromFile        string = "GET_FROM_FILE"
	ShippingPreferenceNoShipping         string = "NO_SHIPPING"
	ShippingPreferenceSetProvidedAddress string = "SET_PROVIDED_ADDRESS"
)

const (
	EventPaymentCaptureCompleted       string = "PAYMENT.CAPTURE.COMPLETED"
	EventPaymentCaptureDenied          string = "PAYMENT.CAPTURE.DENIED"
	EventPaymentCaptureRefunded        string = "PAYMENT.CAPTURE.REFUNDED"
	EventMerchantOnboardingCompleted   string = "MERCHANT.ONBOARDING.COMPLETED"
	EventMerchantPartnerConsentRevoked string = "MERCHANT.PARTNER-CONSENT.REVOKED"
)

const (
	OperationAPIIntegration   string = "API_INTEGRATION"
	ProductExpressCheckout    string = "EXPRESS_CHECKOUT"
	IntegrationMethodPayPal   string = "PAYPAL"
	IntegrationTypeThirdParty string = "THIRD_PARTY"
	ConsentShareData          string = "SHARE_DATA_CONSENT"
)

const (
	FeaturePayment               string = "PAYMENT"
	FeatureRefund                string = "REFUND"
	FeatureFuturePayment         string = "FUTURE_PAYMENT"
	FeatureDirectPayment         string = "DIRECT_PAYMENT"
	FeaturePartnerFee            string = "PARTNER_FEE"
	FeatureDelayFunds            string = "DELAY_FUNDS_DISBURSEMENT"
	FeatureReadSellerDispute     string = "READ_SELLER_DISPUTE"
	FeatureUpdateSellerDispute   string = "UPDATE_SELLER_DISPUTE"
	FeatureDisputeReadBuyer      string = "DISPUTE_READ_BUYER"
	FeatureUpdateCustomerDispute string = "UPDATE_CUSTOMER_DISPUTES"
)

const (
	LinkRelSelf      string = "self"
	LinkRelActionURL string = "action_url"
)

type (
	// JSONTime overrides MarshalJson method to format in ISO8601
	JSONTime time.Time

	// Address struct
	Address struct {
		Line1       string `json:"line1"`
		Line2       string `json:"line2,omitempty"`
		City        string `json:"city"`
		CountryCode string `json:"country_code"`
		PostalCode  string `json:"postal_code,omitempty"`
		State       string `json:"state,omitempty"`
		Phone       string `json:"phone,omitempty"`
	}

	// AgreementDetails struct
	AgreementDetails struct {
		OutstandingBalance AmountPayout `json:"outstanding_balance"`
		CyclesRemaining    int          `json:"cycles_remaining,string"`
		CyclesCompleted    int          `json:"cycles_completed,string"`
		NextBillingDate    time.Time    `json:"next_billing_date"`
		LastPaymentDate    time.Time    `json:"last_payment_date"`
		LastPaymentAmount  AmountPayout `json:"last_payment_amount"`
		FinalPaymentDate   time.Time    `json:"final_payment_date"`
		FailedPaymentCount int          `json:"failed_payment_count,string"`
	}

	// Amount struct
	Amount struct {
		Currency string  `json:"currency"`
		Total    string  `json:"total"`
		Details  Details `json:"details,omitempty"`
	}

	// AmountPayout struct
	AmountPayout struct {
		Currency string `json:"currency"`
		Value    string `json:"value"`
	}

	// PaymentMethod struct
	PaymentMethod struct {
		PayerSelected  string `json:"payer_selected,omitempty"`
		PayeePreferred string `json:"payee_preferred,omitempty"`
	}

	// ApplicationContext struct
	ApplicationContext struct {
		BrandName          string         `json:"brand_name,omitempty"`
		Locale             string         `json:"locale,omitempty"`
		LandingPage        string         `json:"landing_page,omitempty"`
		ShippingPreference string         `json:"shipping_preference,omitempty"`
		UserAction         string         `json:"user_action,omitempty"`
		PaymentMethod      *PaymentMethod `json:"payment_method,omitempty"`
		ReturnURL          string         `json:"return_url,omitempty"`
		CancelURL          string         `json:"cancel_url,omitempty"`
	}

	// Authorization struct
	Authorization struct {
		ID               string                `json:"id,omitempty"`
		CustomID         string                `json:"custom_id,omitempty"`
		InvoiceID        string                `json:"invoice_id,omitempty"`
		Status           string                `json:"status,omitempty"`
		StatusDetails    *CaptureStatusDetails `json:"status_details,omitempty"`
		Amount           *PurchaseUnitAmount   `json:"amount,omitempty"`
		SellerProtection *SellerProtection     `json:"seller_protection,omitempty"`
		CreateTime       *time.Time            `json:"create_time,omitempty"`
		UpdateTime       *time.Time            `json:"update_time,omitempty"`
		ExpirationTime   *time.Time            `json:"expiration_time,omitempty"`
		Links            []Link                `json:"links,omitempty"`
	}

	// AuthorizeOrderResponse .
	AuthorizeOrderResponse struct {
		CreateTime    *time.Time             `json:"create_time,omitempty"`
		UpdateTime    *time.Time             `json:"update_time,omitempty"`
		ID            string                 `json:"id,omitempty"`
		Status        string                 `json:"status,omitempty"`
		Intent        string                 `json:"intent,omitempty"`
		PurchaseUnits []PurchaseUnitRequest  `json:"purchase_units,omitempty"`
		Payer         *PayerWithNameAndPhone `json:"payer,omitempty"`
	}

	// AuthorizeOrderRequest - https://developer.paypal.com/docs/api/orders/v2/#orders_authorize
	AuthorizeOrderRequest struct {
		PaymentSource      *PaymentSource     `json:"payment_source,omitempty"`
		ApplicationContext ApplicationContext `json:"application_context,omitempty"`
	}

	// https://developer.paypal.com/docs/api/payments/v2/#definition-platform_fee
	PlatformFee struct {
		Amount *Money          `json:"amount,omitempty"`
		Payee  *PayeeForOrders `json:"payee,omitempty"`
	}

	// https://developer.paypal.com/docs/api/payments/v2/#definition-payment_instruction
	PaymentInstruction struct {
		PlatformFees     []PlatformFee `json:"platform_fees,omitempty"`
		DisbursementMode string        `json:"disbursement_mode,omitempty"`
	}

	// https://developer.paypal.com/docs/api/payments/v2/#authorizations_capture
	PaymentCaptureRequest struct {
		InvoiceID      string `json:"invoice_id,omitempty"`
		NoteToPayer    string `json:"note_to_payer,omitempty"`
		SoftDescriptor string `json:"soft_descriptor,omitempty"`
		Amount         *Money `json:"amount,omitempty"`
		FinalCapture   bool   `json:"final_capture,omitempty"`
	}

	SellerProtection struct {
		Status            string   `json:"status,omitempty"`
		DisputeCategories []string `json:"dispute_categories,omitempty"`
	}

	// https://developer.paypal.com/docs/api/payments/v2/#definition-capture_status_details
	CaptureStatusDetails struct {
		Reason string `json:"reason,omitempty"`
	}

	PaymentCaptureResponse struct {
		Status           string                `json:"status,omitempty"`
		StatusDetails    *CaptureStatusDetails `json:"status_details,omitempty"`
		ID               string                `json:"id,omitempty"`
		Amount           *Money                `json:"amount,omitempty"`
		InvoiceID        string                `json:"invoice_id,omitempty"`
		FinalCapture     bool                  `json:"final_capture,omitempty"`
		DisbursementMode string                `json:"disbursement_mode,omitempty"`
		Links            []Link                `json:"links,omitempty"`
	}

	// CaptureOrderRequest - https://developer.paypal.com/docs/api/orders/v2/#orders_capture
	CaptureOrderRequest struct {
		PaymentSource *PaymentSource `json:"payment_source"`
	}

	// BatchHeader struct
	BatchHeader struct {
		Amount            *AmountPayout      `json:"amount,omitempty"`
		Fees              *AmountPayout      `json:"fees,omitempty"`
		PayoutBatchID     string             `json:"payout_batch_id,omitempty"`
		BatchStatus       string             `json:"batch_status,omitempty"`
		TimeCreated       *time.Time         `json:"time_created,omitempty"`
		TimeCompleted     *time.Time         `json:"time_completed,omitempty"`
		SenderBatchHeader *SenderBatchHeader `json:"sender_batch_header,omitempty"`
	}

	// BillingAgreement struct
	BillingAgreement struct {
		Name                        string               `json:"name,omitempty"`
		Description                 string               `json:"description,omitempty"`
		StartDate                   JSONTime             `json:"start_date,omitempty"`
		Plan                        BillingPlan          `json:"plan,omitempty"`
		Payer                       Payer                `json:"payer,omitempty"`
		ShippingAddress             *ShippingAddress     `json:"shipping_address,omitempty"`
		OverrideMerchantPreferences *MerchantPreferences `json:"override_merchant_preferences,omitempty"`
	}

	// BillingPlan struct
	BillingPlan struct {
		ID          string `json:"id,omitempty"`
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`

		PaymentDefinitions  []PaymentDefinition  `json:"payment_definitions,omitempty"`
		MerchantPreferences *MerchantPreferences `json:"merchant_preferences,omitempty"`
	}

	// Capture struct
	Capture struct {
		ID                        string                     `json:"id,omitempty"`
		Status                    string                     `json:"status,omitempty"`
		Amount                    *Money                     `json:"amount,omitempty"`
		IsFinalCapture            bool                       `json:"is_final_capture"`
		SellerReceivableBreakdown *SellerReceivableBreakdown `json:"seller_receivable_breakdown,omitempty"`
		Links                     []Link                     `json:"links,omitempty"`
		CreateTime                *time.Time                 `json:"create_time,omitempty"`
		UpdateTime                *time.Time                 `json:"update_time,omitempty"`
	}

	// SellerReceivableBreakdown struct
	SellerReceivableBreakdown struct {
		GrossAmount      *Money `json:"gross_amount,omitempty"`
		PaypalFee        *Money `json:"paypal_fee,omitempty"`
		NetAmount        *Money `json:"net_amount,omitempty"`
		ReceivableAmount *Money `json:"receivable_amount,omitempty"`
		ExchangeRate     *Money `json:"exchange_rate,omitempty"`
		PlatformFees     *Money `json:"platform_fees,omitempty"`
	}

	// ChargeModel struct
	ChargeModel struct {
		Type   string       `json:"type,omitempty"`
		Amount AmountPayout `json:"amount,omitempty"`
	}

	// CreditCard struct
	CreditCard struct {
		ID                 string   `json:"id,omitempty"`
		PayerID            string   `json:"payer_id,omitempty"`
		ExternalCustomerID string   `json:"external_customer_id,omitempty"`
		Number             string   `json:"number"`
		Type               string   `json:"type"`
		ExpireMonth        string   `json:"expire_month"`
		ExpireYear         string   `json:"expire_year"`
		CVV2               string   `json:"cvv2,omitempty"`
		FirstName          string   `json:"first_name,omitempty"`
		LastName           string   `json:"last_name,omitempty"`
		BillingAddress     *Address `json:"billing_address,omitempty"`
		State              string   `json:"state,omitempty"`
		ValidUntil         string   `json:"valid_until,omitempty"`
	}

	// CreditCards GET /v1/vault/credit-cards
	CreditCards struct {
		Items      []CreditCard `json:"items"`
		Links      []Link       `json:"links"`
		TotalItems int          `json:"total_items"`
		TotalPages int          `json:"total_pages"`
	}

	// CreditCardToken struct
	CreditCardToken struct {
		CreditCardID string `json:"credit_card_id"`
		PayerID      string `json:"payer_id,omitempty"`
		Last4        string `json:"last4,omitempty"`
		ExpireYear   string `json:"expire_year,omitempty"`
		ExpireMonth  string `json:"expire_month,omitempty"`
	}

	// CreditCardsFilter struct
	CreditCardsFilter struct {
		PageSize int
		Page     int
	}

	// CreditCardField PATCH /v1/vault/credit-cards/credit_card_id
	CreditCardField struct {
		Operation string `json:"op"`
		Path      string `json:"path"`
		Value     string `json:"value"`
	}

	// Currency struct
	Currency struct {
		Currency string `json:"currency,omitempty"`
		Value    string `json:"value,omitempty"`
	}

	// Details structure used in Amount structures as optional value
	Details struct {
		Subtotal         string `json:"subtotal,omitempty"`
		Shipping         string `json:"shipping,omitempty"`
		Tax              string `json:"tax,omitempty"`
		HandlingFee      string `json:"handling_fee,omitempty"`
		ShippingDiscount string `json:"shipping_discount,omitempty"`
		Insurance        string `json:"insurance,omitempty"`
		GiftWrap         string `json:"gift_wrap,omitempty"`
	}

	// ErrorResponseDetail struct
	ErrorResponseDetail struct {
		Field string `json:"field"`
		Issue string `json:"issue"`
		Links []Link `json:"link"`
	}

	// ErrorResponse https://developer.paypal.com/docs/api/errors/
	ErrorResponse struct {
		Response        *http.Response        `json:"-"`
		Name            string                `json:"name"`
		DebugID         string                `json:"debug_id"`
		Message         string                `json:"message"`
		InformationLink string                `json:"information_link"`
		Details         []ErrorResponseDetail `json:"details"`
	}

	// ExecuteAgreementResponse struct
	ExecuteAgreementResponse struct {
		ID               string           `json:"id"`
		State            string           `json:"state"`
		Description      string           `json:"description,omitempty"`
		Payer            Payer            `json:"payer"`
		Plan             BillingPlan      `json:"plan"`
		StartDate        time.Time        `json:"start_date"`
		ShippingAddress  ShippingAddress  `json:"shipping_address"`
		AgreementDetails AgreementDetails `json:"agreement_details"`
		Links            []Link           `json:"links"`
	}

	// FundingInstrument struct
	FundingInstrument struct {
		CreditCard      *CreditCard      `json:"credit_card,omitempty"`
		CreditCardToken *CreditCardToken `json:"credit_card_token,omitempty"`
	}

	// Item struct
	Item struct {
		Name        string `json:"name"`
		UnitAmount  *Money `json:"unit_amount,omitempty"`
		Tax         *Money `json:"tax,omitempty"`
		Quantity    string `json:"quantity"`
		Description string `json:"description,omitempty"`
		SKU         string `json:"sku,omitempty"`
		Category    string `json:"category,omitempty"`
	}

	// ItemList struct
	ItemList struct {
		Items           []Item           `json:"items,omitempty"`
		ShippingAddress *ShippingAddress `json:"shipping_address,omitempty"`
	}

	// Link struct
	Link struct {
		Href        string `json:"href"`
		Rel         string `json:"rel,omitempty"`
		Method      string `json:"method,omitempty"`
		Description string `json:"description,omitempty"`
		Enctype     string `json:"enctype,omitempty"`
	}

	// PurchaseUnitAmount struct
	PurchaseUnitAmount struct {
		Currency  string                       `json:"currency_code"`
		Value     string                       `json:"value"`
		Breakdown *PurchaseUnitAmountBreakdown `json:"breakdown,omitempty"`
	}

	// PurchaseUnitAmountBreakdown struct
	PurchaseUnitAmountBreakdown struct {
		ItemTotal        *Money `json:"item_total,omitempty"`
		Shipping         *Money `json:"shipping,omitempty"`
		Handling         *Money `json:"handling,omitempty"`
		TaxTotal         *Money `json:"tax_total,omitempty"`
		Insurance        *Money `json:"insurance,omitempty"`
		ShippingDiscount *Money `json:"shipping_discount,omitempty"`
		Discount         *Money `json:"discount,omitempty"`
	}

	// Money struct
	//
	// https://developer.paypal.com/docs/api/orders/v2/#definition-money
	Money struct {
		Currency string `json:"currency_code"`
		Value    string `json:"value"`
	}

	// PurchaseUnit struct
	PurchaseUnit struct {
		ReferenceID string              `json:"reference_id"`
		Amount      *PurchaseUnitAmount `json:"amount,omitempty"`
		Payments    *Payments           `json:"payments,omitempty"`
	}

	// Payments struct
	Payments struct {
		Captures []Capture `json:"captures,omitempty"`
	}

	// TaxInfo used for orders.
	TaxInfo struct {
		TaxID     string `json:"tax_id,omitempty"`
		TaxIDType string `json:"tax_id_type,omitempty"`
	}

	// PhoneWithTypeNumber struct for PhoneWithType
	PhoneWithTypeNumber struct {
		NationalNumber string `json:"national_number,omitempty"`
	}

	// PhoneWithType struct used for orders
	PhoneWithType struct {
		PhoneType   string               `json:"phone_type,omitempty"`
		PhoneNumber *PhoneWithTypeNumber `json:"phone_number,omitempty"`
	}

	// CreateProductInput struct
	CreateProductInput struct {
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
		Type        string `json:"type,omitempty"`
		Category    string `json:"category,omitempty"`
		ImageURL    string `json:"image_url,omitempty"`
		HomeURL     string `json:"home_url,omitempty"`
	}

	// ListProductsItem struct each product in list when fetching products
	ListProductsItem struct {
		ID          string    `json:"id,omitempty"`
		Name        string    `json:"name,omitempty"`
		Description string    `json:"description,omitempty"`
		CreateTime  time.Time `json:"create_time,omitempty"`
		Links       []Link    `json:"links,omitempty"`
	}

	// ProductResp struct response when fetching a single product
	ProductResp struct {
		ID          string    `json:"id,omitempty"`
		Name        string    `json:"name,omitempty"`
		Description string    `json:"description,omitempty"`
		Type        string    `json:"type,omitempty"`
		Category    string    `json:"category,omitempty"`
		ImageURL    string    `json:"image_url,omitempty"`
		HomeURL     string    `json:"home_url,omitempty"`
		CreateTime  time.Time `json:"create_time,omitempty"`
		UpdateTime  time.Time `json:"update_time,omitempty"`
		Links       []Link    `json:"links,omitempty"`
	}

	// CycleExecution struct
	CycleExecution struct {
		TenureType      string `json:"tenure_type,omitempty"`
		Sequence        int    `json:"sequence,omitempty"`
		CyclesCompleted int    `json:"cycles_completed,omitempty"`
		CyclesRemaining int    `json:"cycles_remaining,omitempty"`
		TotalCycles     int    `json:"total_cycles,omitempty"`
	}

	// LastPayment struct
	LastPayment struct {
		Amount *AmountPayout `json:"amount,omitempty"`
		Time   time.Time     `json:"time,omitempty"`
	}

	// BillingInfo struct
	BillingInfo struct {
		OutstandingBalance  *AmountPayout    `json:"outstanding_balance,omitempty"`
		CycleExecutions     []CycleExecution `json:"cycle_executions,omitempty"`
		LastPayment         *LastPayment     `json:"last_payment,omitempty"`
		NextBillingTime     time.Time        `json:"next_billing_time,omitempty"`
		FailedPaymentsCount int              `json:"failed_payments_count,omitempty"`
	}

	// PlanAmount struct
	PlanAmount struct {
		Value        string `json:"value,omitempty"`
		CurrencyCode string `json:"currency_code,omitempty"`
	}

	// PlanTaxes struct
	PlanTaxes struct {
		Percentage string `json:"percentage,omitempty"`
		Inclusive  bool   `json:"inclusive,omitempty"`
	}

	// PricingScheme struct
	PricingScheme struct {
		FixedPrice *PlanAmount `json:"fixed_price,omitempty"`
		Status     string      `json:"status,omitempty"`
		Version    int         `json:"version,omitempty"`
		CreateTime time.Time   `json:"create_time,omitempty"`
		UpdateTime time.Time   `json:"update_time,omitempty"`
	}

	// BillingCycle struct
	BillingCycle struct {
		Frequency     *BillingFrequency `json:"frequency"`
		TenureType    string            `json:"tenure_type"`
		Sequence      int               `json:"sequence"`
		TotalCycles   int               `json:"total_cycles"`
		PricingScheme *PricingScheme    `json:"pricing_scheme,omitempty"`
	}

	// BillingFrequency struct
	BillingFrequency struct {
		IntervalUnit  string `json:"interval_unit,omitempty"`
		IntervalCount int    `json:"interval_count,omitempty"`
	}

	// PaymentPreferences struct
	PaymentPreferences struct {
		AutoBillOutstanding     bool        `json:"auto_bill_outstanding"`
		SetupFee                *PlanAmount `json:"setup_fee,omitempty"`
		SetupFeeFailureAction   string      `json:"setup_fee_failure_action,omitempty"`
		PaymentFailureThreshold int         `json:"payment_failure_threshold,omitempty"`
	}

	// SubscriptionPlanResp struct
	SubscriptionPlanResp struct {
		ID                 string              `json:"id,omitempty"`
		ProductID          string              `json:"product_id,omitempty"`
		Name               string              `json:"name,omitempty"`
		Status             string              `json:"status,omitempty"`
		Description        string              `json:"description,omitempty"`
		BillingCycles      []BillingCycle      `json:"billing_cycles,omitempty"`
		PaymentPreferences *PaymentPreferences `json:"payment_preferences,omitempty"`
		Taxes              *PlanTaxes          `json:"taxes,omitempty"`
		Links              []Link              `json:"links,omitempty"`
	}

	// ListSubscriptionPlan struct
	ListSubscriptionPlan struct {
		ID          string    `json:"id,omitempty"`
		ProductID   string    `json:"product_id,omitempty"`
		Name        string    `json:"name,omitempty"`
		Status      string    `json:"status,omitempty"`
		Description string    `json:"description,omitempty"`
		CreateTime  time.Time `json:"create_time,omitempty"`
		Links       []Link    `json:"links,omitempty"`
	}

	// ListSubscriptionPlansResp struct
	ListSubscriptionPlansResp struct {
		Plans      []ListSubscriptionPlan `json:"plans,omitempty"`
		TotalItems string                 `json:"total_items,omitempty"`
		TotalPages string                 `json:"total_pages,omitempty"`
		Links      []Link                 `json:"links,omitempty"`
	}

	// UpdatedPricingScheme struct
	UpdatedPricingScheme struct {
		BillingCycleSequence int            `json:"billing_cycle_sequence,omitempty"`
		PricingScheme        *PricingScheme `json:"pricing_scheme,omitempty"`
	}

	// UpdatePricingSchemeInput struct
	UpdatePricingSchemeInput struct {
		PricingSchemes []UpdatedPricingScheme
	}

	// SubscriptionPlanInput struct
	SubscriptionPlanInput struct {
		ProductID          string              `json:"product_id"`
		Name               string              `json:"name"`
		Status             string              `json:"status,omitempty"`
		Description        string              `json:"description,omitempty"`
		BillingCycles      []*BillingCycle     `json:"billing_cycles"`
		PaymentPreferences *PaymentPreferences `json:"payment_preferences"`
		Taxes              *PlanTaxes          `json:"taxes,omitempty"`
	}

	// SubscriberResp struct subscriber data returned when creating a new subscription
	SubscriberResp struct {
		Name         *Name  `json:"name,omitempty"`
		EmailAddress string `json:"email_address,omitempty"`
		PayerID      string `json:"payer_id,omitempty"`
	}

	// SubscriberParams struct subscriber params for creating a new subscription
	SubscriberParams struct {
		Name         *Name  `json:"name,omitempty"`
		EmailAddress string `json:"email_address,omitempty"`
		PayerID      string `json:"payer_id,omitempty"`
	}

	// SubscriptionResp struct
	SubscriptionResp struct {
		ID             string          `json:"id"`
		Status         string          `json:"status"`
		PlanID         string          `json:"plan_id"`
		StartTime      string          `json:"start_time"`
		Quantity       string          `json:"quantity"`
		ShippingAmount *Money          `json:"shipping_amount"`
		Subscriber     *SubscriberResp `json:"subscriber"`
		BillingInfo    *BillingInfo    `json:"billing_info"`
		CreateTime     string          `json:"create_time"`
		UpdateTime     string          `json:"update_time"`
		Links          []Link          `json:"links"`
	}

	// SubscriptionInput struct for new subscription input
	SubscriptionInput struct {
		PlanID             string              `json:"plan_id,omitempty"`
		StartTime          string              `json:"start_time,omitempty"`
		Quantity           string              `json:"quantity,omitempty"`
		ShippingAmount     *PlanAmount         `json:"shipping_amount,omitempty"`
		Subscriber         *SubscriberParams   `json:"subscriber,omitempty"`
		ApplicationContext *ApplicationContext `json:"application_context,omitempty"`
	}

	// CreateSubscriptionResp struct response after creating a subscription
	CreateSubscriptionResp struct {
		ID               string          `json:"id,omitempty"`
		Status           string          `json:"status,omitempty"`
		StatusUpdateTime time.Time       `json:"status_update_time,omitempty"`
		PlanID           string          `json:"plan_id,omitempty"`
		StartTime        time.Time       `json:"start_time,omitempty"`
		Quantity         string          `json:"quantity,omitempty"`
		ShippingAmount   *AmountPayout   `json:"shipping_amount,omitempty"`
		Subscriber       *SubscriberResp `json:"subscriber,omitempty"`
		CreateTime       time.Time       `json:"create_time,omitempty"`
		Links            []Link          `json:"links,omitempty"`
	}

	SubscriptionTransactionsResp struct {
		Transactions []*Transaction `json:"transactions,omitempty"`
		Links        []Link         `json:"links,omitempty"`
	}

	Transaction struct {
		Status              string              `json:"status,omitempty"`
		ID                  string              `json:"id,omitempty"`
		AmountWithBreakdown AmountWithBreakdown `json:"amount_with_breakdown,omitempty"`
		PayerName           *Name               `json:"payer_name,omitempty"`
		PayerEmail          string              `json:"payer_email,omitempty"`
		Time                time.Time           `json:"time,omitempty"`
	}

	// SellerReceivableBreakdown struct
	AmountWithBreakdown struct {
		GrossAmount    *Money `json:"gross_amount,omitempty"`
		FeeAmount      *Money `json:"fee_amount,omitempty"`
		ShippingAmount *Money `json:"shipping_amount,omitempty"`
		TaxAmount      *Money `json:"tax_amount,omitempty"`
		NetAmount      *Money `json:"net_amount,omitempty"`
	}

	// CreateOrderPayer used with create order requests
	CreateOrderPayer struct {
		Name         *Name                          `json:"name,omitempty"`
		EmailAddress string                         `json:"email_address,omitempty"`
		PayerID      string                         `json:"payer_id,omitempty"`
		Phone        *PhoneWithType                 `json:"phone,omitempty"`
		BirthDate    string                         `json:"birth_date,omitempty"`
		TaxInfo      *TaxInfo                       `json:"tax_info,omitempty"`
		Address      *ShippingDetailAddressPortable `json:"address,omitempty"`
	}

	// PurchaseUnitRequest struct
	PurchaseUnitRequest struct {
		ReferenceID    string              `json:"reference_id,omitempty"`
		Amount         *PurchaseUnitAmount `json:"amount"`
		Payee          *PayeeForOrders     `json:"payee,omitempty"`
		Description    string              `json:"description,omitempty"`
		CustomID       string              `json:"custom_id,omitempty"`
		InvoiceID      string              `json:"invoice_id,omitempty"`
		SoftDescriptor string              `json:"soft_descriptor,omitempty"`
		Items          []Item              `json:"items,omitempty"`
		Shipping       *ShippingDetail     `json:"shipping,omitempty"`
	}

	// MerchantPreferences struct
	MerchantPreferences struct {
		SetupFee                *AmountPayout `json:"setup_fee,omitempty"`
		ReturnURL               string        `json:"return_url,omitempty"`
		CancelURL               string        `json:"cancel_url,omitempty"`
		AutoBillAmount          string        `json:"auto_bill_amount,omitempty"`
		InitialFailAmountAction string        `json:"initial_fail_amount_action,omitempty"`
		MaxFailAttempts         string        `json:"max_fail_attempts,omitempty"`
	}

	// Order struct
	Order struct {
		ID            string         `json:"id,omitempty"`
		Status        string         `json:"status,omitempty"`
		Intent        string         `json:"intent,omitempty"`
		PurchaseUnits []PurchaseUnit `json:"purchase_units,omitempty"`
		Links         []Link         `json:"links,omitempty"`
		CreateTime    *time.Time     `json:"create_time,omitempty"`
		UpdateTime    *time.Time     `json:"update_time,omitempty"`
	}

	// PayerWithNameAndPhone struct
	PayerWithNameAndPhone struct {
		Name         *Name          `json:"name,omitempty"`
		EmailAddress string         `json:"email_address,omitempty"`
		Phone        *PhoneWithType `json:"phone,omitempty"`
		PayerID      string         `json:"payer_id,omitempty"`
	}

	// CaptureOrderResponse is the response for capture order
	CaptureOrderResponse struct {
		ID            string                 `json:"id,omitempty"`
		Status        string                 `json:"status,omitempty"`
		Payer         *PayerWithNameAndPhone `json:"payer,omitempty"`
		PurchaseUnits []PurchaseUnit         `json:"purchase_units,omitempty"`
	}

	// Payer struct
	Payer struct {
		PaymentMethod      string              `json:"payment_method"`
		FundingInstruments []FundingInstrument `json:"funding_instruments,omitempty"`
		PayerInfo          *PayerInfo          `json:"payer_info,omitempty"`
		Status             string              `json:"payer_status,omitempty"`
	}

	// PayerInfo struct
	PayerInfo struct {
		Email           string           `json:"email,omitempty"`
		FirstName       string           `json:"first_name,omitempty"`
		LastName        string           `json:"last_name,omitempty"`
		PayerID         string           `json:"payer_id,omitempty"`
		Phone           string           `json:"phone,omitempty"`
		ShippingAddress *ShippingAddress `json:"shipping_address,omitempty"`
		TaxIDType       string           `json:"tax_id_type,omitempty"`
		TaxID           string           `json:"tax_id,omitempty"`
		CountryCode     string           `json:"country_code"`
	}

	// PaymentDefinition struct
	PaymentDefinition struct {
		ID                string        `json:"id,omitempty"`
		Name              string        `json:"name,omitempty"`
		Type              string        `json:"type,omitempty"`
		Frequency         string        `json:"frequency,omitempty"`
		FrequencyInterval string        `json:"frequency_interval,omitempty"`
		Amount            AmountPayout  `json:"amount,omitempty"`
		Cycles            string        `json:"cycles,omitempty"`
		ChargeModels      []ChargeModel `json:"charge_models,omitempty"`
	}

	// PaymentOptions struct
	PaymentOptions struct {
		AllowedPaymentMethod string `json:"allowed_payment_method,omitempty"`
	}

	// PaymentPatch PATCH /v2/payments/payment/{payment_id)
	PaymentPatch struct {
		Operation string      `json:"op"`
		Path      string      `json:"path"`
		Value     interface{} `json:"value"`
	}

	// PaymentPayer struct
	PaymentPayer struct {
		PaymentMethod string     `json:"payment_method"`
		Status        string     `json:"status,omitempty"`
		PayerInfo     *PayerInfo `json:"payer_info,omitempty"`
	}

	// PaymentSource structure
	PaymentSource struct {
		Card  *PaymentSourceCard  `json:"card"`
		Token *PaymentSourceToken `json:"token"`
	}

	// PaymentSourceCard structure
	PaymentSourceCard struct {
		ID             string              `json:"id"`
		Name           string              `json:"name"`
		Number         string              `json:"number"`
		Expiry         string              `json:"expiry"`
		SecurityCode   string              `json:"security_code"`
		LastDigits     string              `json:"last_digits"`
		CardType       string              `json:"card_type"`
		BillingAddress *CardBillingAddress `json:"billing_address"`
	}

	// CardBillingAddress structure
	CardBillingAddress struct {
		AddressLine1 string `json:"address_line_1"`
		AddressLine2 string `json:"address_line_2"`
		AdminArea2   string `json:"admin_area_2"`
		AdminArea1   string `json:"admin_area_1"`
		PostalCode   string `json:"postal_code"`
		CountryCode  string `json:"country_code"`
	}

	// PaymentSourceToken structure
	PaymentSourceToken struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	}

	// Payout struct
	Payout struct {
		SenderBatchHeader *SenderBatchHeader `json:"sender_batch_header"`
		Items             []PayoutItem       `json:"items"`
	}

	// PayoutItem struct
	PayoutItem struct {
		RecipientType string        `json:"recipient_type"`
		Receiver      string        `json:"receiver"`
		Amount        *AmountPayout `json:"amount"`
		Note          string        `json:"note,omitempty"`
		SenderItemID  string        `json:"sender_item_id,omitempty"`
	}

	// PayoutItemResponse struct
	PayoutItemResponse struct {
		PayoutItemID      string        `json:"payout_item_id"`
		TransactionID     string        `json:"transaction_id"`
		TransactionStatus string        `json:"transaction_status"`
		PayoutBatchID     string        `json:"payout_batch_id,omitempty"`
		PayoutItemFee     *AmountPayout `json:"payout_item_fee,omitempty"`
		PayoutItem        *PayoutItem   `json:"payout_item"`
		TimeProcessed     *time.Time    `json:"time_processed,omitempty"`
		Links             []Link        `json:"links"`
		Error             ErrorResponse `json:"errors,omitempty"`
	}

	// PayoutResponse struct
	PayoutResponse struct {
		BatchHeader *BatchHeader         `json:"batch_header"`
		Items       []PayoutItemResponse `json:"items"`
		Links       []Link               `json:"links"`
	}

	// RedirectURLs struct
	RedirectURLs struct {
		ReturnURL string `json:"return_url,omitempty"`
		CancelURL string `json:"cancel_url,omitempty"`
	}

	// Refund struct
	Refund struct {
		ID                     string                 `json:"id,omitempty"`
		Status                 string                 `json:"status,omitempty"`
		Amount                 *Money                 `json:"amount,omitempty"`
		SellerPayableBreakdown SellerPayableBreakdown `json:"seller_payable_breakdown"`
		CreateTime             *time.Time             `json:"create_time,omitempty"`
		UpdateTime             *time.Time             `json:"update_time,omitempty"`
	}

	// SellerPayableBreakdown struct
	SellerPayableBreakdown struct {
		GrossAmount *Money `json:"gross_amount,omitempty"`
		PaypalFee   *Money `json:"paypal_fee,omitempty"`
		NetAmount   *Money `json:"net_amount,omitempty"`
	}

	// RefundResponse .
	RefundResponse struct {
		ID     string              `json:"id,omitempty"`
		Amount *PurchaseUnitAmount `json:"amount,omitempty"`
		Status string              `json:"status,omitempty"`
	}

	// SenderBatchHeader struct
	SenderBatchHeader struct {
		SenderBatchID string `json:"sender_batch_id,omitempty"`
		RecipientType string `json:"recipient_type,omitempty"`
		EmailSubject  string `json:"email_subject,omitempty"`
		EmailMessage  string `json:"email_message,omitempty"`
	}

	// ShippingAddress struct
	ShippingAddress struct {
		RecipientName string `json:"recipient_name,omitempty"`
		Type          string `json:"type,omitempty"`
		Line1         string `json:"line1"`
		Line2         string `json:"line2,omitempty"`
		City          string `json:"city"`
		CountryCode   string `json:"country_code"`
		PostalCode    string `json:"postal_code,omitempty"`
		State         string `json:"state,omitempty"`
		Phone         string `json:"phone,omitempty"`
	}

	// ShippingDetailAddressPortable used with create orders
	ShippingDetailAddressPortable struct {
		AddressLine1 string `json:"address_line_1,omitempty"`
		AddressLine2 string `json:"address_line_2,omitempty"`
		AdminArea1   string `json:"admin_area_1,omitempty"`
		AdminArea2   string `json:"admin_area_2,omitempty"`
		PostalCode   string `json:"postal_code,omitempty"`
		CountryCode  string `json:"country_code,omitempty"`
	}

	// Name struct
	Name struct {
		Prefix     string `json:"prefix,omitempty"`
		GivenName  string `json:"given_name,omitempty"`
		Surname    string `json:"surname,omitempty"`
		MiddleName string `json:"middle_name,omitempty"`
		Suffix     string `json:"suffix,omitempty"`
		FullName   string `json:"full_name,omitempty"`
	}

	// ShippingDetail struct
	ShippingDetail struct {
		Name    *Name                          `json:"name,omitempty"`
		Address *ShippingDetailAddressPortable `json:"address,omitempty"`
	}

	expirationTime int64

	// TokenResponse is for API response for the /oauth2/token endpoint
	TokenResponse struct {
		RefreshToken string         `json:"refresh_token"`
		Token        string         `json:"access_token"`
		Type         string         `json:"token_type"`
		ExpiresIn    expirationTime `json:"expires_in"`
	}

	//Payee struct
	Payee struct {
		Email string `json:"email"`
	}

	// PayeeForOrders struct
	PayeeForOrders struct {
		EmailAddress string `json:"email_address,omitempty"`
		MerchantID   string `json:"merchant_id,omitempty"`
	}

	// UserInfo struct
	UserInfo struct {
		ID              string   `json:"user_id"`
		Name            string   `json:"name"`
		GivenName       string   `json:"given_name"`
		FamilyName      string   `json:"family_name"`
		Email           string   `json:"email"`
		Verified        bool     `json:"verified,omitempty,string"`
		Gender          string   `json:"gender,omitempty"`
		BirthDate       string   `json:"birthdate,omitempty"`
		ZoneInfo        string   `json:"zoneinfo,omitempty"`
		Locale          string   `json:"locale,omitempty"`
		Phone           string   `json:"phone_number,omitempty"`
		Address         *Address `json:"address,omitempty"`
		VerifiedAccount bool     `json:"verified_account,omitempty,string"`
		AccountType     string   `json:"account_type,omitempty"`
		AgeRange        string   `json:"age_range,omitempty"`
		PayerID         string   `json:"payer_id,omitempty"`
	}

	// WebProfile represents the configuration of the payment web payment experience
	//
	// https://developer.paypal.com/docs/api/payment-experience/
	WebProfile struct {
		ID           string       `json:"id,omitempty"`
		Name         string       `json:"name"`
		Presentation Presentation `json:"presentation,omitempty"`
		InputFields  InputFields  `json:"input_fields,omitempty"`
		FlowConfig   FlowConfig   `json:"flow_config,omitempty"`
	}

	// Presentation represents the branding and locale that a customer sees on
	// redirect payments
	//
	// https://developer.paypal.com/docs/api/payment-experience/#definition-presentation
	Presentation struct {
		BrandName  string `json:"brand_name,omitempty"`
		LogoImage  string `json:"logo_image,omitempty"`
		LocaleCode string `json:"locale_code,omitempty"`
	}

	// InputFields represents the fields that are displayed to a customer on
	// redirect payments
	//
	// https://developer.paypal.com/docs/api/payment-experience/#definition-input_fields
	InputFields struct {
		AllowNote       bool `json:"allow_note,omitempty"`
		NoShipping      uint `json:"no_shipping,omitempty"`
		AddressOverride uint `json:"address_override,omitempty"`
	}

	// FlowConfig represents the general behaviour of redirect payment pages
	//
	// https://developer.paypal.com/docs/api/payment-experience/#definition-flow_config
	FlowConfig struct {
		LandingPageType   string `json:"landing_page_type,omitempty"`
		BankTXNPendingURL string `json:"bank_txn_pending_url,omitempty"`
		UserAction        string `json:"user_action,omitempty"`
	}

	VerifyWebhookResponse struct {
		VerificationStatus string `json:"verification_status,omitempty"`
	}

	WebhookEvent struct {
		ID              string    `json:"id"`
		CreateTime      time.Time `json:"create_time"`
		ResourceType    string    `json:"resource_type"`
		EventType       string    `json:"event_type"`
		Summary         string    `json:"summary,omitempty"`
		Resource        Resource  `json:"resource"`
		Links           []Link    `json:"links"`
		EventVersion    string    `json:"event_version,omitempty"`
		ResourceVersion string    `json:"resource_version,omitempty"`
	}

	Resource struct {
		// Payment Resource type
		ID                     string                  `json:"id,omitempty"`
		Status                 string                  `json:"status,omitempty"`
		StatusDetails          *CaptureStatusDetails   `json:"status_details,omitempty"`
		Amount                 *PurchaseUnitAmount     `json:"amount,omitempty"`
		UpdateTime             string                  `json:"update_time,omitempty"`
		CreateTime             string                  `json:"create_time,omitempty"`
		ExpirationTime         string                  `json:"expiration_time,omitempty"`
		SellerProtection       *SellerProtection       `json:"seller_protection,omitempty"`
		FinalCapture           bool                    `json:"final_capture,omitempty"`
		SellerPayableBreakdown *CaptureSellerBreakdown `json:"seller_payable_breakdown,omitempty"`
		NoteToPayer            string                  `json:"note_to_payer,omitempty"`
		// merchant-onboarding Resource type
		PartnerClientID string `json:"partner_client_id,omitempty"`
		MerchantID      string `json:"merchant_id,omitempty"`
		// Common
		Links []Link `json:"links,omitempty"`
	}

	CaptureSellerBreakdown struct {
		GrossAmount         PurchaseUnitAmount  `json:"gross_amount"`
		PayPalFee           PurchaseUnitAmount  `json:"paypal_fee"`
		NetAmount           PurchaseUnitAmount  `json:"net_amount"`
		TotalRefundedAmount *PurchaseUnitAmount `json:"total_refunded_amount,omitempty"`
	}

	ReferralRequest struct {
		TrackingID            string                 `json:"tracking_id"`
		PartnerConfigOverride *PartnerConfigOverride `json:"partner_config_override,omitemtpy"`
		Operations            []Operation            `json:"operations,omitempty"`
		Products              []string               `json:"products,omitempty"`
		LegalConsents         []Consent              `json:"legal_consents,omitempty"`
	}

	PartnerConfigOverride struct {
		PartnerLogoURL       string `json:"partner_logo_url,omitempty"`
		ReturnURL            string `json:"return_url,omitempty"`
		ReturnURLDescription string `json:"return_url_description,omitempty"`
		ActionRenewalURL     string `json:"action_renewal_url,omitempty"`
		ShowAddCreditCard    *bool  `json:"show_add_credit_card,omitempty"`
	}

	Operation struct {
		Operation                string              `json:"operation"`
		APIIntegrationPreference *IntegrationDetails `json:"api_integration_preference,omitempty"`
	}

	IntegrationDetails struct {
		RestAPIIntegration *RestAPIIntegration `json:"rest_api_integration,omitempty"`
	}

	RestAPIIntegration struct {
		IntegrationMethod string            `json:"integration_method"`
		IntegrationType   string            `json:"integration_type"`
		ThirdPartyDetails ThirdPartyDetails `json:"third_party_details"`
	}

	ThirdPartyDetails struct {
		Features []string `json:"features"`
	}

	Consent struct {
		Type    string `json:"type"`
		Granted bool   `json:"granted"`
	}
)

// Error method implementation for ErrorResponse struct
func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %s", r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message)
}

// MarshalJSON for JSONTime
func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf(`"%s"`, time.Time(t).UTC().Format(time.RFC3339))
	return []byte(stamp), nil
}

func (e *expirationTime) UnmarshalJSON(b []byte) error {
	var n json.Number
	err := json.Unmarshal(b, &n)
	if err != nil {
		return err
	}
	i, err := n.Int64()
	if err != nil {
		return err
	}
	*e = expirationTime(i)
	return nil
}
