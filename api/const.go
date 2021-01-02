package api

//Media type constants.
const (
	ApplicationFormUrlencodedValue = "application/x-www-form-urlencoded"
	ApplicationJson                = "application/json"
)

//Config parameters constants.
const (
	EndpointOauth2               = "ENDPOINT_OAUTH2"
	EndpointDomesticConsent      = "ENDPOINT_DOMESTIC_CONSENT"
	EndpointAuthorize            = "ENDPOINT_AUTHORIZE"
	EndpointDomesticPayment      = "ENDPOINT_DOMESTIC_PAYMENT"
	EndpointAccountAccessConsent = "ENDPOINT_ACCOUNT_ACCESS_CONSENT"
	EndpointRegister             = "ENDPOINT_REGISTER"
	EndpointAccounts             = "ENDPOINT_ACCOUNTS"
	FapiFinancialId              = "FAPI_FINANCIAL_ID"
	ClientId                     = "CLIENT_ID"
	Iss                          = "ISS"
	TokenEndpointAuthMethod      = "TOKEN_ENDPOINT_AUTH_METHOD"
	Aud                          = "AUD"
	ApplicationType              = "APPLICATION_TYPE"
	RedirectUrl                  = "REDIRECT_URL"
	AuthenticationToken          = "AUTHENTICATION_TOKEN"
)

//Http header constants.
const (
	Accept                        = "Accept"
	AcceptCharset                 = "Accept-Charset"
	AcceptEncoding                = "Accept-Encoding"
	AcceptLanguage                = "Accept-Language"
	Authorization                 = "Authorization"
	CacheControl                  = "Cache-Control"
	ContentLength                 = "Content-Length"
	ContentMD5                    = "Content-MD5"
	ContentType                   = "Content-Type"
	DoNotTrack                    = "DNT"
	IfMatch                       = "If-Match"
	IfModifiedSince               = "If-Modified-Since"
	IfNoneMatch                   = "If-None-Match"
	IfRange                       = "If-Range"
	IfUnmodifiedSince             = "If-Unmodified-Since"
	MaxForwards                   = "Max-Forwards"
	ProxyAuthorization            = "Proxy-Authorization"
	Pragma                        = "Pragma"
	Range                         = "Range"
	Referer                       = "Referer"
	UserAgent                     = "User-Agent"
	TE                            = "TE"
	Via                           = "Via"
	Warning                       = "Warning"
	Cookie                        = "Cookie"
	Origin                        = "Origin"
	AcceptDatetime                = "Accept-Datetime"
	XRequestedWith                = "X-Requested-With"
	AccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	AccessControlAllowMethods     = "Access-Control-Allow-Methods"
	AccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	AccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	AccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	AccessControlMaxAge           = "Access-Control-Max-Age"
	AccessControlRequestMethod    = "Access-Control-Request-Method"
	AccessControlRequestHeaders   = "Access-Control-Request-Headers"
	AcceptPatch                   = "Accept-Patch"
	AcceptRanges                  = "Accept-Ranges"
	Allow                         = "Allow"
	ContentEncoding               = "Content-Encoding"
	ContentLanguage               = "Content-Language"
	ContentLocation               = "Content-Location"
	ContentDisposition            = "Content-Disposition"
	ContentRange                  = "Content-Range"
	ETag                          = "ETag"
	Expires                       = "Expires"
	LastModified                  = "Last-Modified"
	Link                          = "Link"
	Location                      = "Location"
	P3P                           = "P3P"
	ProxyAuthenticate             = "Proxy-Authenticate"
	Refresh                       = "Refresh"
	RetryAfter                    = "Retry-After"
	Server                        = "Server"
	SetCookie                     = "Set-Cookie"
	StrictTransportSecurity       = "Strict-Transport-Security"
	TransferEncoding              = "Transfer-Encoding"
	Upgrade                       = "Upgrade"
	Vary                          = "Vary"
	WWWAuthenticate               = "WWW-Authenticate"

	// Non-Standard
	XFapiFinancialId       = "x-fapi-financial-id"
	XIdempotencyKey        = "x-idempotency-key"
	XFrameOptions          = "X-Frame-Options"
	XXSSProtection         = "X-XSS-Protection"
	ContentSecurityPolicy  = "Content-Security-Policy"
	XContentSecurityPolicy = "X-Content-Security-Policy"
	XWebKitCSP             = "X-WebKit-CSP"
	XContentTypeOptions    = "X-Content-Type-Options"
	XPoweredBy             = "X-Powered-By"
	XUACompatible          = "X-UA-Compatible"
	XForwardedProto        = "X-Forwarded-Proto"
	XHTTPMethodOverride    = "X-HTTP-Method-Override"
	XForwardedFor          = "X-Forwarded-For"
	XRealIP                = "X-Real-IP"
	XCSRFToken             = "X-CSRF-Token"
	XRatelimitLimit        = "X-Ratelimit-Limit"
	XRatelimitRemaining    = "X-Ratelimit-Remaining"
	XRatelimitReset        = "X-Ratelimit-Reset"
)

//open banking constants.
const (
	ScopeAccounts = "openid accounts"
	ScopePayments = "payments"
	ResponseType  = "response_type"
	ClientIdParam = "client_id"
	RedirectUri   = "redirect_uri"
	Request       = "request"
	Scope         = "scope"
	Nonce         = "nonce"
	State         = "state"
)

//outh2 constants
const (
	CodeIdToken  = "code id_token"
	CodeToken    = "code token"
	IdToken      = "id_token"
	IdTokenToken = "id_token token"
)

//Consent types constants.
const (
	Authorised            = "Authorised"
	AwaitingAuthorisation = "AwaitingAuthorisation"
	Rejected              = "Rejected"
	Consumed              = "Consumed"
	Revoked               = "Revoked"
	Expired               = "Expired"

	Account         = "ACCOUNT"
	DomesticPayment = "DOMESTIC_PAYMENT"
)

//cache key id constants.
const (
	InternalSignKey = "internal_sign_key"
	ObSignKey       = "ob_sign_key"
)
