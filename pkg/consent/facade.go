package consent

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/kaanaktas/openbanking-accountinformation/api"
	"github.com/kaanaktas/openbanking-accountinformation/internal/client"
	cfg "github.com/kaanaktas/openbanking-accountinformation/pkg/configmanager"
	"github.com/kaanaktas/openbanking-accountinformation/pkg/token"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type Facade interface {
	CreateConsent(sessionReferenceId, trackingId, aspspId string, consent *ObReadConsent) (string, error)
	GetConsent(cid, aspspId string) (string, error)
	DeleteConsent(consentId string) (string, error)
}

type facade struct {
	serviceRead  ServiceRead
	serviceWrite ServiceWrite
	tokenService token.Service
	cfg          cfg.Service
}

func NewFacade(serviceRead ServiceRead, serviceWrite ServiceWrite, tokenService token.Service, cfg cfg.Service) Facade {
	return &facade{
		serviceRead:  serviceRead,
		serviceWrite: serviceWrite,
		tokenService: tokenService,
		cfg:          cfg,
	}
}

func (f facade) CreateConsent(sessionReferenceId, trackingId, aspspId string, obConsent *ObReadConsent) (string, error) {
	consentResp, err := f.serviceRead.FindByTrackingId(trackingId)
	if (err != nil && !errors.Is(err, sql.ErrNoRows)) || consentResp != nil {
		return "", fmt.Errorf("reference has already been used. please try a new one")
	}

	var errMessage = "error in CreateConsent()"
	obAccessToken, err := f.tokenService.GetAccessToken(aspspId, api.ScopeAccounts)
	if err != nil {
		return "", errors.WithMessage(err, errMessage)
	}

	consentJson, err := json.Marshal(obConsent)
	if err != nil {
		return "", errors.WithMessage(err, errMessage)
	}

	endpointAccountAccessConsent, _ := f.cfg.FindByConfigName(aspspId, api.EndpointAccountAccessConsent)
	fapiFinancialId, _ := f.cfg.FindByConfigName(aspspId, api.FapiFinancialId)

	httpClient, err := client.NewSecureHttpClient(endpointAccountAccessConsent, f.setHeader(obAccessToken, fapiFinancialId))
	if err != nil {
		return "", errors.WithMessage(err, errMessage)
	}

	resp, err := httpClient.Post(bytes.NewBuffer(consentJson))
	if err != nil {
		return "", errors.WithMessage(err, errMessage)
	}

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		var responseData ObReadConsentResponse
		err := json.Unmarshal([]byte(resp.Body), &responseData)
		if err != nil {
			return "", errors.WithMessage(err, errMessage)
		}

		consentId := responseData.Data.ConsentId
		createDateTime := api.ObTime(time.Now())
		tokenStatus := api.AwaitingAuthorisation
		expiresIn := 0
		tokens := []Token{
			{
				AccessToken:    &obAccessToken,
				TokenStatus:    &tokenStatus,
				CreateDateTime: &createDateTime,
				ExpiresIn:      &expiresIn,
			},
		}

		consentDetail := &Consent{
			TrackingId:                     trackingId,
			SessionReferenceId:             sessionReferenceId,
			AspspId:                        aspspId,
			ConsentId:                      consentId,
			ConsentExpirationDateTime:      obConsent.Data.ExpirationDateTime,
			ConsentTransactionFromDateTime: obConsent.Data.TransactionFromDateTime,
			ConsentTransactionToDateTime:   obConsent.Data.TransactionToDateTime,
			ConsentStatusUpdateDateTime:    api.ObTime(time.Now()),
			CreateDateTime:                 api.ObTime(time.Now()),
			UpdateDateTime:                 api.ObTime(time.Now()),
			ConsentStatus:                  api.AwaitingAuthorisation,
			ConsentType:                    api.Account,
			Tokens:                         tokens,
		}

		err = f.serviceWrite.SaveConsent(consentDetail)
		if err != nil {
			return "", errors.WithMessage(err, errMessage)
		}

		redirectUrl, err := f.cfg.FindByConfigName(aspspId, api.RedirectUrl)
		if err != nil {
			return "", errors.WithMessage(err, errMessage)
		}
		clientId, err := f.cfg.FindByConfigName(aspspId, api.ClientId)
		if err != nil {
			return "", errors.WithMessage(err, errMessage)
		}
		iss, err := f.cfg.FindByConfigName(aspspId, api.Iss)
		if err != nil {
			return "", errors.WithMessage(err, errMessage)
		}
		aud, err := f.cfg.FindByConfigName(aspspId, api.Aud)
		if err != nil {
			return "", errors.WithMessage(err, errMessage)
		}
		endpointAuthorize, err := f.cfg.FindByConfigName(aspspId, api.EndpointAuthorize)
		if err != nil {
			return "", errors.WithMessage(err, errMessage)
		}

		authRedirectUrl, err := authorizeConsentId(clientId, iss, aud, endpointAuthorize, api.ScopeAccounts, trackingId, uuid.New().String(), consentId, redirectUrl)
		if err != nil {
			return "", errors.WithMessage(err, errMessage)
		}

		return authRedirectUrl, nil
	} else {
		return "", fmt.Errorf("unexpected result from the token facade. response: %v", *resp)
	}
}

func (f facade) GetConsent(cid, aspspId string) (string, error) {
	consentResp, err := f.serviceRead.FindByCid(cid)
	if err == sql.ErrNoRows {
		return "", err
	}

	var errMessage = "error in GetConsent()"
	if err != nil {
		return "", errors.WithMessage(err, errMessage)
	}
	consentId := consentResp.ConsentId

	accessToken, err := f.tokenService.GetAccessToken(aspspId, api.ScopeAccounts)
	if err != nil {
		return "", errors.WithMessage(err, errMessage)
	}

	endpointAccountAccessConsent, _ := f.cfg.FindByConfigName(aspspId, api.EndpointAccountAccessConsent)
	fapiFinancialId, _ := f.cfg.FindByConfigName(aspspId, api.FapiFinancialId)

	httpClient, err := client.NewSecureHttpClient(endpointAccountAccessConsent+"/"+consentId,
		f.setHeader(accessToken, fapiFinancialId))
	if err != nil {
		return "", errors.WithMessage(err, errMessage)
	}

	resp, err := httpClient.Get(nil)
	if err != nil {
		return "", errors.WithMessage(err, errMessage)
	}

	if resp.StatusCode == 200 {
		return resp.Body, err
	} else {
		return "", fmt.Errorf("unexpected result from the token facade. response: %v", *resp)
	}
}

func (f facade) DeleteConsent(consentId string) (string, error) {
	panic("implement me")
}

func (f facade) setHeader(obAccessToken, xFapiFinancialId string) http.Header {
	header := http.Header{}
	header.Set(api.Authorization, "Bearer "+obAccessToken)
	header.Set(api.Accept, api.ApplicationJson)
	header.Set(api.ContentType, api.ApplicationJson)
	header.Set(api.CacheControl, "no-cache")
	header.Set(api.XFapiFinancialId, xFapiFinancialId)
	header.Set(api.XIdempotencyKey, uuid.New().String())

	return header
}
