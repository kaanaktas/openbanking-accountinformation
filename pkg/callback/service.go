package callback

import (
	"fmt"
	"github.com/kaanaktas/openbanking-accountinformation/api"
	"github.com/kaanaktas/openbanking-accountinformation/internal/cache"
	"github.com/kaanaktas/openbanking-accountinformation/pkg/consent"
	"github.com/kaanaktas/openbanking-accountinformation/pkg/token"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

type Service interface {
	ProcessCallBack(code, state string) error
}

type service struct {
	repository         Repository
	consentServiceRead consent.ServiceRead
	tokenService       token.Service
	chInRedis          cache.Cache
}

func NewService(repository Repository, consentServiceRead consent.ServiceRead, tokenService token.Service, chInRedis cache.Cache) Service {
	return &service{
		repository:         repository,
		consentServiceRead: consentServiceRead,
		tokenService:       tokenService,
		chInRedis:          chInRedis,
	}
}

func (s service) ProcessCallBack(code, state string) error {
	cons, err := s.consentServiceRead.FindByTrackingId(state)
	if err != nil {
		return err
	}

	if cons.ConsentStatus == api.AwaitingAuthorisation {
		if cons.Tokens == nil {
			log.Errorf("Unsatisfied token list in Consent. referenceId: %v", code)
		}

		tokenResponse, err := s.tokenService.GetResourceAccessRefreshToken(cons.AspspId, code)
		if err != nil {
			return errors.WithMessage(err, "error in ProcessCallBack()")
		}

		tokenExpiresInSecond := tokenResponse.ExpiresIn - 300
		tokenExpirationDateTime := api.ObTime(time.Now().Add(time.Second * time.Duration(tokenExpiresInSecond)))

		updateParameters := map[string]interface{}{
			"id":                      cons.Id,
			"resourceAccessToken":     tokenResponse.AccessToken,
			"resourceRefreshToken":    tokenResponse.RefreshToken,
			"status":                  api.Authorised,
			"expiresIn":               tokenResponse.ExpiresIn,
			"tokenExpirationDateTime": tokenExpirationDateTime,
			"updateTime":              api.ObTime(time.Now()),
		}

		err = s.repository.saveResourceAccessAndRefreshToken(updateParameters)
		if err != nil {
			return errors.WithMessage(err, "error in ProcessCallBack()")
		}

		log.Info("Resource access token has been saved successfully")

		err = s.chInRedis.Set(strconv.FormatInt(cons.Id, 10), tokenResponse.AccessToken, time.Duration(tokenExpiresInSecond))
		if err == nil {
			log.Info("Resource access token has been cached successfully")
		}

		return nil
	} else {
		return fmt.Errorf("consent is not in AwaitingAuthorisation status. referenceId: %v", state)
	}
}
