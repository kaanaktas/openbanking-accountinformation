package accounts

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/kaanaktas/openbanking-accountinformation/api"
	"github.com/kaanaktas/openbanking-accountinformation/internal/client"
	"github.com/kaanaktas/openbanking-accountinformation/pkg/authmanager"
	cfg "github.com/kaanaktas/openbanking-accountinformation/pkg/configmanager"
	"github.com/pkg/errors"
	"net/http"
)

type Service interface {
	Account(cid, aspspId, accountId string) (string, error)
	Accounts(cid, aspspId string) (string, error)
}

type service struct {
	authManager authmanager.AuthManager
	cfg         cfg.Service
}

func NewService(authManager authmanager.AuthManager, cfg cfg.Service) Service {
	return &service{
		authManager: authManager,
		cfg:         cfg,
	}
}

func (s service) Account(cid, aspspId, accountId string) (string, error) {
	endpointAccounts, err := s.cfg.FindByConfigName(aspspId, api.EndpointAccounts)
	if err != nil {
		return "", err
	}
	endpointAccounts = endpointAccounts + "/" + accountId

	return s.processCall(cid, aspspId, endpointAccounts)
}

func (s service) Accounts(cid, aspspId string) (string, error) {
	endpointAccounts, err := s.cfg.FindByConfigName(aspspId, api.EndpointAccounts)
	if err != nil {
		return "", errors.WithMessage(err, "error in Accounts()")
	}

	return s.processCall(cid, aspspId, endpointAccounts)
}

func (s service) processCall(cid, aspspId, endpointAccounts string) (string, error) {
	resourceAccessToken, err := s.authManager.GetAuthorisedTokenByCid(aspspId, cid)
	if err != nil {
		return "", err
	}

	fapiFinancialId, err := s.cfg.FindByConfigName(aspspId, api.FapiFinancialId)
	if err != nil {
		return "", err
	}

	httpClient, err := client.NewSecureHttpClient(endpointAccounts, s.setHeader(resourceAccessToken, fapiFinancialId))
	if err != nil {
		return "", errors.WithMessage(err, "error in processCall()")
	}

	resp, err := httpClient.Get(nil)
	if err != nil {
		return "", errors.WithMessage(err, "error in processCall()")
	}

	switch resp.StatusCode {
	case 200, 201:
		return resp.Body, nil
	//TODO case 401 is an anomaly and needs to be taken care either by revoking the consent or refreshing the token
	default:
		return "", fmt.Errorf("unexpected result from the accounts service. resp: %v", *resp)
	}
}

func (s service) setHeader(resourceAccessToken, fapiFinancialId string) http.Header {
	header := http.Header{}
	header.Set(api.Accept, api.ApplicationJson)
	header.Set(api.Authorization, "Bearer "+resourceAccessToken)
	header.Set(api.ContentType, api.ApplicationJson)
	header.Set(api.CacheControl, "no-cache")
	header.Set(api.XIdempotencyKey, uuid.New().String())
	header.Set(api.XFapiFinancialId, fapiFinancialId)

	return header
}
