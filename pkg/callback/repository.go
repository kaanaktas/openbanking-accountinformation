package callback

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Repository interface {
	saveResourceAccessAndRefreshToken(parameters map[string]interface{}) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r repository) saveResourceAccessAndRefreshToken(parameters map[string]interface{}) error {
	tx := r.db.MustBegin()
	_, err := tx.NamedExec(`UPDATE consent_token_table SET resource_access_token=:resourceAccessToken, update_date_time=:updateTime, resource_refresh_token=:resourceRefreshToken, token_status=:status,expires_in=:expiresIn, token_expiration_date_time=:tokenExpirationDateTime where consent_tid=:id`, parameters)
	if err != nil {
		tx.Rollback()
		return errors.WithMessage(err, "error in saveResourceAccessAndRefreshToken() while updating token")
	}

	_, err = tx.NamedExec(`UPDATE consent_table SET consent_status=:status, consent_status_update_date_time=:updateTime, update_date_time=:updateTime where id=:id`, parameters)
	if err != nil {
		tx.Rollback()
		return errors.WithMessage(err, "error in saveResourceAccessAndRefreshToken() while updating consent")
	}

	err = tx.Commit()
	if err != nil {
		return errors.WithMessage(err, "error in saveResourceAccessAndRefreshToken() while committing transactions")
	}

	return nil
}
