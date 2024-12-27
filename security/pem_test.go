// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package security

import (
	"crypto/rsa"
	"encoding/json"
	"testing"

	"github.com/tommika/gorilla/assert"
)

func TestJWKSFromPEMFile(t *testing.T) {
	privateKey, err := readTestRSAPrivateKey()
	assert.Nil(t, err)
	assert.NotNil(t, privateKey)

	var keyId = "TestJWKSFromPEMFile"
	t.Logf("generating JWKS from private key: keyId=%s", keyId)
	jwk := JWKFromRSAPrivateKey(keyId, privateKey)
	jwks := NewJWKS(jwk)
	bytes, err := json.Marshal(jwks)
	assert.Nil(t, err)
	assert.NotNil(t, bytes)
	t.Logf("jwks: %+v", string(bytes))
}

func TestJWKSFromBadPEMFile(t *testing.T) {
	// Not a PEM file
	notPemFile := mustOpen("./test-data/not-a-key.txt")
	defer notPemFile.Close()
	_, err := ReadPrivateKeyFromPEMFile[rsa.PrivateKey](notPemFile)
	assert.NotNil(t, err)

	// EOF
	_, err = ReadPrivateKeyFromPEMFile[rsa.PrivateKey](notPemFile)
	assert.NotNil(t, err)

	// Not supported by Go
	unsupportedPemFile := mustOpen("./test-data/test-key-ed488.pem")
	defer unsupportedPemFile.Close()
	_, err = ReadPrivateKeyFromPEMFile[rsa.PrivateKey](unsupportedPemFile)
	assert.NotNil(t, err)

	// Not RSA
	notRsaPemFile := mustOpen("./test-data/test-key-ed25519.pem")
	defer notRsaPemFile.Close()
	_, err = ReadPrivateKeyFromPEMFile[rsa.PrivateKey](notRsaPemFile)
	assert.NotNil(t, err)
}
