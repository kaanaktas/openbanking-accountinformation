package security

import (
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

func initTest() {
	_ = os.Setenv("INTERNAL_SIGN_KEY", "./testdata/internal_signing.key")
	_ = os.Setenv("OB_SIGN_KEY", "./testdata/test_key.pem")
	_ = os.Setenv("KID", "kid_test")
}

func Test_GenerateJwtWithClaims(t *testing.T) {
	initTest()

	var claims = map[string]interface{}{
		"grant_types":                  []string{"authorization_code", "refresh_token", "client_credentials"},
		"redirect_uris":                []string{"redirect_uris_1", "redirect_uris_2"},
		"application_type":             "web",
		"iss":                          "iss",
		"token_endpoint_auth_method":   "tls_client_auth",
		"tls_client_auth_dn":           "tls_client_auth_dn",
		"software_id":                  "software_id",
		"software_statement":           "test_ssa",
		"aud":                          "https://obp-api.danskebank.com/open-banking/private",
		"scope":                        "openid accounts payments",
		"jti":                          "40ec08a9-8645-4e4a-ae90-21c473a2a0b8",
		"id_token_signed_response_alg": "PS256",
		"request_object_signing_alg":   "PS256",
		"iat":                          1582717153,
		"exp":                          1582725153,
	}

	type args struct {
		claims           jwt.MapClaims
		signingAlgorithm jwt.SigningMethod
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"sign_with_private_key",
			args{
				claims:           claims,
				signingAlgorithm: jwt.SigningMethodPS256,
			},
			"Token is expired",
			false,
		},
		{"sign_with_secret_key",
			args{
				claims:           claims,
				signingAlgorithm: jwt.SigningMethodHS256,
			},
			"Token is expired",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, vError := GenerateJwtWithClaims(tt.args.claims, tt.args.signingAlgorithm)
			if (vError != nil) != tt.wantErr {
				t.Errorf("generateJwt() error = %v\n, wantErr %v", vError, tt.wantErr)
				return
			}

			if tt.name == "sign_with_private_key" {
				token, vError := jwt.Parse(got, func(t *jwt.Token) (interface{}, error) {
					certData, _ := ioutil.ReadFile("./testdata/test_cert.pem")
					cert, err := jwt.ParseRSAPublicKeyFromPEM(certData)
					if err != nil {
						log.Fatalln("couldn't retrieve the pem file.", err)
					}
					return cert, err
				})

				if token != nil && token.Header["alg"] != "PS256" {
					log.Fatalln("algo doesn't match.")
				}

				if vError != nil && vError.Error() != tt.want {
					t.Errorf("generateJwt() got = %v\n,want = %v", got, tt.want)
				}
			} else {
				token, vError := jwt.Parse(got, func(t *jwt.Token) (interface{}, error) {
					key, err := ioutil.ReadFile("./testdata/internal_signing.key")

					return key, err
				})

				if token != nil && token.Header["alg"] != "HS256" {
					log.Fatalln("algo doesn't match.")
				}

				if vError != nil && vError.Error() != tt.want {
					t.Errorf("generateJwt() got = %v\n,want = %v", got, tt.want)
				}
			}
		})
	}
}

func Test_GenerateJwtWithJsonString(t *testing.T) {
	initTest()

	type args struct {
		jsonBody      string
		signingMethod jwt.SigningMethod
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"sign_json_string_with_private_key",
			args{
				jsonBody:      `{"name":"test", "lastname":"test_2"}`,
				signingMethod: jwt.SigningMethodPS256,
			},
			"eyJhbGciOiJQUzI1NiIsImtpZCI6ImtpZF90ZXN0IiwidHlwIjoiSldUIn0",
			false,
		},
		{"sign_json_string_with_secret_key",
			args{
				jsonBody:      `{"name":"test", "lastname":"test_2"}`,
				signingMethod: jwt.SigningMethodHS256,
			},
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateJwtWithJsonString(tt.args.jsonBody, tt.args.signingMethod)
			if (err != nil) != tt.wantErr {
				t.Errorf("signJsonStringWithPrivateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if strings.Split(got, ".")[0] != tt.want {
				t.Errorf("signJsonStringWithPrivateKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_verifyToken(t *testing.T) {
	initTest()

	keyData, err := ioutil.ReadFile(os.Getenv("INTERNAL_SIGN_KEY"))
	if err != nil {
		t.Errorf("couldn't retrieve key from the file. err: %v", err)
		return
	}

	type args struct {
		token         string
		signingMethod jwt.SigningMethod
		key           interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"verify_token_token_expired",
			args{
				token:         "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDczNzY3MzMsImlhdCI6MTYwNzM3MzEzMywidGlkIjoiMmNmYWZmMjItNDZiYy00NGVhLTk4NDAtZjBiZGNmZWE2NGU2IiwidHBwSWQiOiJUcHBfMSJ9.a2J_Xt6_Gha0-F2UbY4B9I2QGgn4qw8yxMoTLLbRdQ4",
				signingMethod: jwt.SigningMethodHS256,
				key:           keyData,
			},
			"Token is expired",
			true,
		},
		{
			"verify_token_unexpected_signingmethod",
			args{
				token:         "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDczNzY3MzMsImlhdCI6MTYwNzM3MzEzMywidGlkIjoiMmNmYWZmMjItNDZiYy00NGVhLTk4NDAtZjBiZGNmZWE2NGU2IiwidHBwSWQiOiJUcHBfMSJ9.a2J_Xt6_Gha0-F2UbY4B9I2QGgn4qw8yxMoTLLbRdQ4",
				signingMethod: jwt.SigningMethodPS256,
				key:           keyData,
			},
			"unexpected signing method",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got error
			if got = VerifyJwt(tt.args.token, tt.args.signingMethod, tt.args.key); (got != nil) != tt.wantErr {
				t.Errorf("verifyToken() = %v", got)
				return
			}
			if got != nil && !strings.Contains(got.Error(), tt.want) {
				t.Errorf("verifyToken() = %v. want: %v", got, tt.want)
				return
			}
		})
	}
}
