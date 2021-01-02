package client

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"
)

func initTest() {
	_ = os.Setenv("CLIENT_CA_CERT_PEM", "./testdata/postman_issuer.cer")
	_ = os.Setenv("CLIENT_CERT_PEM", "./testdata/test_cert.pem")
	_ = os.Setenv("CLIENT_KEY_PEM", "./testdata/test_key.pem")
}

func Test_PostService(t *testing.T) {
	initTest()

	postBody, _ := json.Marshal(map[string]string{
		"name": "test",
	})

	type args struct {
		endpoint string
		payload  io.Reader
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"http_post_expect_200",
			args{"https://postman-echo.com/post",
				bytes.NewReader(postBody)},
			200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewSecureHttpClient(tt.args.endpoint, nil)
			if err != nil {
				t.Fatalf("could not create secure client, %v", err)
			}

			if got, err := client.Post(tt.args.payload); got.StatusCode != tt.want {
				t.Error(err)
				t.Errorf("callService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetService(t *testing.T) {
	initTest()

	type args struct {
		endpoint string
		payload  io.Reader
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"http_post_expect_200",
			args{"https://postman-echo.com/get?foo1=bar1&foo2=bar2",
				nil},
			200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewSecureHttpClient(tt.args.endpoint, nil)
			if err != nil {
				t.Fatalf("could not create secure client, %v", err)
			}

			if got, err := client.Get(nil); got.StatusCode != tt.want {
				t.Error(err)
				t.Errorf("callService() = %v, want %v", got, tt.want)
			}
		})
	}
}
