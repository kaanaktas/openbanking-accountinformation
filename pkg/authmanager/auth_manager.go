package authmanager

import (
	"fmt"
	"github.com/kaanaktas/openbanking-accountinformation/api"
	"github.com/kaanaktas/openbanking-accountinformation/internal/cache"
	"github.com/kaanaktas/openbanking-accountinformation/pkg/consent"
	"github.com/kaanaktas/openbanking-accountinformation/pkg/token"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"time"
)

type AuthManager interface {
	GetAuthorisedTokenByCid(aspspId, cid string) (string, error)
}

type authManager struct {
	consentServiceRead  consent.ServiceRead
	consentServiceWrite consent.ServiceWrite
	tokenService        token.Service
	chRedis             cache.Cache
}

func NewAuthManager(consentServiceRead consent.ServiceRead, consentServiceWrite consent.ServiceWrite, tokenService token.Service, chRedis cache.Cache) AuthManager {
	return &authManager{
		consentServiceRead:  consentServiceRead,
		consentServiceWrite: consentServiceWrite,
		tokenService:        tokenService,
		chRedis:             chRedis,
	}
}

func (s authManager) GetAuthorisedTokenByCid(aspspId, cid string) (string, error) {
	if value, found := s.chRedis.Get(cid); found {
		return value.(string), nil
	}

	consentResp, err := s.consentServiceRead.FindConsentByCidAndStatus(cid, api.Authorised)
	if err != nil || consentResp.AspspId != aspspId {
		return "", errors.WithMessagef(err, "couldn't retrieve the consentResp. cid: %v aspspId: %v", cid, aspspId)
	}

	var consentExpirationDateTime time.Time
	consentExpirationDateTime, err = time.Parse(time.RFC3339, consentResp.ConsentExpirationDateTime)
	if err != nil {
		return "", errors.WithMessage(err, "error in GetAuthorisedTokenByCid()")
	}

	//if no authorised authorisedToken, revoke consentResp
	if consentExpirationDateTime.Before(time.Now()) || consentResp.Tokens == nil || len(consentResp.Tokens) < 1 {
		err := s.consentServiceWrite.ChangeConsentStateByCid(cid, api.Revoked)
		if err != nil {
			log.Errorf("unexpected error while revoking the consentResp for cid: %v. consentResp revoking will be tried with the new request. err: %v", cid, err)
		}

		return "", fmt.Errorf("consentResp expired or doesn't have authorised token, it has been revoked. cid: %v", cid)
	}

	authorisedToken := consentResp.Tokens[0]
	var tokenExpirationDateTime time.Time
	tokenExpirationDateTime, err = time.Parse(time.RFC3339, *authorisedToken.TokenExpirationDateTime)
	if err != nil {
		return "", errors.WithMessage(err, "error in GetAuthorisedTokenByCid()")
	}

	log.Infof("Checking Resource Token expiry time. resourceAccessToken: %v, consentId: %v", *authorisedToken.ResourceAccessToken, consentResp.ConsentId)
	timeNow := time.Now()
	if tokenExpirationDateTime.After(timeNow) {
		//token is valid. cache it again then return the resource access token
		tokenExpiresInSecond := tokenExpirationDateTime.Sub(timeNow).Seconds()
		_ = s.chRedis.Set(cid, *authorisedToken.ResourceAccessToken, time.Duration(tokenExpiresInSecond))
		return *authorisedToken.ResourceAccessToken, nil
	} else {
		log.Infof("Resource token has been expired. resourceAccessToken: %v. Requesting a new resource token for the existing consentResp. resourceRefreshToken: %v",
			*authorisedToken.ResourceAccessToken, *authorisedToken.ResourceRefreshToken)

		refreshToken := *authorisedToken.ResourceRefreshToken
		//authorisedToken expired call refresh authorisedToken
		tokenResp, err := s.tokenService.RefreshAccessToken(aspspId, api.ScopeAccounts, refreshToken)
		if err != nil {
			return "", errors.WithMessage(err, "error in GetAuthorisedTokenByCid()")
		}

		tokenExpiresInSecond := tokenResp.ExpiresIn - 300
		tokenExpirationDateTime := api.ObTime(time.Now().Add(time.Second * time.Duration(tokenExpiresInSecond)))

		tokenStatus := api.Authorised
		dateTime := api.ObTime(time.Now())
		newToken := &consent.Token{
			AccessToken:             authorisedToken.AccessToken,
			ResourceAccessToken:     &tokenResp.AccessToken,
			ResourceRefreshToken:    &tokenResp.RefreshToken,
			TokenStatus:             &tokenStatus,
			ExpiresIn:               &tokenExpiresInSecond,
			CreateDateTime:          &dateTime,
			UpdateDateTime:          &dateTime,
			TokenExpirationDateTime: &tokenExpirationDateTime,
			ConsentTid:              authorisedToken.ConsentTid,
		}

		err = s.consentServiceWrite.InvalidateAuthorisedTokenByConsentTid(*authorisedToken.ConsentTid, api.Expired)
		if err != nil {
			return "", errors.WithMessage(err, "error in GetAuthorisedTokenByCid()")
		}

		log.Info("All authorised tokens set to EXPIRED. consentTid:", *authorisedToken.ConsentTid)
		err = s.consentServiceWrite.SaveToken(newToken)
		if err != nil {
			return "", errors.WithMessage(err, "error in GetAuthorisedTokenByCid()")
		}

		log.Info("Resource access token refreshed successfully. refreshAccessToken:", tokenResp.AccessToken)

		err = s.chRedis.Set(cid, tokenResp.AccessToken, time.Duration(tokenExpiresInSecond))
		if err == nil {
			log.Info("Resource access token cached successfully")
		} else {
			log.Error("Resource access token couldn't be cached. this will be tried again with the next request")
		}

		return tokenResp.AccessToken, nil
	}
}
