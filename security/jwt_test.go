// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package security

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/tommika/gorilla/assert"
	"github.com/tommika/gorilla/must"
)

// REVIEW: consider moving this to some test utils package
func mustOpen(path string) *os.File {
	return must.NotBeAnError(os.Open(path))
}

// To generate the test key file used here:
// $ openssl genpkey -quiet -algorithm RSA  > test-key-rsa.pem

func readTestRSAPrivateKey() (*rsa.PrivateKey, error) {
	pemFile := mustOpen("./test-data/test-key-rsa.pem")
	defer pemFile.Close()
	return ReadPrivateKeyFromPEMFile[rsa.PrivateKey](pemFile)
}

func TestJWKSFromRandomPrivateKey(t *testing.T) {
	t.Log("generating RSA private key")
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(t, err)
	assert.NotNil(t, privateKey)

	var keyId = "TestJWKSFromRandomPrivateKey"
	t.Logf("generating JWKS from private key: keyId=%s", keyId)
	jwk := JWKFromRSAPrivateKey(keyId, privateKey)
	jwks := NewJWKS(jwk)
	bytes, err := json.Marshal(jwks)
	assert.Nil(t, err)
	assert.NotNil(t, bytes)
	t.Logf("jwks: %+v", string(bytes))
}

func TestCreateSignedJwt(t *testing.T) {
	privateKey, err := readTestRSAPrivateKey()
	assert.Nil(t, err)
	assert.NotNil(t, privateKey)

	var keyId = "TestCreateSignedJwt"
	iat := time.Now()
	exp := iat.Add(24 * time.Hour)
	var claims = jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(iat),
		ExpiresAt: jwt.NewNumericDate(exp),
	}
	tok, err := CreateSignedJWT(keyId, privateKey, claims)
	assert.Nil(t, err)
	t.Logf("jwt: %+v", tok)

}
