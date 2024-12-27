// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package security

import (
	"crypto/rsa"
	"encoding/base64"

	"github.com/golang-jwt/jwt/v4"
)

func CreateSignedJWT(keyId, privateKey any, claims jwt.Claims) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	jwtToken.Header["kid"] = keyId
	return jwtToken.SignedString(privateKey)
}

type JWK struct {
	KeyType   string `json:"kty"`
	Use       string `json:"use"`
	KeyId     string `json:"kid"`
	Modulus   string `json:"n"`
	Exponent  string `json:"e"`
	Algorithm string `json:"alg"`
}

type JWKS struct {
	Keys []*JWK `json:"keys"`
}

// Create a new JWK from the given RSA private key
func JWKFromRSAPrivateKey(keyId string, privateKey *rsa.PrivateKey) *JWK {
	modulesByes := privateKey.N.Bytes()
	modulesB64 := base64.RawURLEncoding.EncodeToString(modulesByes)
	return &JWK{
		KeyType:   "RSA",
		Use:       "sig",
		KeyId:     keyId,
		Algorithm: "RS256",
		Modulus:   modulesB64,
		Exponent:  "AQAB", // REVIEW: should create this from the exponent
	}
}

// Create new Key Set from given keys
func NewJWKS(keys ...*JWK) *JWKS {
	return &JWKS{
		Keys: keys,
	}
}
