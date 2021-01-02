package configmanager

import (
	"github.com/kaanaktas/openbanking-accountinformation/internal/cache"
	"github.com/kaanaktas/openbanking-accountinformation/internal/store"
	"testing"
)

func Test_service_FindByConfigName(t *testing.T) {
	repo := NewRepository(store.LoadDBConnection())

	type fields struct {
		repo Repository
		ch   cache.Cache
	}
	type args struct {
		aspspId    string
		configName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			"test",
			fields{repo: repo, ch: cache.LoadInMemory()},
			args{
				aspspId:    "danske",
				configName: "FAPI_FINANCIAL_ID",
			},
			"0015800000jf7AeAAI",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service{
				repository: tt.fields.repo,
				ch:         tt.fields.ch,
			}
			got, err := s.FindByConfigName(tt.args.aspspId, tt.args.configName)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByConfigName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FindByConfigName() got = %v, want %v", got, tt.want)
			}
		})
	}
}
