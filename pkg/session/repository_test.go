package session

import (
	"github.com/kaanaktas/openbanking-accountinformation/internal/store"
	"testing"
)

func Test_repository_FindByInternalAccessToken(t *testing.T) {
	repo := NewRepository(store.LoadDBConnection())

	type args struct {
		accessToken string
	}
	tests := []struct {
		name    string
		repo    Repository
		args    args
		want    string
		wantErr bool
	}{
		{"find_by_internal_access_token_success",
			repo,
			args{"123456789"},
			"eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.repo.findByInternalAccessToken(tt.args.accessToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("findByInternalAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil || got.ReferenceId != tt.want {
				t.Errorf("findByInternalAccessToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}
