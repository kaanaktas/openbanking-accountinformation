package configmanager

import (
	"github.com/kaanaktas/openbanking-accountinformation/internal/cache"
	"github.com/pkg/errors"
)

type Service interface {
	FindByConfigName(aspspId, configName string) (string, error)
}

type service struct {
	ch         cache.Cache
	repository Repository
}

func NewService(repository Repository, ch cache.Cache) Service {
	return &service{
		repository: repository,
		ch:         ch,
	}
}

func (s service) FindByConfigName(aspspId, configName string) (string, error) {
	cacheId := aspspId + "_" + configName
	if value, found := s.ch.Get(cacheId); found {
		return value.(string), nil
	}

	configValue, err := s.repository.findByConfigName(aspspId, configName)
	if err != nil {
		return "", errors.WithMessage(err, "error in FindByConfigName()")
	}

	_ = s.ch.Set(cacheId, configValue, cache.NoExpiration)
	return configValue, nil
}
