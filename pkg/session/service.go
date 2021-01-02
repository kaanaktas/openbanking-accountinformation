package session

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/kaanaktas/openbanking-accountinformation/api"
	"github.com/kaanaktas/openbanking-accountinformation/internal/security"
	"github.com/pkg/errors"
	"time"
)

type Service interface {
	InitiateSession(userId, tppId, tid string) (map[string]interface{}, error)
	FindByInternalAccessToken(accessToken string) (*Session, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s service) InitiateSession(userId, tppId, tid string) (map[string]interface{}, error) {
	referenceId, internalAccessToken, err := s.retrieveSession(userId, tppId, tid)
	if err != nil {
		return nil, errors.WithMessage(err, "error in InitiateSession()")
	}

	response := map[string]interface{}{
		"internal_access_token": internalAccessToken,
		"reference_id":          referenceId,
	}

	return response, nil
}

func (s service) retrieveSession(userId, tppId, tid string) (string, string, error) {
	session, err := s.repo.findSessionByTidAndTppId(tid, tppId)
	if errors.Is(err, sql.ErrNoRows) {
		return s.createNewSession(userId, tppId, tid)
	} else {
		switch err {
		case nil:
			return session.ReferenceId, session.InternalAccessToken, nil
		default:
			return "", "", errors.WithMessage(err, "error in retrieveSession()")
		}
	}
}

func (s service) createNewSession(userId, tppId, tid string) (string, string, error) {
	claims := map[string]interface{}{
		"tppId": tppId,
		"tid":   tid,
		"iat":   security.CreateTokenTime(0),
		"exp":   security.CreateTokenTime(60),
	}

	internalAccessToken, err := security.GenerateJwtWithClaims(claims, jwt.SigningMethodHS256)
	if err != nil {
		return "", "", errors.WithMessage(err, "error in createNewSession()")
	}

	plainText := internalAccessToken + uuid.New().String()
	hasher := sha1.New()
	hasher.Write([]byte(plainText)) //nolint:errcheck
	referenceId := hex.EncodeToString(hasher.Sum(nil))

	session := &Session{
		UserId:              userId,
		Tid:                 tid,
		TppId:               tppId,
		InternalAccessToken: internalAccessToken,
		ReferenceId:         referenceId,
		CreateDateTime:      api.ObTime(time.Now()),
		UpdateDateTime:      api.ObTime(time.Now()),
	}

	err = s.repo.saveSession(session)
	if err != nil {
		return "", "", errors.WithMessage(err, "error in createNewSession()")
	}

	return referenceId, internalAccessToken, nil
}

func (s service) FindByInternalAccessToken(accessToken string) (*Session, error) {
	return s.repo.findByInternalAccessToken(accessToken)
}
