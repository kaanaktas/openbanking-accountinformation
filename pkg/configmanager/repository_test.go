package configmanager

import (
	"github.com/kaanaktas/openbanking-accountinformation/internal/store"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func Test_repository_FindByConfigName(t *testing.T) {
	repo := NewRepository(store.LoadDBConnection())

	type args struct {
		aspspId    string
		configName string
	}
	tests := []struct {
		name    string
		repo    Repository
		args    args
		want    string
		wantErr bool
	}{
		{"get_config_success",
			repo,
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
			got, err := repo.findByConfigName(tt.args.aspspId, tt.args.configName)
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
