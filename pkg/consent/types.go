package consent

import (
	"time"
)

type ObReadConsent struct {
	Data *ObReadConsentData `json:"Data"`
	Risk *ObRisk            `json:"Risk"`
}

type ObReadConsentData struct {
	Permissions []string `json:"Permissions"`
	// Specified date and time the permissions will expire.
	//If this is not populated, the permissions will be open ended.All dates in the JSON payloads are represented in ISO 8601 date-time format.
	//All date-time fields in responses must include the timezone. An example is below: 2017-04-05T10:43:07+00:00
	ExpirationDateTime string `json:"ExpirationDateTime,omitempty"`
	// Specified start date and time for the transaction query period.
	//If this is not populated, the start date will be open ended, and data will be returned from the earliest available transaction.
	//All dates in the JSON payloads are represented in ISO 8601 date-time format. All date-time fields in responses must include the timezone.
	//An example is below: 2017-04-05T10:43:07+00:00
	TransactionFromDateTime string `json:"TransactionFromDateTime,omitempty"`
	// Specified end date and time for the transaction query period.
	//If this is not populated, the end date will be open ended, and data will be returned to the latest available transaction.
	//All dates in the JSON payloads are represented in ISO 8601 date-time format. All date-time fields in responses must include the timezone.
	//An example is below: 2017-04-05T10:43:07+00:00
	TransactionToDateTime string `json:"TransactionToDateTime,omitempty"`
}

// The Risk section is sent by the initiating party to the ASPSP. It is used to specify additional details for risk scoring for Account Info.
type ObRisk struct {
}

type ObReadConsentResponse struct {
	Data  *ObReadConsentResponseData `json:"Data"`
	Risk  *ObRisk                    `json:"Risk"`
	Links *Links                     `json:"Links,omitempty"`
	Meta  *Meta                      `json:"Meta,omitempty"`
}

type ObReadConsentResponseData struct {
	// Unique identification as assigned to identify the account access consent resource.
	ConsentId        string `json:"ConsentId"`
	CreationDateTime string `json:"CreationDateTime"`
	// Specifies the status of consent resource in code form.
	Status               string   `json:"Status"`
	StatusUpdateDateTime string   `json:"StatusUpdateDateTime"`
	Permissions          []string `json:"Permissions"`
	// Specified date and time the permissions will expire. If this is not populated, the permissions will be open ended.All dates in the JSON payloads are represented in ISO 8601 date-time format.  All date-time fields in responses must include the timezone. An example is below: 2017-04-05T10:43:07+00:00
	ExpirationDateTime time.Time `json:"ExpirationDateTime,omitempty"`
	// Specified start date and time for the transaction query period. If this is not populated, the start date will be open ended, and data will be returned from the earliest available transaction.All dates in the JSON payloads are represented in ISO 8601 date-time format.  All date-time fields in responses must include the timezone. An example is below: 2017-04-05T10:43:07+00:00
	TransactionFromDateTime time.Time `json:"TransactionFromDateTime,omitempty"`
	// Specified end date and time for the transaction query period. If this is not populated, the end date will be open ended, and data will be returned to the latest available transaction.All dates in the JSON payloads are represented in ISO 8601 date-time format.  All date-time fields in responses must include the timezone. An example is below: 2017-04-05T10:43:07+00:00
	TransactionToDateTime time.Time `json:"TransactionToDateTime,omitempty"`
}

// Links relevant to the payload
type Links struct {
	Self  string `json:"Self"`
	First string `json:"First,omitempty"`
	Prev  string `json:"Prev,omitempty"`
	Next  string `json:"Next,omitempty"`
	Last  string `json:"Last,omitempty"`
}

// Meta Data relevant to the payload
type Meta struct {
	TotalPages             int32  `json:"TotalPages,omitempty"`
	FirstAvailableDateTime string `json:"FirstAvailableDateTime,omitempty"`
	LastAvailableDateTime  string `json:"LastAvailableDateTime,omitempty"`
}

// ******** TYPES OF AUTHORIZE CONSENT JWT *************
type Acr struct {
	Value     string `json:"value"`
	Essential bool   `json:"essential"`
}

type OpenBankingIntentId struct {
	Value     string `json:"value"`
	Essential bool   `json:"essential"`
}

type IdToken struct {
	OpenBankingIntentId OpenBankingIntentId `json:"openbanking_intent_id"`
	Acr                 Acr                 `json:"acr"`
}

type Userinfo struct {
	OpenBankingIntentId OpenBankingIntentId `json:"openbanking_intent_id"`
}

type Claims struct {
	Userinfo Userinfo `json:"userinfo"`
	IdToken  IdToken  `json:"id_token"`
}

type AuthorisedConsent struct {
	Iss          string `json:"iss"`
	Aud          string `json:"aud"`
	ResponseType string `json:"response_type"`
	ClientId     string `json:"client_id"`
	RedirectUri  string `json:"redirect_uri"`
	Scope        string `json:"scope"`
	Nonce        string `json:"nonce"`
	State        string `json:"state"`
	Exp          int64  `json:"exp"`
	Iat          int64  `json:"iat"`
	Claims       Claims `json:"claims"`
}

type ActiveConsent struct {
	ConsentTid int64  `json:"consentTid"`
	AspspId    string `json:"aspspId"`
}

type Token struct {
	Id                      *int64  `db:"token_tid"`
	AccessToken             *string `db:"access_token"`
	ResourceAccessToken     *string `db:"resource_access_token"`
	ResourceRefreshToken    *string `db:"resource_refresh_token"`
	TokenStatus             *string `db:"token_status"`
	ExpiresIn               *int    `db:"expires_in"`
	CreateDateTime          *string `db:"create_date_time"`
	UpdateDateTime          *string `db:"update_date_time"`
	TokenExpirationDateTime *string `db:"token_expiration_date_time"`
	ConsentTid              *int64  `db:"consent_tid"`
}

type Consent struct {
	Id                             int64  `db:"id"`
	SessionReferenceId             string `db:"session_reference_id"`
	TrackingId                     string `db:"tracking_id"`
	AspspId                        string `db:"aspsp_id"`
	ConsentId                      string `db:"consent_id"`
	ConsentExpirationDateTime      string `db:"consent_expiration_date_time"`
	ConsentTransactionFromDateTime string `db:"consent_transaction_from_date_time"`
	ConsentTransactionToDateTime   string `db:"consent_transaction_to_date_time"`
	ConsentStatusUpdateDateTime    string `db:"consent_status_update_date_time"`
	CreateDateTime                 string `db:"create_date_time"`
	UpdateDateTime                 string `db:"update_date_time"`
	ConsentStatus                  string `db:"consent_status"`
	ConsentType                    string `db:"consent_type"`
	Tokens                         []Token
}

type TokensInConsent struct {
	Consent
	Token
}
