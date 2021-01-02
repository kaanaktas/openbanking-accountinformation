package session

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Repository interface {
	findSessionByTidAndTppId(tid string, tppId string) (*Session, error)
	findByInternalAccessToken(accessToken string) (*Session, error)
	saveSession(session *Session) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r repository) findSessionByTidAndTppId(tid string, tppId string) (*Session, error) {
	var session Session
	err := r.db.Get(&session, `SELECT reference_id, internal_access_token from session_table where tid=$1 AND tpp_id=$2`, tid, tppId)
	if err != nil {
		return nil, errors.WithMessage(err, "error in findSessionByTidAndTppId()")
	}

	return &session, nil
}

func (r repository) findByInternalAccessToken(accessToken string) (*Session, error) {
	var session Session
	err := r.db.Get(&session, `SELECT * from session_table where internal_access_token=$1`, accessToken)
	if err != nil {
		return nil, errors.WithMessage(err, "error in findByInternalAccessToken()")
	}

	return &session, nil
}

func (r repository) saveSession(session *Session) error {
	_, err := r.db.NamedExec(`INSERT INTO session_table(tid, tpp_id, user_id, aspsp_id, internal_access_token, reference_id, create_date_time, update_date_time) 
												VALUES (:tid, :tpp_id, :user_id, :aspsp_id, :internal_access_token, :reference_id, :create_date_time, :update_date_time)`, session)
	if err != nil {
		return errors.WithMessage(err, "error in saveSession()")
	}

	return nil
}
