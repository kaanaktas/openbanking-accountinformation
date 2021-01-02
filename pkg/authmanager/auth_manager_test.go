package authmanager

import (
	"github.com/kaanaktas/openbanking-accountinformation/internal/cache"
	"github.com/kaanaktas/openbanking-accountinformation/internal/store"
	cfg "github.com/kaanaktas/openbanking-accountinformation/pkg/configmanager"
	"github.com/kaanaktas/openbanking-accountinformation/pkg/consent"
	"github.com/kaanaktas/openbanking-accountinformation/pkg/token"
	"reflect"
	"testing"
)

func Test_manager_GetAuthorisedTokenByCid(t *testing.T) {
	type fields struct {
		serviceWrite consent.ServiceWrite
		serviceRead  consent.ServiceRead
		tokenService token.Service
		ch           cache.Cache
	}
	type args struct {
		cid     string
		aspspId string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantToken string
		wantErr   bool
	}{
		{
			"get_authorised_token_by_cid_success",
			fields{
				serviceRead:  consent.NewServiceRead(consent.NewRepositoryRead(store.LoadDBConnection())),
				serviceWrite: consent.NewServiceWrite(consent.NewRepositoryWrite(store.LoadDBConnection())),
				tokenService: token.NewService(cfg.NewService(cfg.NewRepository(store.LoadDBConnection()), cache.LoadInMemory())),
				ch:           cache.LoadInMemory(),
			},
			args{aspspId: "danske", cid: "1"},
			"11111111111111",
			false,
		},
		{
			"get_authorised_token_by_cid_fail",
			fields{
				serviceRead:  consent.NewServiceRead(consent.NewRepositoryRead(store.LoadDBConnection())),
				serviceWrite: consent.NewServiceWrite(consent.NewRepositoryWrite(store.LoadDBConnection())),
				tokenService: token.NewService(cfg.NewService(cfg.NewRepository(store.LoadDBConnection()), cache.LoadInMemory())),
				ch:           cache.LoadInMemory(),
			},
			args{aspspId: "danske", cid: "2"},
			"11111111111111",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := authManager{
				consentServiceRead:  tt.fields.serviceRead,
				consentServiceWrite: tt.fields.serviceWrite,
				tokenService:        tt.fields.tokenService,
				chRedis:             tt.fields.ch,
			}
			got, err := s.GetAuthorisedTokenByCid(tt.args.aspspId, tt.args.cid)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindAuthorisedTokenByCid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.wantToken) {
				t.Errorf("FindAuthorisedTokenByCid() got = %v, want %v", got, tt.wantToken)
			}
		})
	}
}
