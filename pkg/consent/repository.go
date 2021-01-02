package consent

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/kaanaktas/openbanking-accountinformation/api"
	"github.com/pkg/errors"
	"time"
)

type RepositoryRead interface {
	findByTrackingId(trackingId string) (*Consent, error)
	findConsentByCidAndStatus(cid, status string) (*Consent, error)
	findConsentByUserIdAndTppIdAndStatus(userId, tppId, status string) ([]Consent, error)
	findByCid(cid string) (*Consent, error)
}

type RepositoryWrite interface {
	saveConsent(consent *Consent) error
	saveToken(token *Token) error
	invalidateAuthorisedTokenByConsentTid(tid int64, status string) error
	changeConsentStateByCid(cid, status string) error
}

type repositoryRead struct {
	db *sqlx.DB
}

type repositoryWrite struct {
	db *sqlx.DB
}

func NewRepositoryRead(db *sqlx.DB) RepositoryRead {
	return &repositoryRead{db: db}
}

func NewRepositoryWrite(db *sqlx.DB) RepositoryWrite {
	return &repositoryWrite{db: db}
}

func (r repositoryRead) findByTrackingId(trackingId string) (*Consent, error) {
	var consentToken []TokensInConsent

	err := r.db.Select(&consentToken, `SELECT ctt.id as token_tid, ctt.*, ct.* from consent_table ct INNER JOIN consent_token_table ctt on ct.id = ctt.consent_tid WHERE ct.tracking_id = $1`, trackingId)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, errors.WithMessage(err, "error in findByTrackingId()")
	}
	if consentToken == nil {
		return nil, sql.ErrNoRows
	}

	var tokens = make([]Token, len(consentToken))
	for k, v := range consentToken {
		tokens[k] = v.Token
	}
	var consent = consentToken[0].Consent
	consent.Tokens = tokens

	return &consent, nil
}

func (r repositoryRead) findByCid(cid string) (*Consent, error) {
	var consent Consent
	err := r.db.Get(&consent, `SELECT ct.* FROM consent_table ct WHERE ct.id = $1`, cid)
	if err != nil {
		return nil, err
	}

	return &consent, nil
}

func (r repositoryRead) findConsentByCidAndStatus(cid, status string) (*Consent, error) {
	var consentToken []TokensInConsent
	err := r.db.Select(&consentToken, `SELECT ct.*, ctt.id as token_tid, ctt.* from consent_table ct LEFT JOIN consent_token_table ctt on ct.id = ctt.consent_tid AND ctt.token_status = 'Authorised' WHERE ct.consent_status = $1 AND ct.id = $2`, status, cid)
	if err != nil {
		return nil, errors.WithMessage(err, "error in findConsentByCidAndStatus()")
	}
	if consentToken == nil {
		return nil, sql.ErrNoRows
	}

	var tokens []Token
	for _, v := range consentToken {
		if v.Token.Id != nil {
			tokens = append(tokens, v.Token)
		}
	}
	var consent = consentToken[0].Consent
	if len(tokens) > 0 {
		consent.Tokens = tokens
	}

	return &consent, nil
}

func (r repositoryRead) findConsentByUserIdAndTppIdAndStatus(userId, tppId, status string) ([]Consent, error) {
	var consents []Consent
	err := r.db.Select(&consents, `SELECT ct.id as id, ct.* FROM consent_table ct INNER JOIN session_table st on ct.session_reference_id = st.reference_id
													WHERE st.user_id = $1 AND st.tpp_id = $2 AND ct.consent_status = $3`, userId, tppId, status)
	if err != nil {
		return nil, errors.WithMessage(err, "error in findConsentByUserIdAndTppIdAndStatus()")
	}

	return consents, nil
}

func (r repositoryWrite) saveConsent(consent *Consent) error {
	tx := r.db.MustBegin()

	stmt, err := tx.PrepareNamed(`INSERT INTO consent_table(aspsp_id, consent_expiration_date_time, consent_id, consent_status, consent_status_update_date_time, 
                          consent_transaction_from_date_time, consent_transaction_to_date_time, consent_type, create_date_time, session_reference_id, tracking_id, update_date_time) 
                          VALUES (:aspsp_id, :consent_expiration_date_time, :consent_id, :consent_status, :consent_status_update_date_time, :consent_transaction_from_date_time, 
                                  :consent_transaction_to_date_time, :consent_type, :create_date_time, :session_reference_id, :tracking_id, :update_date_time) RETURNING id`)
	if err != nil {
		tx.Rollback()
		return errors.WithMessage(err, "error in saveConsent() while inserting the new consent")
	}

	var lastInsertedId int64
	err = stmt.Get(&lastInsertedId, consent)
	if err != nil {
		tx.Rollback()
		return errors.WithMessage(err, "error in saveConsent() while trying to retrieve lastInsertedId")
	}

	tokens := consent.Tokens

	for _, token := range tokens {
		token.ConsentTid = &lastInsertedId
		_, err = tx.NamedExec(`INSERT INTO consent_token_table(access_token, create_date_time, expires_in, resource_access_token, resource_refresh_token, token_expiration_date_time, 
                                token_status, update_date_time, consent_tid) VALUES (:access_token, :create_date_time, :expires_in, :resource_access_token, 
                                                                                     :resource_refresh_token, :token_expiration_date_time, :token_status, :update_date_time, :consent_tid)`, token)

		if err != nil {
			tx.Rollback()
			return errors.WithMessage(err, "error in saveConsent() while inserting token details")
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.WithMessage(err, "error in saveConsent() while committing transactions")
	}

	return nil
}

func (r repositoryWrite) saveToken(token *Token) error {
	_, err := r.db.NamedExec("INSERT INTO consent_token_table(access_token, create_date_time, expires_in, resource_access_token, resource_refresh_token, "+
		"token_expiration_date_time, token_status, update_date_time, consent_tid) "+
		"VALUES (:access_token, :create_date_time, :expires_in, :resource_access_token, :resource_refresh_token, "+
		":token_expiration_date_time, :token_status, :update_date_time, :consent_tid)", token)
	if err != nil {
		return errors.WithMessage(err, "error in saveToken() while inserting the new token")
	}

	return nil
}

func (r repositoryWrite) invalidateAuthorisedTokenByConsentTid(tid int64, status string) error {
	_, err := r.db.NamedExec(`UPDATE consent_token_table SET token_status=:status WHERE token_status = 'Authorised' AND consent_tid=:tid`, map[string]interface{}{"status": status, "tid": tid})
	if err != nil {
		return errors.WithMessage(err, "error in invalidateAuthorisedTokenByConsentTid()")
	}

	return nil
}

func (r repositoryWrite) changeConsentStateByCid(cid, status string) error {
	updateTime := api.ObTime(time.Now())
	_, err := r.db.NamedExec(`UPDATE consent_table SET consent_status=:status, update_date_time=:updateTime WHERE id=:cid`,
		map[string]interface{}{"status": status, "cid": cid, "updateTime": updateTime})

	if err != nil {
		return errors.WithMessage(err, "error in changeConsentStateByCid()")
	}

	return nil
}
