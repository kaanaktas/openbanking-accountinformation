package consent

import (
	"github.com/jmoiron/sqlx"
	"github.com/kaanaktas/openbanking-accountinformation/api"
	"github.com/kaanaktas/openbanking-accountinformation/internal/store"
	"testing"
	"time"
)

func Test_repository_SaveConsent(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		consent *Consent
	}

	accessToken := "333333333333333333333333"
	resourceAccessToken := "444444444444444444444444"
	resourceRefreshToken := "555555555555555555555555"
	tokenStatus := api.AwaitingAuthorisation
	expiresIn := 60
	createDateTime := api.ObTime(time.Now())
	updateDateTime := api.ObTime(time.Now())
	tokenExpirationDateTime := api.ObTime(time.Now().Add(time.Hour * 2))
	consentTid := int64(0)

	token1 := Token{}
	token1.AccessToken = &accessToken
	token1.ResourceAccessToken = &resourceAccessToken
	token1.ResourceRefreshToken = &resourceRefreshToken
	token1.ExpiresIn = &expiresIn
	token1.ConsentTid = &consentTid
	token1.CreateDateTime = &createDateTime
	token1.TokenStatus = &tokenStatus
	token1.UpdateDateTime = &updateDateTime
	token1.TokenExpirationDateTime = &tokenExpirationDateTime

	accessToken = "666666666666666666666666"
	resourceAccessToken = "777777777777777777777777"
	resourceRefreshToken = "888888888888888888888888"
	tokenStatus = api.Consumed
	expiresIn = 45
	createDateTime = api.ObTime(time.Now())
	updateDateTime = api.ObTime(time.Now())
	tokenExpirationDateTime = api.ObTime(time.Now().Add(time.Hour * 2))

	token2 := Token{}
	token2.AccessToken = &accessToken
	token2.ResourceAccessToken = &resourceAccessToken
	token2.ResourceRefreshToken = &resourceRefreshToken
	token2.ExpiresIn = &expiresIn
	token2.ConsentTid = &consentTid
	token2.CreateDateTime = &createDateTime
	token2.TokenStatus = &tokenStatus
	token2.UpdateDateTime = &updateDateTime
	token2.TokenExpirationDateTime = &tokenExpirationDateTime

	tokens := []Token{token1, token2}

	consent := &Consent{
		AspspId:                        "danske",
		TrackingId:                     "111111111111111111111111111",
		SessionReferenceId:             "222222222222222222222222222",
		ConsentId:                      "123456789",
		ConsentExpirationDateTime:      api.ObTime(time.Now()),
		ConsentTransactionFromDateTime: api.ObTime(time.Now()),
		ConsentTransactionToDateTime:   api.ObTime(time.Now().Add(time.Hour * 24 * 10)),
		ConsentStatusUpdateDateTime:    api.ObTime(time.Now()),
		CreateDateTime:                 api.ObTime(time.Now()),
		UpdateDateTime:                 api.ObTime(time.Now()),
		ConsentStatus:                  api.AwaitingAuthorisation,
		ConsentType:                    api.Account,
		Tokens:                         tokens,
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"save_consent_success",
			fields{db: store.LoadDBConnection()},
			args{consent: consent},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := repositoryWrite{
				db: tt.fields.db,
			}
			if err := r.saveConsent(tt.args.consent); (err != nil) != tt.wantErr {
				t.Errorf("saveConsent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_repository_findByTrackingId(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		trackingId string
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantAccessToken string
		wantErr         bool
	}{
		{
			"find_by_reference_id_success",
			fields{db: store.LoadDBConnection()},
			args{trackingId: "163b39483659429edd305c8c9228a55456a9e42b"},
			"eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := repositoryRead{
				db: tt.fields.db,
			}
			got, err := r.findByTrackingId(tt.args.trackingId)
			if (err != nil) != tt.wantErr {
				t.Errorf("findByTrackingId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil || got.SessionReferenceId != tt.wantAccessToken {
				t.Errorf("findByTrackingId() got = %v, want %v", got, tt.wantAccessToken)
				return
			}
		})
	}
}

func Test_repository_findConsentByCidAndStatus(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		cid    string
		status string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		tokenExist bool
		wantErr    bool
	}{
		{
			"find_by_cid_success",
			fields{db: store.LoadDBConnection()},
			args{cid: "1", status: "Authorised"},
			true,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := repositoryRead{
				db: tt.fields.db,
			}
			got, err := r.findConsentByCidAndStatus(tt.args.cid, tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("findConsentByCidAndStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil || ((got.Tokens != nil) != tt.tokenExist) {
				t.Errorf("findByTrackingId() -> tokenExists? got = %v, want %v", got, tt.tokenExist)
				return
			}
		})
	}
}

func Test_repository_findConsentByUserIdAndTppIdAndStatus(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		userId string
		tppId  string
		status string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			"findConsentByUserIdAndTppIdAndStatus_success",
			fields{db: store.LoadDBConnection()},
			args{
				userId: "kaan",
				tppId:  "Tpp_1",
				status: api.Authorised,
			},
			3,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := repositoryRead{
				db: tt.fields.db,
			}
			got, err := r.findConsentByUserIdAndTppIdAndStatus(tt.args.userId, tt.args.tppId, tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("findConsentByUserIdAndTppIdAndStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil || (len(got) != tt.want) {
				t.Errorf("findConsentByUserIdAndTppIdAndStatus() got = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func Test_repository_findByCid(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		cid string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			"findByCid_success",
			fields{db: store.LoadDBConnection()},
			args{cid: "1"},
			1,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := repositoryRead{
				db: tt.fields.db,
			}
			got, err := r.findByCid(tt.args.cid)
			if (err != nil) != tt.wantErr {
				t.Errorf("findByCid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil || (got.Id != tt.want) {
				t.Errorf("findByCid() got = %v, want %v", got, tt.want)
			}
		})
	}
}
