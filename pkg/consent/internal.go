package consent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/kaanaktas/openbanking-accountinformation/api"
	"github.com/kaanaktas/openbanking-accountinformation/internal/client"
	"github.com/kaanaktas/openbanking-accountinformation/internal/security"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
)

func authorizeConsentId(clientId, iss, aud, endpointAuthorize, scope, state, nonce, consentId, redirectUrl string) (string, error) {
	acr := Acr{
		Value:     "urn:openbanking:psd2:sca",
		Essential: true,
	}
	openBankingIntentId := OpenBankingIntentId{
		Value:     consentId,
		Essential: true,
	}
	idToken := IdToken{
		OpenBankingIntentId: openBankingIntentId,
		Acr:                 acr,
	}
	userInfo := Userinfo{OpenBankingIntentId: openBankingIntentId}

	claims := Claims{
		Userinfo: userInfo,
		IdToken:  idToken,
	}
	authorisedConsent := AuthorisedConsent{
		Iss:          iss,
		Aud:          aud,
		ResponseType: api.CodeIdToken,
		ClientId:     clientId,
		RedirectUri:  redirectUrl,
		Scope:        scope,
		Nonce:        nonce,
		State:        state,
		Exp:          security.CreateTokenTime(240),
		Iat:          security.CreateTokenTime(0),
		Claims:       claims,
	}

	var errMessage = "error occurred in authorizeConsentId()"
	jsonData, err := json.Marshal(authorisedConsent)
	if err != nil {
		return "", errors.WithMessage(err, errMessage)
	}

	signedJwt, err := security.GenerateJwtWithJsonString(bytes.NewBuffer(jsonData).String(), jwt.SigningMethodPS256)
	if err != nil {
		return "", errors.WithMessage(err, errMessage)
	}

	parameters := url.Values{}
	parameters.Set(api.ResponseType, api.CodeIdToken)
	parameters.Set(api.ClientIdParam, clientId)
	parameters.Set(api.RedirectUri, redirectUrl)
	parameters.Set(api.Scope, "openid "+scope)
	parameters.Set(api.Nonce, nonce)
	parameters.Set(api.State, state)
	parameters.Set(api.Request, signedJwt)

	header := http.Header{}
	header.Set(api.CacheControl, "no-cache")
	httpClient, err := client.NewSecureHttpClient(endpointAuthorize, header)
	if err != nil {
		return "", errors.WithMessage(err, errMessage)
	}

	resp, err := httpClient.Get(parameters)
	if err != nil {
		return "", errors.WithMessage(err, errMessage)
	}

	if resp.StatusCode != 302 || resp.Header.Get(api.Location) == "" {
		return "", fmt.Errorf("unexpected error while calling authorizeConsentId(). resp: %v", *resp)
	}

	return resp.Header.Get(api.Location), nil
}
