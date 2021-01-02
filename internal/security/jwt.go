package security

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/kaanaktas/openbanking-accountinformation/api"
	"github.com/kaanaktas/openbanking-accountinformation/internal/cache"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var cacheMem = cache.LoadInMemory()

func init() {
	//fix for the wrong signature problem for PS-*.
	//follow up issue: https://github.com/dgrijalva/jwt-go/pull/305
	jwt.SigningMethodPS256.Options.SaltLength = rsa.PSSSaltLengthEqualsHash
}

func GenerateJwtWithClaims(claims jwt.MapClaims, signingMethod jwt.SigningMethod) (string, error) {
	if signingMethod == jwt.SigningMethodRS256 || signingMethod == jwt.SigningMethodPS256 {
		return signWithPrivateKey(claims, signingMethod)
	}

	if signingMethod == jwt.SigningMethodHS256 {
		return signWithSecret(claims, signingMethod)
	}

	return "", fmt.Errorf("error in GenerateJwtWithClaims(). unsupported signing algorithm: %v", signingMethod.Alg())
}

func GenerateJwtWithJsonString(jsonBody string, signingMethod jwt.SigningMethod) (string, error) {
	if signingMethod == jwt.SigningMethodRS256 || signingMethod == jwt.SigningMethodPS256 {
		return signJsonStringWithPrivateKey(jsonBody, signingMethod)
	}

	if signingMethod == jwt.SigningMethodHS256 {
		return signJsonStringWithSecret(jsonBody, signingMethod)
	}

	return "", fmt.Errorf("error in GenerateJwtWithJsonString(). unsupported signing algorithm: %v", signingMethod.Alg())
}

func signWithSecret(claims jwt.MapClaims, signingMethod jwt.SigningMethod) (string, error) {
	keyData, err := GetSecretKey(api.InternalSignKey, os.Getenv("INTERNAL_SIGN_KEY"))
	if err != nil {
		return "", errors.WithMessage(err, "error in signWithSecret()")
	}

	token := jwt.NewWithClaims(signingMethod, claims)

	return token.SignedString(keyData)
}

func signJsonStringWithSecret(jsonBody string, signingMethod jwt.SigningMethod) (string, error) {
	keyData, err := GetSecretKey(api.InternalSignKey, os.Getenv("INTERNAL_SIGN_KEY"))
	if err != nil {
		return "", errors.WithMessage(err, "error in signJsonStringWithSecret()")
	}

	headers := map[string]interface{}{
		"typ": "JWT",
		"alg": signingMethod.Alg(),
	}

	return signJwtString(headers, jsonBody, signingMethod, keyData)
}

func signWithPrivateKey(claims jwt.MapClaims, signingMethod jwt.SigningMethod) (string, error) {
	key, err := GetPrivateKey(api.ObSignKey, os.Getenv("OB_SIGN_KEY"))
	if err != nil {
		return "", errors.WithMessage(err, "error in signWithPrivateKey()")
	}

	token := jwt.NewWithClaims(signingMethod, claims)
	kid := os.Getenv("KID")
	token.Header["kid"] = kid

	return token.SignedString(key)
}

func signJsonStringWithPrivateKey(jsonBody string, signingMethod jwt.SigningMethod) (string, error) {
	key, err := GetPrivateKey(api.ObSignKey, os.Getenv("OB_SIGN_KEY"))
	if err != nil {
		return "", errors.WithMessage(err, "error in signJsonStringWithPrivateKey()")
	}

	kid := os.Getenv("KID")
	headers := map[string]interface{}{
		"typ": "JWT",
		"alg": signingMethod.Alg(),
		"kid": kid,
	}

	return signJwtString(headers, jsonBody, signingMethod, key)
}

func CreateTokenTime(addMinute time.Duration) int64 {
	return time.Now().Add(time.Minute * addMinute).Unix()
}

func GetPrivateKey(cacheId, keyAddress string) (*rsa.PrivateKey, error) {
	if value, found := cacheMem.Get(cacheId); found {
		return value.(*rsa.PrivateKey), nil
	}

	keyData, err := ioutil.ReadFile(keyAddress)
	if err != nil {
		return nil, errors.WithMessage(err, "error in GetPrivateKey(). couldn't parse private key file")
	}
	var key *rsa.PrivateKey
	key, err = jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		return nil, errors.WithMessage(err, "error in GetPrivateKey(). couldn't parse private key")
	}
	_ = cacheMem.Set(cacheId, key, cache.NoExpiration)

	return key, nil
}

func GetSecretKey(cacheId, keyAddress string) ([]byte, error) {
	if value, found := cacheMem.Get(cacheId); found {
		return value.([]byte), nil
	}

	keyData, err := ioutil.ReadFile(keyAddress)
	if err != nil {
		return nil, errors.WithMessage(err, "error in GetSecretKey(). couldn't parse secret key file")
	}
	_ = cacheMem.Set(cacheId, keyData, cache.NoExpiration)

	return keyData, nil
}

//signs encoded json body.
//This is an alternative of signing with claims, if you want directly sign json string
//Returns a signed Jwt
func signJwtString(headers map[string]interface{}, bodyJson string, signingMethod jwt.SigningMethod, key interface{}) (string, error) {
	var err error
	var headerJson []byte
	if headerJson, err = json.Marshal(headers); err != nil {
		return "", errors.WithMessage(err, "error in signJwtString()")
	}

	parts := make([]string, 2)
	parts[0] = jwt.EncodeSegment(headerJson)
	parts[1] = jwt.EncodeSegment([]byte(bodyJson))
	jwtString := strings.Join(parts, ".")

	sig, err := signingMethod.Sign(jwtString, key)
	if err != nil {
		return "", errors.WithMessage(err, "error in signJwtString()")
	}

	return strings.Join([]string{jwtString, sig}, "."), nil
}

func VerifyJwt(tokenString string, signingMethod jwt.SigningMethod, key interface{}) error {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != signingMethod {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return key, nil
	})
	if err != nil {
		return errors.WithMessage(err, "error in VerifyJwt(). couldn't verify Jwt")
	}

	return nil
}
