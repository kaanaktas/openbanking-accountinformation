package session

import (
	"github.com/google/uuid"
	"github.com/kaanaktas/openbanking-accountinformation/internal/store"
	"reflect"
	"testing"
)

func Test_service_InitiateSession(t *testing.T) {
	type fields struct {
		repo Repository
	}
	type args struct {
		userId string
		tppId  string
		tid    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			"initiate_new_session_success",
			fields{repo: NewRepository(store.LoadDBConnection())},
			args{
				userId: "kaan",
				tppId:  "Tpp_1",
				tid:    uuid.New().String(),
			},
			map[string]interface{}{},
			false,
		},
		{
			"retrieve_session_success",
			fields{repo: NewRepository(store.LoadDBConnection())},
			args{
				userId: "kaan",
				tppId:  "Tpp_1",
				tid:    "1",
			},
			map[string]interface{}{
				"internal_access_token": "123456789",
				"reference_id":          "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service{
				repo: tt.fields.repo,
			}
			got, err := s.InitiateSession(tt.args.userId, tt.args.tppId, tt.args.tid)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitiateSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.name == "retrieve_session_success" && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitiateSession() got = %v, want %v", got, tt.want)
				return
			}
			if tt.name == "initiate_new_session_success" && (got == nil || got["internal_access_token"] == "") {
				t.Errorf("InitiateSession() got = %v, want %v", got, tt.want)
				return
			}
		})
	}
}
