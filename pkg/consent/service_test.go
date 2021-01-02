package consent

import (
	"github.com/kaanaktas/openbanking-accountinformation/api"
	"github.com/kaanaktas/openbanking-accountinformation/internal/cache"
	"github.com/kaanaktas/openbanking-accountinformation/internal/store"
	cfg "github.com/kaanaktas/openbanking-accountinformation/pkg/configmanager"
	"github.com/kaanaktas/openbanking-accountinformation/pkg/token"
	"strings"
	"testing"
	"time"
)

func Test_service_CreateConsent(t *testing.T) {
	type fields struct {
		serviceRead  ServiceRead
		serviceWrite ServiceWrite
		tokenService token.Service
		cfg          cfg.Service
	}

	type args struct {
		sessionReferenceId string
		trackingId         string
		aspspId            string
		consent            *ObReadConsent
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"create_consent_success",
			fields{
				serviceRead:  NewServiceRead(NewRepositoryRead(store.LoadDBConnection())),
				serviceWrite: NewServiceWrite(NewRepositoryWrite(store.LoadDBConnection())),
				tokenService: token.NewService(cfg.NewService(cfg.NewRepository(store.LoadDBConnection()), cache.LoadInMemory())),
				cfg:          cfg.NewService(cfg.NewRepository(store.LoadDBConnection()), cache.LoadInMemory()),
			},
			args{
				sessionReferenceId: "",
				trackingId:         "",
				aspspId:            "danske",
				consent: &ObReadConsent{
					Data: &ObReadConsentData{
						Permissions:             []string{"ReadAccountsBasic"},
						ExpirationDateTime:      api.ObTime(time.Now().Add(time.Hour * 24 * 10)),
						TransactionFromDateTime: api.ObTime(time.Now()),
						TransactionToDateTime:   api.ObTime(time.Now()),
					},
					Risk: &ObRisk{},
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewFacade(tt.fields.serviceRead, tt.fields.serviceWrite, tt.fields.tokenService, tt.fields.cfg)

			got, err := s.CreateConsent(tt.args.sessionReferenceId, tt.args.trackingId, tt.args.aspspId, tt.args.consent)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateConsent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && (got == "" || !strings.HasPrefix(got, "https://")) {
				t.Errorf("CreateConsent() got = %v", got)
			}
		})
	}
}

func Test_service_FindConsentByUserIdAndTppIdAndStatus(t *testing.T) {
	type fields struct {
		repo RepositoryRead
	}
	type args struct {
		userId string
		tppId  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			"FindAuthorisedConsentByUserIdAndTppId_success",
			fields{
				repo: NewRepositoryRead(store.LoadDBConnection()),
			},
			args{
				userId: "kaan",
				tppId:  "Tpp_1",
			},
			2,
			false,
		},
		{
			"FindAuthorisedConsentByUserIdAndTppId_expect_no_data",
			fields{
				repo: NewRepositoryRead(store.LoadDBConnection()),
			},
			args{
				userId: "kaan",
				tppId:  "123",
			},
			0,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := serviceRead{
				repo: tt.fields.repo,
			}
			got, err := s.FindAuthorisedConsentByUserIdAndTppId(tt.args.userId, tt.args.tppId)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindConsentByUserIdAndTppIdAndStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("FindConsentByUserIdAndTppIdAndStatus() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_facade_GetConsent(t *testing.T) {
	type fields struct {
		serviceRead  ServiceRead
		serviceWrite ServiceWrite
		tokenService token.Service
		cfg          cfg.Service
	}
	type args struct {
		cid     string
		aspspId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			"GetConsent_success",
			fields{
				serviceRead:  NewServiceRead(NewRepositoryRead(store.LoadDBConnection())),
				serviceWrite: NewServiceWrite(NewRepositoryWrite(store.LoadDBConnection())),
				tokenService: token.NewService(cfg.NewService(cfg.NewRepository(store.LoadDBConnection()), cache.LoadInMemory())),
				cfg:          cfg.NewService(cfg.NewRepository(store.LoadDBConnection()), cache.LoadInMemory()),
			},
			args{
				cid:     "1",
				aspspId: "danske",
			},
			1,
			false,
		},
		{
			"GetConsent_no_data",
			fields{
				serviceRead:  NewServiceRead(NewRepositoryRead(store.LoadDBConnection())),
				serviceWrite: NewServiceWrite(NewRepositoryWrite(store.LoadDBConnection())),
				tokenService: token.NewService(cfg.NewService(cfg.NewRepository(store.LoadDBConnection()), cache.LoadInMemory())),
				cfg:          cfg.NewService(cfg.NewRepository(store.LoadDBConnection()), cache.LoadInMemory()),
			},
			args{
				cid:     "-1",
				aspspId: "danske",
			},
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := facade{
				serviceRead:  NewServiceRead(NewRepositoryRead(store.LoadDBConnection())),
				serviceWrite: NewServiceWrite(NewRepositoryWrite(store.LoadDBConnection())),
				tokenService: tt.fields.tokenService,
				cfg:          tt.fields.cfg,
			}
			got, err := p.GetConsent(tt.args.cid, tt.args.aspspId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConsent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) < tt.want {
				t.Errorf("GetConsent() got = %v, want %v", len(got), tt.want)
			}
		})
	}
}
