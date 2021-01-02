package configmanager

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Repository interface {
	findByConfigName(aspspId, configName string) (string, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r repository) findByConfigName(aspspId, configName string) (string, error) {
	var configValue string
	err := r.db.Get(&configValue, "SELECT config_value from config_table where config_name=$1 and aspsp_id=$2", configName, aspspId)
	if err != nil {
		return "", errors.WithMessage(err, "error in findByConfigName()")
	}

	return configValue, nil
}
