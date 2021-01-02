package callback

import (
	"github.com/jmoiron/sqlx"
	"github.com/kaanaktas/openbanking-accountinformation/api"
	"github.com/kaanaktas/openbanking-accountinformation/internal/store"
	"testing"
	"time"
)

func Test_repository_saveResourceAccessAndRefreshToken(t *testing.T) {
	updateParameters := map[string]interface{}{
		"id":                      3,
		"resourceAccessToken":     "987654321",
		"resourceRefreshToken":    "123456789",
		"status":                  api.Authorised,
		"expiresIn":               60,
		"tokenExpirationDateTime": api.ObTime(time.Now().Add(time.Hour * 720)),
		"updateTime":              api.ObTime(time.Now()),
	}

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		parameters map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"save_success",
			fields{db: store.LoadDBConnection()},
			args{parameters: updateParameters},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := repository{
				db: tt.fields.db,
			}
			if err := r.saveResourceAccessAndRefreshToken(tt.args.parameters); (err != nil) != tt.wantErr {
				t.Errorf("saveResourceAccessAndRefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
